package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) CreateAttachmentFromQuiz(q *Quiz) []*model.SlackAttachment {
	attachment := model.SlackAttachment{
		Title:   "Quiz creation",
		Actions: []*model.PostAction{},
	}

	renameAction := model.PostAction{
		Type: "button",
		Name: "Name quiz",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathNameQuiz,
			Context: map[string]interface{}{
				AttachmentContextFieldID: q.ID,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &renameAction)

	if q.Name == "" {
		attachment.Text = "First name your quiz"
		return p.finishCreateAttachmentForQuiz(&attachment, q)
	}

	attachment.Text = "Quiz: " + q.Name
	renameAction.Name = "Rename quiz"

	changeTypeAction := model.PostAction{
		Type: "button",
		Name: "Select type",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathChangeType,
			Context: map[string]interface{}{
				AttachmentContextFieldID: q.ID,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &changeTypeAction)

	if q.Type == "" {
		attachment.Text += "\nSelect the quiz type"
		return p.finishCreateAttachmentForQuiz(&attachment, q)
	}

	attachment.Text += "\nType: " + string(q.Type)
	changeTypeAction.Name = "Change type"

	addQuestionAction := model.PostAction{
		Type: "button",
		Name: "Add question",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathAddQuestion,
			Context: map[string]interface{}{
				AttachmentContextFieldID: q.ID,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &addQuestionAction)

	validQuestions := q.ValidQuestions()
	allQuestions := len(q.Questions)

	if allQuestions > 0 {
		reviewQuestionsAction := model.PostAction{
			Type: "button",
			Name: "Review questions",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathReviewQuestions,
				Context: map[string]interface{}{
					AttachmentContextFieldID: q.ID,
				},
			},
		}
		attachment.Actions = append(attachment.Actions, &reviewQuestionsAction)

		removeQuestionsAction := model.PostAction{
			Type:  "button",
			Name:  "Remove questions",
			Style: "danger",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathRemoveQuestion,
				Context: map[string]interface{}{
					AttachmentContextFieldID: q.ID,
				},
			},
		}
		attachment.Actions = append(attachment.Actions, &removeQuestionsAction)
	}

	if allQuestions > 0 {
		saveAction := model.PostAction{
			Type:  "button",
			Name:  "Save quiz",
			Style: "good",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathSave,
				Context: map[string]interface{}{
					AttachmentContextFieldID: q.ID,
				},
			},
		}
		attachment.Actions = append(attachment.Actions, &saveAction)
	}

	attachment.Text += fmt.Sprintf("\nNumber of questions: %d", validQuestions)
	if validQuestions < len(q.Questions) {
		attachment.Text += fmt.Sprintf("\nWarning! %d out of %d questions are invalid", allQuestions-validQuestions, allQuestions)
	}

	return p.finishCreateAttachmentForQuiz(&attachment, q)
}

func (p *Plugin) finishCreateAttachmentForQuiz(attachment *model.SlackAttachment, q *Quiz) []*model.SlackAttachment {
	cancelAction := model.PostAction{
		Type:  "button",
		Name:  "Cancel",
		Style: "danger",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathDelete,
			Context: map[string]interface{}{
				AttachmentContextFieldID: q.ID,
			},
		},
	}

	attachment.Actions = append(attachment.Actions, &cancelAction)
	return []*model.SlackAttachment{attachment}
}

