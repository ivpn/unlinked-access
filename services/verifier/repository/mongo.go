package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

// mongoSubscription is used for decoding MongoDB documents whose _id is stored
// as a UUID binary (BSON subtype 0x04), which cannot be decoded directly into
// a string field. It mirrors model.Subscription but uses bson.Binary for _id.
type mongoSubscription struct {
	ID          bson.Binary `bson:"_id,omitempty"`
	TokenHash   string      `bson:"token_hash"`
	IsActive    bool        `bson:"is_active"`
	ActiveUntil time.Time   `bson:"active_until"`
	Tier        string      `bson:"tier"`
}

func toModelSubscription(ms mongoSubscription) model.Subscription {
	id, err := uuid.FromBytes(ms.ID.Data)
	if err != nil {
		// Fall back to raw hex representation so callers always get a non-empty ID.
		id = uuid.UUID(ms.ID.Data)
	}
	return model.Subscription{
		ID:          id.String(),
		TokenHash:   ms.TokenHash,
		IsActive:    ms.IsActive,
		ActiveUntil: ms.ActiveUntil,
		Tier:        ms.Tier,
	}
}

func uuidBinaryFromString(id string) (bson.Binary, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return bson.Binary{}, fmt.Errorf("invalid UUID %q: %w", id, err)
	}
	return bson.Binary{Subtype: 0x04, Data: u[:]}, nil
}

const (
	DbConnTimeout       = 20 * time.Second
	DbPingTimeout       = 20 * time.Second
	DbDisconnectTimeout = 3 * time.Second
)

type MongoDB struct {
	cfg            config.NoSQLDBConfig
	Client         *mongo.Client
	DBName         string
	CollectionName string
}

func NewMongoDB(cfg config.Config) (*MongoDB, error) {
	db := &MongoDB{
		cfg:            cfg.NoSQLDB,
		DBName:         cfg.NoSQLDB.Name,
		CollectionName: cfg.NoSQLDB.Collection,
	}
	if err := db.connect(); err != nil {
		return nil, err
	}
	return db, nil
}

func (m *MongoDB) connect() error {
	c := m.cfg

	var uri string
	if strings.HasPrefix(c.Host, "mongodb://") || strings.HasPrefix(c.Host, "mongodb+srv://") {
		uri = c.Host
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", c.Host, c.Port)
	}

	clientOpts := options.Client().ApplyURI(uri)
	if c.User != "" && c.Password != "" {
		clientOpts.SetAuth(buildCredentials(c))
	}

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DbPingTimeout)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	log.Println("MongoDB connection OK")
	m.Client = client
	return nil
}

func buildCredentials(c config.NoSQLDBConfig) options.Credential {
	authSource := c.AuthSource
	if authSource == "" {
		authSource = "admin"
	}
	return options.Credential{
		Username:   c.User,
		Password:   c.Password,
		AuthSource: authSource,
	}
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), DbDisconnectTimeout)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) collection() *mongo.Collection {
	return m.Client.Database(m.DBName).Collection(m.CollectionName)
}

func (m *MongoDB) GetSubscriptions() ([]model.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.collection().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var raw []mongoSubscription
	if err := cursor.All(ctx, &raw); err != nil {
		return nil, err
	}

	subs := make([]model.Subscription, len(raw))
	for i, ms := range raw {
		subs[i] = toModelSubscription(ms)
	}
	return subs, nil
}

func (m *MongoDB) UpdateSubscriptions(subs []model.Subscription) error {
	if len(subs) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	models := make([]mongo.WriteModel, 0, len(subs))
	for _, sub := range subs {
		idBin, err := uuidBinaryFromString(sub.ID)
		if err != nil {
			log.Printf("skipping subscription with invalid ID %q: %v", sub.ID, err)
			continue
		}
		filter := bson.D{{Key: "_id", Value: idBin}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "is_active", Value: sub.IsActive},
			{Key: "active_until", Value: sub.ActiveUntil},
			{Key: "tier", Value: sub.Tier},
			{Key: "updated_at", Value: time.Now()},
		}}}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := m.collection().BulkWrite(ctx, models, opts)
	return err
}
