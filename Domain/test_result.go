package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAnswer represents a single question answer in a test
type TestAnswer struct {
	Question        string         `bson:"question" json:"question"`
	Options         QuestionOption `bson:"options" json:"options"`
	CorrectAnswer   string         `bson:"correct_answer" json:"correct_answer"`
	UserAnswer      string         `bson:"user_answer" json:"user_answer"`
	IsCorrect       bool           `bson:"is_correct" json:"is_correct"`
	ChapterTitle    string         `bson:"chapter_title" json:"chapter_title"`
	Explanation     string         `bson:"explanation" json:"explanation"`
	TimeSpent       int            `bson:"time_spent" json:"time_spent"` // in seconds
}

// TestResult represents a complete test submission
type TestResult struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	NotebookID   primitive.ObjectID `bson:"notebook_id" json:"notebook_id"`
	PrepPilotID  primitive.ObjectID `bson:"prep_pilot_id" json:"prep_pilot_id"`
	TestAnswers  []TestAnswer       `bson:"test_answers" json:"test_answers"`
	Score        float64            `bson:"score" json:"score"`           // percentage score (0-100)
	TotalQuestions int              `bson:"total_questions" json:"total_questions"`
	CorrectAnswers int              `bson:"correct_answers" json:"correct_answers"`
	TotalTimeSpent int              `bson:"total_time_spent" json:"total_time_spent"` // in seconds
	StartedAt    time.Time          `bson:"started_at" json:"started_at"`
	CompletedAt  time.Time          `bson:"completed_at" json:"completed_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// TestResultRepository interface for database operations
type TestResultRepository interface {
	Create(testResult *TestResult) error
	GetByID(id primitive.ObjectID) (*TestResult, error)
	GetByUserID(userID primitive.ObjectID) ([]*TestResult, error)
	GetByNotebookID(notebookID primitive.ObjectID) ([]*TestResult, error)
	GetByPrepPilotID(prepPilotID primitive.ObjectID) ([]*TestResult, error)
	GetByUserAndNotebook(userID, notebookID primitive.ObjectID) ([]*TestResult, error)
	Update(testResult *TestResult) error
	Delete(id primitive.ObjectID) error
}

// TestResultUseCase interface for business logic
type TestResultUseCase interface {
	SubmitTestResult(userID string, testResult *TestResult) error
	GetTestResultByID(userID string, testResultID string) (*TestResult, error)
	GetUserTestResults(userID string) ([]*TestResult, error)
	GetNotebookTestResults(userID string, notebookID string) ([]*TestResult, error)
	GetTestResultStats(userID string, notebookID string) (*TestStats, error)
}

// TestStats represents aggregated test statistics
type TestStats struct {
	TotalTests      int     `json:"total_tests"`
	AverageScore    float64 `json:"average_score"`
	BestScore       float64 `json:"best_score"`
	WorstScore      float64 `json:"worst_score"`
	TotalTimeSpent  int     `json:"total_time_spent"`
	AverageTime     float64 `json:"average_time"`
	ImprovementRate float64 `json:"improvement_rate"` // percentage improvement from first to last test
}
