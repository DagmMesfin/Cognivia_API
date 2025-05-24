package controllers

import (
	"net/http"

	domain "cognivia-api/Domain"

	"github.com/gin-gonic/gin"
)

type NotebookHandler struct {
	notebookUseCase domain.NotebookUseCase
}

func NewNotebookHandler(notebookUseCase domain.NotebookUseCase) *NotebookHandler {
	return &NotebookHandler{
		notebookUseCase: notebookUseCase,
	}
}

func (h *NotebookHandler) CreateNotebook(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var notebook domain.Notebook
	if err := c.ShouldBindJSON(&notebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pass userID to the use case
	if err := h.notebookUseCase.CreateNotebook(userID.(string), &notebook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created notebook details
	c.JSON(http.StatusCreated, notebook)
}

func (h *NotebookHandler) GetNotebook(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("id")

	notebook, err := h.notebookUseCase.GetNotebookByID(userID.(string), notebookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if notebook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found"})
		return
	}

	c.JSON(http.StatusOK, notebook)
}

func (h *NotebookHandler) GetNotebooksByUserID(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebooks, err := h.notebookUseCase.GetNotebooksByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notebooks)
}

func (h *NotebookHandler) UpdateNotebook(c *gin.Context) {
	notebookID := c.Param("id")

	updateReq := domain.UpdateRequest{}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notebookUseCase.UpdateNotebook(notebookID, updateReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notebook updated successfully"})
}

func (h *NotebookHandler) DeleteNotebook(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("id")

	if err := h.notebookUseCase.DeleteNotebook(userID.(string), notebookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notebook deleted successfully"})
}

// GetSnapnotes returns the snapnotes for a notebook
func (h *NotebookHandler) GetSnapnotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("id")
	if notebookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notebook ID is required"})
		return
	}

	snapnotes, err := h.notebookUseCase.GetSnapnotes(userID.(string), notebookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, snapnotes)
}

// GetPrepPilot returns the prep pilot for a notebook
func (h *NotebookHandler) GetPrepPilot(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("id")
	if notebookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notebook ID is required"})
		return
	}

	prepPilot, err := h.notebookUseCase.GetPrepPilot(userID.(string), notebookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prepPilot)
}

// Add other handler methods as needed
