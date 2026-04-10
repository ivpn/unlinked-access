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
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type MongoDB struct {
	Client         *mongo.Client
	DBName         string
	CollectionName string
}

func NewMongoDB(cfg config.Config) (*MongoDB, error) {
	c := cfg.NoSQLDB

	var uri string
	if strings.HasPrefix(c.Host, "mongodb://") || strings.HasPrefix(c.Host, "mongodb+srv://") {
		// HOST is already a full connection string — use it directly
		uri = c.Host
	} else {
		authSource := c.AuthSource
		if authSource == "" {
			authSource = "admin"
		}
		uri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/?authSource=%s",
			c.User, c.Password, c.Host, c.Port, authSource,
		)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	log.Println("MongoDB connection OK")

	return &MongoDB{
		Client:         client,
		DBName:         c.Name,
		CollectionName: c.Collection,
	}, nil
}

func (m *MongoDB) Close() error {
	return m.Client.Disconnect(context.Background())
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
