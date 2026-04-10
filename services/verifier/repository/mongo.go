package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

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

	var subs []model.Subscription
	if err := cursor.All(ctx, &subs); err != nil {
		return nil, err
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
		filter := bson.D{{Key: "_id", Value: sub.ID}}
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