func (p *Plugin) GameAttachment(g *Game) []*model.SlackAttachment {
	currentQuestion := g.RemainingQuestions[0]
	attachment := &model.SlackAttachment{
		Title:   "Quiz: " + g.Quiz.Name,
		Text:    currentQuestion.Question,
		Footer:  fmt.Sprintf("Question %d out of %d.", g.NQuestions-len(g.RemainingQuestions)+1, g.NQuestions),
		Actions: []*model.PostAction{},
	}

	if g.Type == GameTypeParty {
		attachment.Text += fmt.Sprintf("\n\n%d people already answered.", len(g.AlreadyAnswered))
	}

	if g.Quiz.Type == QuizTypeMultipleChoice {
		for i, answer := range g.CurrentAnswers {
			attachment.Text += fmt.Sprintf("\n\nAnswer %d: %s", i+1, answer)
			attachment.Actions = append(attachment.Actions, &model.PostAction{
				Type: "button",
				Name: fmt.Sprintf("Answer %d", i+1),
				Integration: &model.PostActionIntegration{
					URL: p.getAttachmentURL() + AttachmentPathSelectAnswer,
					Context: map[string]interface{}{
						AttachmentContextFieldCorrect:    i == g.CorrectAnswer,
						AttachmentContextFieldGameID:     g.RootPostID,
						AttachmentContextFieldQuestionID: currentQuestion.ID,
					},
				},
			})
		}
	} else {
		attachment.Actions = append(attachment.Actions, &model.PostAction{
			Type: "button",
			Name: "Answer",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathAnswer,
				Context: map[string]interface{}{
					AttachmentContextFieldGameID:     g.RootPostID,
					AttachmentContextFieldQuestionID: currentQuestion.ID,
				},
			},
		})
	}

	attachment.Actions = append(attachment.Actions, &model.PostAction{
		Type: "button",
		Name: "Score",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathScore,
			Context: map[string]interface{}{
				AttachmentContextFieldGameID: g.RootPostID,
			},
		},
	})

	if g.Type == GameTypeParty {
		attachment.Actions = append(attachment.Actions, &model.PostAction{
			Type: "button",
			Name: "Next",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathNext,
				Context: map[string]interface{}{
					AttachmentContextFieldGameID:     g.RootPostID,
					AttachmentContextFieldQuestionID: currentQuestion.ID,
				},
			},
		})
	}

	return []*model.SlackAttachment{attachment}
}

func (p *Plugin) GameSolutionAttachment(g *Game) []*model.SlackAttachment {
	currentQuestion := g.RemainingQuestions[0]
	attachment := &model.SlackAttachment{
		Title:  "Quiz: " + g.Quiz.Name,
		Text:   currentQuestion.Question,
		Footer: fmt.Sprintf("Question %d out of %d.", g.NQuestions-len(g.RemainingQuestions)+1, g.NQuestions),
	}
	attachment.Text += fmt.Sprintf("\n\nThe correct answer was: %s", g.RemainingQuestions[0].CorrectAnswer)
	if g.Type == GameTypeParty {
		toAdd := "\n\nThe following users were right: "
		firstRun := true
		for _, name := range g.RightPlayers {
			if !firstRun {
				toAdd += ", "
			}
			toAdd += "@" + name
			firstRun = false
		}

		if len(g.RightPlayers) == 0 {
			toAdd = "\n\nNo users were right"
		}
		attachment.Text += toAdd
	}
	return []*model.SlackAttachment{attachment}
}

func (p *Plugin) GameEndAttachment(g *Game) []*model.SlackAttachment {
	attachment := &model.SlackAttachment{
		Title: "Quiz: " + g.Quiz.Name,
		Text:  "The quiz has finished\n\n" + getScores(g),
	}
	return []*model.SlackAttachment{attachment}
}

func getScores(g *Game) string {
	out := ""
	if g.Type == GameTypeParty {
		type row struct {
			name  string
			score int
		}
		rows := []row{}
		for name, score := range g.Score {
			rows = append(rows, row{name: name, score: score})
		}

		sort.Slice(rows, func(i, j int) bool { return rows[i].score > rows[j].score })
		out += "Scores:"
		for _, scoreRow := range rows {
			out += fmt.Sprintf("\n\n@%s: %d", scoreRow.name, scoreRow.score)
		}
		return out
	}

	out = "Your score: "
	score := 0
	for _, v := range g.Score {
		score = v
		break
	}
	out += strconv.Itoa(score)
	return out
}
