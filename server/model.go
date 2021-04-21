package main

type QuizType string

const (
	QuizTypeSingleAnswer   QuizType = "single-answer"
	QuizTypeMultipleChoice QuizType = "multiple-choice"
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
