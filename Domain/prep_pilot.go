package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionOption struct {
	A string `bson:"A" json:"A"`
	B string `bson:"B" json:"B"`
	C string `bson:"C" json:"C"`
	D string `bson:"D" json:"D"`
}

type Question struct {
	Question    string         `bson:"question" json:"question"`
	Options     QuestionOption `bson:"options" json:"options"`
	Answer      string         `bson:"answer" json:"answer"`
	Explanation string         `bson:"explanation" json:"explanation"`
}

type Chapter struct {
	ChapterTitle string     `bson:"chapterTitle" json:"chapterTitle"`
	Questions    []Question `bson:"questions" json:"questions"`
}

type PrepPilot struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NotebookID primitive.ObjectID `bson:"notebook_id" json:"notebook_id"`
	Chapters   []Chapter          `bson:"chapters" json:"chapters"`
}

type PrepPilotRepository interface {
	GetByID(id primitive.ObjectID) (*PrepPilot, error)
	GetByNotebookID(notebookID primitive.ObjectID) (*PrepPilot, error)
	Create(prepPilot *PrepPilot) error
	Update(prepPilot *PrepPilot) error
	Delete(id primitive.ObjectID) error
}
