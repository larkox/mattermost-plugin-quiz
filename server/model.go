package main

type QuizType string

const (
	QuizTypeSingleAnswer   QuizType = "single-answer"
	QuizTypeMultipleChoice QuizType = "multiple-choice"
)

type GameType string

const (
	GameTypeSolo  GameType = "solo"
	GameTypeParty GameType = "party"
)

type ScoringType string

const (
	ScoringTypeAll   ScoringType = "all"
	ScoringTypeFirst ScoringType = "first"
)

type ResourceType string

const (
	ResourceTypeText  ResourceType = "text"
	ResourceTypeLink  ResourceType = "link"
	ResourceTypeVideo ResourceType = "video"
	ResourceTypeQuiz  ResourceType = "quiz"
)

type Quiz struct {
	ID        string
	Name      string
	Type      QuizType
	Questions []Question
}

type Question struct {
	ID               string
	Question         string
	CorrectAnswer    string
	IncorrectAnswers []string
}

type Game struct {
	Quiz               Quiz
	GM                 string
	Score              map[string]int
	RemainingQuestions []Question
	RootPostID         string
	CurrentPostID      string
	Type               GameType
	ScoringType        ScoringType
	AlreadyAnswered    map[string]bool
	NQuestions         int
	CurrentAnswers     []string
	CorrectAnswer      int
	RightPlayers       []string
}

func (q Quiz) ValidQuestions() int {
	if q.Type == QuizTypeSingleAnswer {
		return len(q.Questions)
	}

	n := 0
	for _, question := range q.Questions {
		if len(question.IncorrectAnswers) >= IncorrectAnswerCount {
			n++
		}
	}

	return n
}

type Course struct {
	ID          string
	Name        string
	Description string
	Lessons     []*Lesson
}

type Lesson struct {
	Name         string
	Introduction string
	Resources    []*Resource
}

type Resource struct {
	Name    string
	Type    string
	Content string
	Pretext string
}
