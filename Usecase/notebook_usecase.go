package usecase

import (
	"errors"

	domain "cognivia-api/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type notebookUseCase struct {
	notebookRepo  domain.NotebookRepository
	snapnotesRepo domain.SnapnotesRepository
	prepPilotRepo domain.PrepPilotRepository
}

func NewNotebookUseCase(
	notebookRepo domain.NotebookRepository,
	snapnotesRepo domain.SnapnotesRepository,
	prepPilotRepo domain.PrepPilotRepository,
) domain.NotebookUseCase {
	return &notebookUseCase{
		notebookRepo:  notebookRepo,
		snapnotesRepo: snapnotesRepo,
		prepPilotRepo: prepPilotRepo,
	}
}

func (u *notebookUseCase) CreateNotebook(userID string, notebook *domain.Notebook) error {

	// Convert userID string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	notebook.UserID = objectID

	// Add any business logic here before creating the notebook
	return u.notebookRepo.Create(notebook)
}

func (u *notebookUseCase) GetNotebookByID(userID string, notebookID string) (*domain.Notebook, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	objectNotebookID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		return nil, err
	}

	notebook, err := u.notebookRepo.GetByID(objectNotebookID)
	if err != nil {
		return nil, err
	}

	if notebook == nil || notebook.UserID != objectUserID {
		return nil, errors.New("notebook not found or does not belong to user")
	}

	return notebook, nil
}

func (u *notebookUseCase) UpdateNotebook(notebookID string, notebook domain.UpdateRequest) error {
	objectNotebookID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		return err
	}

	existingNotebook, err := u.notebookRepo.GetByID(objectNotebookID)
	if err != nil {
		return err
	}

	if existingNotebook == nil {
		return errors.New("notebook not found")
	}

	// Update allowed fields
	if notebook.Name != nil {
		existingNotebook.Name = *notebook.Name
	}
	if notebook.Icon != nil {
		existingNotebook.Icon = *notebook.Icon
	}
	if notebook.Color != nil {
		existingNotebook.Color = *notebook.Color
	}
	if notebook.Type != nil {
		existingNotebook.Type = *notebook.Type
	}
	if notebook.GoogleDriveLink != nil {
		existingNotebook.GoogleDriveLink = notebook.GoogleDriveLink
	}
	if notebook.SnapnotesID != nil {
		snapnotesID, err := primitive.ObjectIDFromHex(*notebook.SnapnotesID)
		if err != nil {
			return err
		}
		existingNotebook.SnapnotesID = &snapnotesID
	}
	if notebook.PrepPilotID != nil {
		prepPilotID, err := primitive.ObjectIDFromHex(*notebook.PrepPilotID)
		if err != nil {
			return err
		}
		existingNotebook.PrepPilotID = &prepPilotID
	}

	return u.notebookRepo.Update(existingNotebook)
}

func (u *notebookUseCase) DeleteNotebook(userID string, notebookID string) error {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	objectNotebookID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		return err
	}

	existingNotebook, err := u.notebookRepo.GetByID(objectNotebookID)
	if err != nil {
		return err
	}

	if existingNotebook == nil || existingNotebook.UserID != objectUserID {
		return errors.New("notebook not found or does not belong to user")
	}

	return u.notebookRepo.Delete(objectNotebookID)
}

func (u *notebookUseCase) GetNotebooksByUserID(userID string) ([]*domain.Notebook, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return u.notebookRepo.GetByUserID(objectUserID)
}

func (u *notebookUseCase) GetSnapnotes(userID string, notebookID string) (*domain.Snapnotes, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	objectNotebookID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		return nil, err
	}

	notebook, err := u.notebookRepo.GetByID(objectNotebookID)
	if err != nil {
		return nil, err
	}

	if notebook == nil || notebook.UserID != objectUserID {
		return nil, errors.New("notebook not found or does not belong to user")
	}

	if notebook.SnapnotesID == nil {
		return nil, errors.New("no snapnotes associated with this notebook")
	}

	return u.snapnotesRepo.GetByID(*notebook.SnapnotesID)
}

func (u *notebookUseCase) GetPrepPilot(userID string, notebookID string) (*domain.PrepPilot, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	objectNotebookID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		return nil, err
	}

	notebook, err := u.notebookRepo.GetByID(objectNotebookID)
	if err != nil {
		return nil, err
	}

	if notebook == nil || notebook.UserID != objectUserID {
		return nil, errors.New("notebook not found or does not belong to user")
	}

	if notebook.PrepPilotID == nil {
		return nil, errors.New("no prep pilot associated with this notebook")
	}

	return u.prepPilotRepo.GetByID(*notebook.PrepPilotID)
}

// Add other use case methods as needed
