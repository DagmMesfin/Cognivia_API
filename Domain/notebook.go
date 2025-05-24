package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notebook struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID  `bson:"user_id" json:"user_id"`
	SnapnotesID     *primitive.ObjectID `bson:"snapnotes_id,omitempty" json:"snapnotes_id,omitempty"`
	PrepPilotID     *primitive.ObjectID `bson:"prep_pilot_id,omitempty" json:"prep_pilot_id,omitempty"`
	Name            string              `bson:"name" json:"name" binding:"required"`
	Icon            string              `bson:"icon" json:"icon"`
	Color           string              `bson:"color" json:"color"`
	Type            string              `bson:"type" json:"type"`
	GoogleDriveLink *string             `bson:"google_drive_link,omitempty" json:"google_drive_link,omitempty"`
	CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at" json:"updated_at"`
}

type UpdateRequest struct {
	Name            string  `json:"name"`
	Icon            string  `json:"icon"`
	Color           string  `json:"color"`
	Type            string  `json:"type"`
	GoogleDriveLink *string `json:"google_drive_link"`
}

type NotebookRepository interface {
	Create(notebook *Notebook) error
	GetByID(id primitive.ObjectID) (*Notebook, error)
	GetByUserID(userID primitive.ObjectID) ([]*Notebook, error)
	Update(notebook *Notebook) error
	Delete(id primitive.ObjectID) error
	// Add other repository methods as needed (e.g., GetByUserID)
}

type NotebookUseCase interface {
	CreateNotebook(userID string, notebook *Notebook) error
	GetNotebookByID(userID string, notebookID string) (*Notebook, error)
	GetNotebooksByUserID(userID string) ([]*Notebook, error)
	GetSnapnotes(userID string, notebookID string) (*Snapnotes, error)
	GetPrepPilot(userID string, notebookID string) (*PrepPilot, error)
	UpdateNotebook(notebookID string, notebook UpdateRequest) error
	DeleteNotebook(userID string, notebookID string) error
	// Add other use case methods as needed
}
