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

type notebookRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewNotebookRepository(db *mongo.Database) domain.NotebookRepository {
	return &notebookRepository{
		db:         db,
		collection: db.Collection("notebooks"),
	}
}

func (r *notebookRepository) Create(notebook *domain.Notebook) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate new ObjectID for the notebook
	notebook.ID = primitive.NewObjectID()
	notebook.CreatedAt = time.Now()
	notebook.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, notebook)
	if err != nil {
		return err
	}

	return nil
}

func (r *notebookRepository) GetByID(id primitive.ObjectID) (*domain.Notebook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var notebook domain.Notebook
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&notebook)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &notebook, nil
}

func (r *notebookRepository) GetByUserID(userID primitive.ObjectID) ([]*domain.Notebook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var notebooks []*domain.Notebook
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &notebooks); err != nil {
		return nil, err
	}

	return notebooks, nil
}

func (r *notebookRepository) Update(notebook *domain.Notebook) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	notebook.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": notebook.ID},
		notebook,
	)
	return err
}

func (r *notebookRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Add other repository methods as needed
