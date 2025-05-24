package mongodb

import (
	"context"
	"errors"
	"time"

	domain "cognivia-api/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type prepPilotRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewPrepPilotRepository(db *mongo.Database) domain.PrepPilotRepository {
	return &prepPilotRepository{
		db:         db,
		collection: db.Collection("prep_pilot"),
	}
}

func (r *prepPilotRepository) GetByID(id primitive.ObjectID) (*domain.PrepPilot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var prepPilot domain.PrepPilot
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&prepPilot)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &prepPilot, nil
}

func (r *prepPilotRepository) GetByNotebookID(notebookID primitive.ObjectID) (*domain.PrepPilot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var prepPilot domain.PrepPilot
	err := r.collection.FindOne(ctx, bson.M{"notebook_id": notebookID}).Decode(&prepPilot)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &prepPilot, nil
}

func (r *prepPilotRepository) Create(prepPilot *domain.PrepPilot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, prepPilot)
	if err != nil {
		return err
	}
	prepPilot.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *prepPilotRepository) Update(prepPilot *domain.PrepPilot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": prepPilot.ID},
		bson.M{"$set": prepPilot},
	)
	return err
}

func (r *prepPilotRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
