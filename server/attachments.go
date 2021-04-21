package main

import (
	"fmt"

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
