package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChapterSummary struct {
	ChapterTitle string   `bson:"chapterTitle" json:"chapterTitle"`
	Summary      string   `bson:"summary" json:"summary"`
	KeyPoints    []string `bson:"keyPoints" json:"keyPoints"`
}

type Flashcard struct {
	KeyTerm    string `bson:"key term" json:"key term"`
	Definition string `bson:"definition" json:"definition"`
}

type ChapterFlashcards struct {
	ChapterTitle string      `bson:"chapterTitle" json:"chapterTitle"`
	Flashcards   []Flashcard `bson:"flashcards" json:"flashcards"`
}

type Snapnotes struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Title            string              `bson:"title" json:"title"`
	SummaryByChapter []ChapterSummary    `bson:"summaryByChapter" json:"summaryByChapter"`
	Flashcards       []ChapterFlashcards `bson:"flashcards" json:"flashcards"`
}

type SnapnotesRepository interface {
	GetByID(id primitive.ObjectID) (*Snapnotes, error)
	Create(content *Snapnotes) error
}
