package mongodb

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"manager/api/hash"
	"manager/internal/model"
	"manager/internal/store"
	"time"
)

type RequestRepository struct {
	requestsCollection *mongo.Collection
	logger             *logrus.Logger
}

func NewRequestRepository(db *mongo.Database, collectionName string, logger *logrus.Logger) store.RequestRepository {
	return &RequestRepository{
		requestsCollection: db.Collection(collectionName),
		logger:             logger,
	}
}

func (rr RequestRepository) Create(Hash string, MaxLength int) (string, error) {
	req := &model.HashRequest{
		ID:        "",
		Hash:      Hash,
		MaxLength: MaxLength,
		Data:      []string{},
		Status:    model.Created,
		DateTime:  time.Now(),
	}

	res, err := rr.requestsCollection.InsertOne(context.Background(), req)
	if err != nil {
		rr.logger.Infof("mongodb error: failed to create user due to error: %v", err)
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		rr.logger.Infof("mongodb error: failed to convert InsertedID to ObjectID")
		return "", err
	}

	return oid.Hex(), nil
}

func (rr RequestRepository) GetRequestById(id string) (*model.HashRequest, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rr.logger.Infof("mongodb error: failed to convert user ID to ObgectID. ID=%s", id)
		return nil, err
	}

	filter := bson.M{"_id": oid}

	res := rr.requestsCollection.FindOne(context.Background(), filter)
	if res.Err() != nil {
		rr.logger.Infof("mongodb error: failed to find user with ID: %s", id)
		return nil, store.ErrorRecordNotFound
	}

	req := &model.HashRequest{}

	if err = res.Decode(&req); err != nil {
		rr.logger.Infof("mongodb error: failed to decode user from DB due to error %v", err)
		return nil, err
	}

	return req, nil
}

func (rr RequestRepository) GetRequestsByStatus(status string) ([]*model.HashRequest, error) {
	filter := bson.M{"status": status}

	cursor, err := rr.requestsCollection.Find(context.Background(), filter)
	if err != nil {
		rr.logger.Infof("mongodb error: failed to find user with status: %s", status)
		return nil, store.ErrorRecordNotFound
	}

	var results []*model.HashRequest
	if err = cursor.All(context.Background(), &results); err != nil {
		rr.logger.Infof("mongodb error: failed to decode document into results: %s", status)
	}
	for _, result := range results {
		fmt.Println(result)
	}

	return results, nil
}

func (rr RequestRepository) SetStatus(id, status string) (*model.HashRequest, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rr.logger.Infof("mongodb error: failed to convert user ID to ObgectID. ID=%s", id)
	}

	filter := bson.M{"_id": oid}

	update := bson.D{
		{"$set", bson.D{
			{"status", status},
		}},
	}

	_, err = rr.requestsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		rr.logger.Infof("mongodb error: failed to execute update user query. error: %v", err)
	}

	return rr.GetRequestById(id)
}

func (rr RequestRepository) SetResults(r *hash.CrackHashWorkerResponse) error {
	oid, err := primitive.ObjectIDFromHex(r.RequestId)
	if err != nil {
		rr.logger.Infof("mongodb error: failed to convert user ID to ObgectID. ID=%s", r.RequestId)
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$push": bson.M{
			"data": bson.M{
				"$each": r.Answers.Words}},
	}

	res, err := rr.requestsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		rr.logger.Infof("mongodb error: failed to execute update user query. error: %v", err)
		return err
	}
	if res.MatchedCount == 0 {
		rr.logger.Infof("mongodb error: matched and replaced an existing document")
		return err
	}

	_, _ = rr.SetStatus(r.RequestId, string(model.Ready))

	return nil
}
