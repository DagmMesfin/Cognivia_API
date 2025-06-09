package controllers

import (
	"net/http"
	"time"

	domain "cognivia-api/Domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestResultHandler struct {
	testResultUseCase domain.TestResultUseCase
}

func NewTestResultHandler(testResultUseCase domain.TestResultUseCase) *TestResultHandler {
	return &TestResultHandler{
		testResultUseCase: testResultUseCase,
	}
}

// SubmitTestResult handles POST /api/v1/test-results
func (h *TestResultHandler) SubmitTestResult(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var testResult domain.TestResult
	if err := c.ShouldBindJSON(&testResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if testResult.NotebookID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notebook_id is required"})
		return
	}

	if testResult.PrepPilotID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prep_pilot_id is required"})
		return
	}

	if len(testResult.TestAnswers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test_answers cannot be empty"})
		return
	}

	// Set started_at if not provided
	if testResult.StartedAt.IsZero() {
		testResult.StartedAt = time.Now()
	}

	// Submit the test result
	if err := h.testResultUseCase.SubmitTestResult(userID.(string), &testResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Test result submitted successfully",
		"test_result": gin.H{
			"id":              testResult.ID.Hex(),
			"score":           testResult.Score,
			"correct_answers": testResult.CorrectAnswers,
			"total_questions": testResult.TotalQuestions,
			"total_time_spent": testResult.TotalTimeSpent,
		},
	})
}

// GetTestResult handles GET /api/v1/test-results/:id
func (h *TestResultHandler) GetTestResult(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	testResultID := c.Param("id")
	if testResultID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Test result ID is required"})
		return
	}

	testResult, err := h.testResultUseCase.GetTestResultByID(userID.(string), testResultID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testResult)
}

// GetUserTestResults handles GET /api/v1/test-results/user
func (h *TestResultHandler) GetUserTestResults(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	testResults, err := h.testResultUseCase.GetUserTestResults(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testResults)
}

// GetNotebookTestResults handles GET /api/v1/test-results/notebook/:notebook_id
func (h *TestResultHandler) GetNotebookTestResults(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("notebook_id")
	if notebookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notebook ID is required"})
		return
	}

	testResults, err := h.testResultUseCase.GetNotebookTestResults(userID.(string), notebookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testResults)
}

// GetTestResultStats handles GET /api/v1/test-results/notebook/:notebook_id/stats
func (h *TestResultHandler) GetTestResultStats(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("notebook_id")
	if notebookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notebook ID is required"})
		return
	}

	stats, err := h.testResultUseCase.GetTestResultStats(userID.(string), notebookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// TestResultRequest represents the request structure for submitting test results
type TestResultRequest struct {
	NotebookID     string                `json:"notebook_id" binding:"required"`
	PrepPilotID    string                `json:"prep_pilot_id" binding:"required"`
	TestAnswers    []domain.TestAnswer   `json:"test_answers" binding:"required"`
	StartedAt      *time.Time            `json:"started_at,omitempty"`
	CompletedAt    *time.Time            `json:"completed_at,omitempty"`
}

// SubmitTestResultV2 handles POST /api/v1/test-results with string IDs
func (h *TestResultHandler) SubmitTestResultV2(c *gin.Context) {
	// Extract user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var request TestResultRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert string IDs to ObjectIDs
	notebookID, err := primitive.ObjectIDFromHex(request.NotebookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notebook_id format"})
		return
	}

	prepPilotID, err := primitive.ObjectIDFromHex(request.PrepPilotID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prep_pilot_id format"})
		return
	}

	// Create test result
	testResult := domain.TestResult{
		NotebookID:  notebookID,
		PrepPilotID: prepPilotID,
		TestAnswers: request.TestAnswers,
	}

	// Set timestamps
	if request.StartedAt != nil {
		testResult.StartedAt = *request.StartedAt
	} else {
		testResult.StartedAt = time.Now()
	}

	if request.CompletedAt != nil {
		testResult.CompletedAt = *request.CompletedAt
	} else {
		testResult.CompletedAt = time.Now()
	}

	// Submit the test result
	if err := h.testResultUseCase.SubmitTestResult(userID.(string), &testResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Test result submitted successfully",
		"test_result": gin.H{
			"id":               testResult.ID.Hex(),
			"score":            testResult.Score,
			"correct_answers":  testResult.CorrectAnswers,
			"total_questions":  testResult.TotalQuestions,
			"total_time_spent": testResult.TotalTimeSpent,
		},
	})
}
