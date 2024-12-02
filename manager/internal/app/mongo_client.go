package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewClient(ctx context.Context, databaseName, url string, logger *logrus.Logger) (*mongo.Client, *mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(url)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Infof("mongodb connect error: %v", err)
		return nil, nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Infof("mongodb unavailable, ping error: %v", err)
		return nil, nil, err
	}

	return client, client.Database(databaseName), nil
}
