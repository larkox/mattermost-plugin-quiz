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

type scoreRow struct {
	name  string
	score int
}

func getScoreRows(g *Game) []scoreRow {
	rows := []scoreRow{}
	for name, score := range g.Score {
		rows = append(rows, scoreRow{name: name, score: score})
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].score > rows[j].score })
	return rows
}

func getScores(g *Game) string {
	out := ""
	if g.Type == GameTypeParty {
		rows := getScoreRows(g)
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

func (p *Plugin) CreateAttachmentFromCourse(c *Course) []*model.SlackAttachment {
	attachment := model.SlackAttachment{
		Title:   "Course creation",
		Actions: []*model.PostAction{},
	}

	renameAction := model.PostAction{
		Type: "button",
		Name: "Name course",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathNameCourse,
			Context: map[string]interface{}{
				AttachmentContextFieldID: c.ID,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &renameAction)

	if c.Name == "" {
		attachment.Text = "First name your course"
		return p.finishCreateAttachmentForCourse(&attachment, c)
	}

	attachment.Text = "Course: " + c.Name
	renameAction.Name = "Rename Course"

	changeTypeAction := model.PostAction{
		Type: "button",
		Name: "Add course description",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathCourseDescription,
			Context: map[string]interface{}{
				AttachmentContextFieldID: c.ID,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &changeTypeAction)

	if c.Description == "" {
		attachment.Text += "\nAdd a description to the course"
		return p.finishCreateAttachmentForCourse(&attachment, c)
	}

	attachment.Text += "\nDescription: " + string(c.Description)
	changeTypeAction.Name = "Change description"

	addQuestionAction := model.PostAction{
		Type: "button",
		Name: "Add lesson",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathAddLesson,
			Context: map[string]interface{}{
				AttachmentContextFieldID: c.ID,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &addQuestionAction)

	allLessons := len(c.Lessons)

	if allLessons > 0 {
		reviewQuestionsAction := model.PostAction{
			Type: "button",
			Name: "Edit lesson",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathEditLesson,
				Context: map[string]interface{}{
					AttachmentContextFieldID: c.ID,
				},
			},
		}
		attachment.Actions = append(attachment.Actions, &reviewQuestionsAction)

		saveAction := model.PostAction{
			Type:  "button",
			Name:  "Save course",
			Style: "good",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathSaveCourse,
				Context: map[string]interface{}{
					AttachmentContextFieldID: c.ID,
				},
			},
		}
		attachment.Actions = append(attachment.Actions, &saveAction)
	}

	attachment.Text += fmt.Sprintf("\nNumber of lessons: %d", allLessons)

	return p.finishCreateAttachmentForCourse(&attachment, c)
}

func (p *Plugin) CreateLessonAttachmentFromCourse(c *Course, index int) []*model.SlackAttachment {
	attachment := model.SlackAttachment{
		Title:   "Lesson creation",
		Actions: []*model.PostAction{},
	}
	attachment.Text = "Course: " + c.Name

	if index >= len(c.Lessons) {
		attachment.Text = "\nLesson not found."
		return p.finishCreateLessonAttachmentForCourse(&attachment, c, index)
	}

	lesson := c.Lessons[index]

	attachment.Text = "\nLesson name: " + lesson.Name
	attachment.Text = "\nLesson introduction: " + lesson.Introduction

	renameAction := model.PostAction{
		Type: "button",
		Name: "Rename lesson",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathNameLesson,
			Context: map[string]interface{}{
				AttachmentContextFieldID:          c.ID,
				AttachmentContextFieldLessonIndex: index,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &renameAction)

	changeIntroduction := model.PostAction{
		Type: "button",
		Name: "Change lesson introduction",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathLessonIntroduction,
			Context: map[string]interface{}{
				AttachmentContextFieldID:          c.ID,
				AttachmentContextFieldLessonIndex: index,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &changeIntroduction)

	addResourceAction := model.PostAction{
		Type: "button",
		Name: "Add resource",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathAddResource,
			Context: map[string]interface{}{
				AttachmentContextFieldID:          c.ID,
				AttachmentContextFieldLessonIndex: index,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &addResourceAction)

	addQuizResourceAction := model.PostAction{
		Type: "button",
		Name: "Add quiz resource",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathAddQuizResource,
			Context: map[string]interface{}{
				AttachmentContextFieldID:          c.ID,
				AttachmentContextFieldLessonIndex: index,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &addQuizResourceAction)

	allResources := len(lesson.Resources)

	if allResources > 0 {
		reviewQuestionsAction := model.PostAction{
			Type: "button",
			Name: "Remove resources",
			Integration: &model.PostActionIntegration{
				URL: p.getAttachmentURL() + AttachmentPathRemoveResources,
				Context: map[string]interface{}{
					AttachmentContextFieldID:          c.ID,
					AttachmentContextFieldLessonIndex: index,
				},
			},
		}
		attachment.Actions = append(attachment.Actions, &reviewQuestionsAction)
	}

	attachment.Text += fmt.Sprintf("\nNumber of resources: %d", allResources)

	return p.finishCreateLessonAttachmentForCourse(&attachment, c, index)
}

func (p *Plugin) finishCreateAttachmentForCourse(attachment *model.SlackAttachment, c *Course) []*model.SlackAttachment {
	cancelAction := model.PostAction{
		Type:  "button",
		Name:  "Cancel",
		Style: "danger",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathCourseDelete,
			Context: map[string]interface{}{
				AttachmentContextFieldID: c.ID,
			},
		},
	}

	attachment.Actions = append(attachment.Actions, &cancelAction)
	return []*model.SlackAttachment{attachment}
}

func (p *Plugin) finishCreateLessonAttachmentForCourse(attachment *model.SlackAttachment, c *Course, index int) []*model.SlackAttachment {
	cancelAction := model.PostAction{
		Type: "button",
		Name: "Back",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathLessonBack,
			Context: map[string]interface{}{
				AttachmentContextFieldID: c.ID,
			},
		},
	}

	attachment.Actions = append(attachment.Actions, &cancelAction)

	deleteAction := model.PostAction{
		Type:  "button",
		Name:  "Delete",
		Style: "danger",
		Integration: &model.PostActionIntegration{
			URL: p.getAttachmentURL() + AttachmentPathLessonDelete,
			Context: map[string]interface{}{
				AttachmentContextFieldID:          c.ID,
				AttachmentContextFieldLessonIndex: index,
			},
		},
	}
	attachment.Actions = append(attachment.Actions, &deleteAction)

	return []*model.SlackAttachment{attachment}
}
