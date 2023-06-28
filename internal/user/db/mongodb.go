package db

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rest-api-tutorial/internal/apperror"
	"rest-api-tutorial/internal/user"
	"rest-api-tutorial/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d db) FindAll(ctx context.Context) (u []user.User, e error) {

	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return nil, fmt.Errorf("failed to find all users: %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. err: %v", err)
	}

	return u, nil
}

func (d db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}

	d.logger.Debug("convert InsertedId to objectId")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert object_id to hex. probably object_id: %s. due to error: %v", oid, err)
}

func (d db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to object_id. Hex: %s, due to error: %v", id, err)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {

		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to find user by id: %s", id)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to deccode user(id: %s) from db, due to error: %v", oid, err)
	}
	return u, nil
}

func (d db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert userID to object_id. userID: %s, due to error: %v", user.ID, err)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err2 := bson.Marshal(user)
	if err2 != nil {
		return fmt.Errorf("failed to marshal user, due to error: %v", err2)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal userBytes, due to error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err3 := d.collection.UpdateOne(ctx, filter, update)
	if err3 != nil {
		return fmt.Errorf("failed to execute updating user query, due to error: %v", err3)
	}

	if result.MatchedCount == 0 {
		//todo ErrEntityNotFound
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("matched %d documents and modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert userID to object_id. userID: %s, due to error: %v", id, err)
	}

	filter := bson.M{"_id": objectID}

	result, err2 := d.collection.DeleteOne(ctx, filter)
	if err2 != nil {
		return fmt.Errorf("failed to execute deleting user query, due to error: %v", err2)
	}

	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
