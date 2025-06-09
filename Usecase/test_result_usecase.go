package usecase

import (
	"errors"
	"math"
	"time"

	domain "cognivia-api/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testResultUseCase struct {
	testResultRepo domain.TestResultRepository
	notebookRepo   domain.NotebookRepository
	prepPilotRepo  domain.PrepPilotRepository
}

func NewTestResultUseCase(
	testResultRepo domain.TestResultRepository,
	notebookRepo domain.NotebookRepository,
	prepPilotRepo domain.PrepPilotRepository,
) domain.TestResultUseCase {
	return &testResultUseCase{
		testResultRepo: testResultRepo,
		notebookRepo:   notebookRepo,
		prepPilotRepo:  prepPilotRepo,
	}
}

func (u *testResultUseCase) SubmitTestResult(userID string, testResult *domain.TestResult) error {
	// Convert userID string to ObjectID
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Validate that the notebook belongs to the user
	notebook, err := u.notebookRepo.GetByID(testResult.NotebookID)
	if err != nil {
		return err
	}
	if notebook == nil || notebook.UserID != objectUserID {
		return errors.New("notebook not found or does not belong to user")
	}

	// Validate that the prep pilot exists and belongs to the notebook
	prepPilot, err := u.prepPilotRepo.GetByID(testResult.PrepPilotID)
	if err != nil {
		return err
	}
	if prepPilot == nil || prepPilot.NotebookID != testResult.NotebookID {
		return errors.New("prep pilot not found or does not belong to notebook")
	}

	// Set the user ID
	testResult.UserID = objectUserID

	// Calculate score and validate answers
	correctAnswers := 0
	totalQuestions := len(testResult.TestAnswers)
	totalTimeSpent := 0

	for i := range testResult.TestAnswers {
		answer := &testResult.TestAnswers[i]
		
		// Check if the answer is correct
		answer.IsCorrect = answer.UserAnswer == answer.CorrectAnswer
		if answer.IsCorrect {
			correctAnswers++
		}
		
		// Add to total time spent
		totalTimeSpent += answer.TimeSpent
	}

	// Calculate score as percentage
	var score float64
	if totalQuestions > 0 {
		score = (float64(correctAnswers) / float64(totalQuestions)) * 100
	}

	// Set calculated values
	testResult.Score = math.Round(score*100) / 100 // Round to 2 decimal places
	testResult.TotalQuestions = totalQuestions
	testResult.CorrectAnswers = correctAnswers
	testResult.TotalTimeSpent = totalTimeSpent

	// Set completion time if not already set
	if testResult.CompletedAt.IsZero() {
		testResult.CompletedAt = time.Now()
	}

	// Create the test result
	return u.testResultRepo.Create(testResult)
}

func (u *testResultUseCase) GetTestResultByID(userID string, testResultID string) (*domain.TestResult, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	objectTestResultID, err := primitive.ObjectIDFromHex(testResultID)
	if err != nil {
		return nil, err
	}

	testResult, err := u.testResultRepo.GetByID(objectTestResultID)
	if err != nil {
		return nil, err
	}

	if testResult == nil || testResult.UserID != objectUserID {
		return nil, errors.New("test result not found or does not belong to user")
	}

	return testResult, nil
}

func (u *testResultUseCase) GetUserTestResults(userID string) ([]*domain.TestResult, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return u.testResultRepo.GetByUserID(objectUserID)
}

func (u *testResultUseCase) GetNotebookTestResults(userID string, notebookID string) ([]*domain.TestResult, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	objectNotebookID, err := primitive.ObjectIDFromHex(notebookID)
	if err != nil {
		return nil, err
	}

	// Validate that the notebook belongs to the user
	notebook, err := u.notebookRepo.GetByID(objectNotebookID)
	if err != nil {
		return nil, err
	}
	if notebook == nil || notebook.UserID != objectUserID {
		return nil, errors.New("notebook not found or does not belong to user")
	}

	return u.testResultRepo.GetByUserAndNotebook(objectUserID, objectNotebookID)
}

func (u *testResultUseCase) GetTestResultStats(userID string, notebookID string) (*domain.TestStats, error) {
	testResults, err := u.GetNotebookTestResults(userID, notebookID)
	if err != nil {
		return nil, err
	}

	if len(testResults) == 0 {
		return &domain.TestStats{}, nil
	}

	// Calculate statistics
	totalTests := len(testResults)
	var totalScore, totalTime float64
	var bestScore, worstScore float64
	bestScore = -1
	worstScore = 101

	for _, result := range testResults {
		totalScore += result.Score
		totalTime += float64(result.TotalTimeSpent)

		if result.Score > bestScore {
			bestScore = result.Score
		}
		if result.Score < worstScore {
			worstScore = result.Score
		}
	}

	averageScore := totalScore / float64(totalTests)
	averageTime := totalTime / float64(totalTests)

	// Calculate improvement rate (compare first and last test)
	var improvementRate float64
	if totalTests >= 2 {
		// testResults are sorted by created_at descending, so reverse order for improvement
		firstScore := testResults[totalTests-1].Score  // oldest test
		lastScore := testResults[0].Score              // newest test
		
		if firstScore > 0 {
			improvementRate = ((lastScore - firstScore) / firstScore) * 100
		}
	}

	return &domain.TestStats{
		TotalTests:      totalTests,
		AverageScore:    math.Round(averageScore*100) / 100,
		BestScore:       bestScore,
		WorstScore:      worstScore,
		TotalTimeSpent:  int(totalTime),
		AverageTime:     math.Round(averageTime*100) / 100,
		ImprovementRate: math.Round(improvementRate*100) / 100,
	}, nil
}
