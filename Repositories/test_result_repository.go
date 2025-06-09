package mongodb

import (
	"context"
	"errors"
	"time"

	domain "cognivia-api/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testResultRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewTestResultRepository(db *mongo.Database) domain.TestResultRepository {
	return &testResultRepository{
		db:         db,
		collection: db.Collection("test_results"),
	}
}

func (r *testResultRepository) Create(testResult *domain.TestResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testResult.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, testResult)
	if err != nil {
		return err
	}

	testResult.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *testResultRepository) GetByID(id primitive.ObjectID) (*domain.TestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var testResult domain.TestResult
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&testResult)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &testResult, nil
}

func (r *testResultRepository) GetByUserID(userID primitive.ObjectID) ([]*domain.TestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Sort by created_at descending to get most recent tests first
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var testResults []*domain.TestResult
	for cursor.Next(ctx) {
		var testResult domain.TestResult
		if err := cursor.Decode(&testResult); err != nil {
			return nil, err
		}
		testResults = append(testResults, &testResult)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return testResults, nil
}

func (r *testResultRepository) GetByNotebookID(notebookID primitive.ObjectID) ([]*domain.TestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"notebook_id": notebookID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var testResults []*domain.TestResult
	for cursor.Next(ctx) {
		var testResult domain.TestResult
		if err := cursor.Decode(&testResult); err != nil {
			return nil, err
		}
		testResults = append(testResults, &testResult)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return testResults, nil
}

func (r *testResultRepository) GetByPrepPilotID(prepPilotID primitive.ObjectID) ([]*domain.TestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"prep_pilot_id": prepPilotID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var testResults []*domain.TestResult
	for cursor.Next(ctx) {
		var testResult domain.TestResult
		if err := cursor.Decode(&testResult); err != nil {
			return nil, err
		}
		testResults = append(testResults, &testResult)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return testResults, nil
}

func (r *testResultRepository) GetByUserAndNotebook(userID, notebookID primitive.ObjectID) ([]*domain.TestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":     userID,
		"notebook_id": notebookID,
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var testResults []*domain.TestResult
	for cursor.Next(ctx) {
		var testResult domain.TestResult
		if err := cursor.Decode(&testResult); err != nil {
			return nil, err
		}
		testResults = append(testResults, &testResult)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return testResults, nil
}

func (r *testResultRepository) Update(testResult *domain.TestResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": testResult.ID}
	update := bson.M{"$set": testResult}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *testResultRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
