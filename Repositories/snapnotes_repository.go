package mongodb

import (
	domain "cognivia-api/Domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type snapnotesRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewSnapnotesRepository(db *mongo.Database) domain.SnapnotesRepository {
	return &snapnotesRepository{
		db:         db,
		collection: db.Collection("snapnotes"),
	}
}

func (r *snapnotesRepository) GetByID(id primitive.ObjectID) (*domain.Snapnotes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var snapnotes domain.Snapnotes
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&snapnotes)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &snapnotes, nil
}

func (r *snapnotesRepository) Create(snapnotes *domain.Snapnotes) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, snapnotes)
	if err != nil {
		return err
	}
	snapnotes.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
