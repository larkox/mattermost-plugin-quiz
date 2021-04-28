package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
)

// HTTPHandlerFuncWithUser is http.HandleFunc but userID is already exported
type HTTPHandlerFuncWithUser func(w http.ResponseWriter, r *http.Request, userID string)

// ResponseType indicates type of response returned by api
type ResponseType string

const (
	// ResponseTypeJSON indicates that response type is json
	ResponseTypeJSON ResponseType = "JSON_RESPONSE"
	// ResponseTypePlain indicates that response type is text plain
	ResponseTypePlain ResponseType = "TEXT_RESPONSE"
	// ResponseTypeDialog indicate that response is a DialogResponse object
	ResponseTypeDialog ResponseType = "DIALOG"
	// ResponseTypeAttachment indicates that response is a
	ResponseTypeAttachment ResponseType = "ATTACHMENT"
)

type APIErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
type Endpoint struct {
	Path    string
	Handler HTTPHandlerFuncWithUser
	Method  string
}

func (p *Plugin) initializeAPI() {
	p.router = mux.NewRouter()
	p.router.Use(p.withRecovery)

	dialogRouter := p.router.PathPrefix(DialogPath).Subrouter()
	attachmentRouter := p.router.PathPrefix(AttachmentPath).Subrouter()

	dialogRouterEndpoints := []Endpoint{
		{
			Path:    DialogPathNameQuiz,
			Handler: p.dialogNameQuiz,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathChangeType,
			Handler: p.dialogChangeType,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathDelete,
			Handler: p.dialogDelete,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathAddQuestion,
			Handler: p.dialogAddQuestion,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathReviewQuestions,
			Handler: p.dialogReviewQuestions,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathRemoveQuestion,
			Handler: p.dialogRemoveQuestions,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathGameStart,
			Handler: p.dialogGameStart,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathScore,
			Handler: p.dialogScore,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathAnswer,
			Handler: p.dialogAnswer,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathNameCourse,
			Handler: p.dialogNameCourse,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathCourseDescription,
			Handler: p.dialogCourseDescription,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathCourseDelete,
			Handler: p.dialogCourseDelete,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathAddLesson,
			Handler: p.dialogAddLesson,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathEditLesson,
			Handler: p.dialogEditLesson,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathNameLesson,
			Handler: p.dialogNameLesson,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathLessonIntroduction,
			Handler: p.dialogLessonIntroduction,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathAddResource,
			Handler: p.dialogAddResource,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathAddQuizResource,
			Handler: p.dialogAddQuizResource,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathRemoveResources,
			Handler: p.dialogRemoveResources,
			Method:  http.MethodPost,
		},
		{
			Path:    DialogPathLessonDelete,
			Handler: p.dialogLessonDelete,
			Method:  http.MethodPost,
		},
	}

	for _, e := range dialogRouterEndpoints {
		dialogRouter.HandleFunc(e.Path, p.extractUserMiddleWare(e.Handler, ResponseTypeDialog)).Methods(e.Method)
	}

	attachmentRouterEndpoints := []Endpoint{
		{
			Path:    AttachmentPathNameQuiz,
			Handler: p.attachmentNameQuiz,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathChangeType,
			Handler: p.attachmentChangeType,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathDelete,
			Handler: p.attachmentDelete,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathAddQuestion,
			Handler: p.attachmentAddQuestion,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathReviewQuestions,
			Handler: p.attachmentReviewQuestions,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathRemoveQuestion,
			Handler: p.attachmentRemoveQuestions,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathSave,
			Handler: p.attachmentSave,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathSelectAnswer,
			Handler: p.attachmentSelectAnswer,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathScore,
			Handler: p.attachmentScore,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathNext,
			Handler: p.attachmentNext,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathAnswer,
			Handler: p.attachmentAnswer,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathNameCourse,
			Handler: p.attachmentNameCourse,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathCourseDelete,
			Handler: p.attachmentCourseDelete,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathCourseDescription,
			Handler: p.attachmentCourseDescription,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathAddLesson,
			Handler: p.attachmentAddLesson,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathEditLesson,
			Handler: p.attachmentEditLesson,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathSaveCourse,
			Handler: p.attachmentSaveCourse,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathNameLesson,
			Handler: p.attachmentNameLesson,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathLessonBack,
			Handler: p.attachmentLessonBack,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathLessonIntroduction,
			Handler: p.attachmentLessonIntroduction,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathAddResource,
			Handler: p.attachmentAddResource,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathAddQuizResource,
			Handler: p.attachmentAddQuizResource,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathRemoveResources,
			Handler: p.attachmentRemoveResources,
			Method:  http.MethodPost,
		},
		{
			Path:    AttachmentPathLessonDelete,
			Handler: p.attachmentLessonDelete,
			Method:  http.MethodPost,
		},
	}

	for _, e := range attachmentRouterEndpoints {
		attachmentRouter.HandleFunc(e.Path, p.extractUserMiddleWare(e.Handler, ResponseTypeDialog)).Methods(e.Method)
	}

	p.router.PathPrefix(StaticPath).Handler(http.StripPrefix("/", http.FileServer(http.FS(staticAssets))))

	p.router.PathPrefix("/").HandlerFunc(p.defaultHandler)
}

func (p *Plugin) defaultHandler(w http.ResponseWriter, r *http.Request) {
	p.mm.Log.Debug("Unexpected call", "url", r.URL, "method", r.Method)
	w.WriteHeader(http.StatusNotFound)
}

func (p *Plugin) dialogNameQuiz(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	name, ok := req.Submission[DialogSubmissionFieldName].(string)
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldName: "Invalid name",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	q, err := p.store.GetQuiz(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if q == nil {
		dialogError(w, "quiz not found", nil)
		return
	}

	q.Name = name
	err = p.store.StoreQuiz(q)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromQuiz(q))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogChangeType(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	qType, ok := req.Submission[DialogSubmissionFieldType].(string)
	if !ok {
		errors := map[string]string{
			DialogSubmissionFieldType: "Could not get type",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	q, err := p.store.GetQuiz(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if q == nil {
		dialogError(w, "quiz not found", nil)
		return
	}

	q.Type = QuizType(qType)
	err = p.store.StoreQuiz(q)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromQuiz(q))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogDelete(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	err := p.mm.Post.DeletePost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	err = p.store.DeleteQuiz(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: req.ChannelId,
		Message:   "Quiz creation cancelled",
	}

	p.mm.Post.SendEphemeralPost(actingUserID, post)
	dialogOK(w)
}

func (p *Plugin) dialogAddQuestion(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	question, ok := req.Submission[DialogSubmissionFieldQuestion].(string)
	question = strings.TrimSpace(question)
	if !ok || question == "" {
		errors := map[string]string{
			DialogSubmissionFieldQuestion: "Could not get question",
		}
		dialogError(w, "Missing some value", errors)
		return
	}
	answer, ok := req.Submission[DialogSubmissionFieldAnswer].(string)
	answer = strings.TrimSpace(answer)
	if !ok || answer == "" {
		errors := map[string]string{
			DialogSubmissionFieldAnswer: "Could not get answer",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	q, err := p.store.GetQuiz(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if q == nil {
		dialogError(w, "quiz not found", nil)
		return
	}

	wrongAnsers := []string{}
	if q.Type == QuizTypeMultipleChoice {
		for i := 0; i < IncorrectAnswerCount; i++ {
			fieldName := DialogSubmissionFieldWrongAnswer + strconv.Itoa(i)
			wrongAnswer, ok := req.Submission[fieldName].(string)
			wrongAnswer = strings.TrimSpace(wrongAnswer)
			if !ok || wrongAnswer == "" {
				errors := map[string]string{
					fieldName: "Could not get answer",
				}
				dialogError(w, "Missing some value", errors)
				return
			}
			wrongAnsers = append(wrongAnsers, wrongAnswer)
		}
	}

	newQuestion := Question{
		ID:               model.NewId(),
		Question:         question,
		CorrectAnswer:    answer,
		IncorrectAnswers: wrongAnsers,
	}

	if q.Questions == nil {
		q.Questions = make([]Question, 0, 1)
	}
	q.Questions = append(q.Questions, newQuestion)

	err = p.store.StoreQuiz(q)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromQuiz(q))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogRemoveQuestions(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	q, err := p.store.GetQuiz(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if q == nil {
		dialogError(w, "quiz not found", nil)
		return
	}

	for toDeleteID, value := range req.Submission {
		v, ok := value.(bool)
		if !ok || !v {
			continue
		}

		for i, question := range q.Questions {
			if question.ID == toDeleteID {
				q.Questions = append(q.Questions[:i], q.Questions[i+1:]...)
				break
			}
		}
	}

	err = p.store.StoreQuiz(q)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromQuiz(q))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	dialogOK(w)
}

func (p *Plugin) dialogReviewQuestions(w http.ResponseWriter, r *http.Request, actingUserID string) {
	dialogOK(w)
}

func (p *Plugin) dialogGameStart(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)

	quizID, ok := req.Submission[DialogSubmissionFieldGameQuiz].(string)
	quizID = strings.TrimSpace(quizID)
	if !ok || quizID == "" {
		errors := map[string]string{
			DialogSubmissionFieldGameQuiz: "Could not get quiz",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	gameType, ok := req.Submission[DialogSubmissionFieldGameType].(string)
	gameType = strings.TrimSpace(gameType)
	if !ok || gameType == "" {
		errors := map[string]string{
			DialogSubmissionFieldGameType: "Could not get type",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	if gameType != string(GameTypeSolo) && gameType != string(GameTypeParty) {
		errors := map[string]string{
			DialogSubmissionFieldGameType: "Type not recognized",
		}
		dialogError(w, "Unrecognized value", errors)
		return
	}

	scoring, ok := req.Submission[DialogSubmissionFieldGameScoring].(string)
	scoring = strings.TrimSpace(scoring)
	if !ok || scoring == "" {
		errors := map[string]string{
			DialogSubmissionFieldGameScoring: "Could not get scoring",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	if scoring != string(ScoringTypeAll) && scoring != string(ScoringTypeFirst) {
		errors := map[string]string{
			DialogSubmissionFieldGameScoring: "Scoring type not recognized",
		}
		dialogError(w, "Unrecognized value", errors)
		return
	}

	nQuestionsFloat, ok := req.Submission[DialogSubmissionFieldNumberOfQuestions].(float64)
	if !ok {
		errors := map[string]string{
			DialogSubmissionFieldNumberOfQuestions: "Could not get the number of questions",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	nQuestions := int(nQuestionsFloat)

	quiz, err := p.store.GetQuiz(quizID)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	validQuestions := quiz.ValidQuestions()

	if nQuestions <= 0 {
		nQuestions = validQuestions
	}

	if nQuestions > validQuestions {
		nQuestions = validQuestions
	}

	questions := quiz.Questions
	if quiz.Type == QuizTypeMultipleChoice {
		questions = []Question{}
		for _, question := range quiz.Questions {
			if len(question.IncorrectAnswers) >= IncorrectAnswerCount {
				questions = append(questions, question)
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

	game := &Game{
		Quiz:               *quiz,
		GM:                 actingUserID,
		Score:              map[string]int{},
		Type:               GameType(gameType),
		ScoringType:        ScoringType(scoring),
		RemainingQuestions: questions[:nQuestions],
		NQuestions:         nQuestions,
		AlreadyAnswered:    map[string]bool{},
	}

	if quiz.Type == QuizTypeMultipleChoice {
		game.CurrentAnswers, game.CorrectAnswer = getRandomAnswers(game.RemainingQuestions[0])
	}

	post := &model.Post{
		Message: "New quiz",
	}

	model.ParseSlackAttachment(post, p.GameAttachment(game))
	if gameType == string(GameTypeSolo) {
		err = p.mm.Post.DM(p.BotUserID, actingUserID, post)
		if err != nil {
			dialogError(w, err.Error(), nil)
			return
		}
	} else {
		post.ChannelId = req.ChannelId
		post.UserId = p.BotUserID
		err = p.mm.Post.CreatePost(post)
		if err != nil {
			dialogError(w, err.Error(), nil)
			return
		}
	}

	game.RootPostID = post.Id
	game.CurrentPostID = post.Id

	err = p.store.StoreGame(game)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogScore(w http.ResponseWriter, r *http.Request, actingUserID string) {
	dialogOK(w)
}

func (p *Plugin) dialogAnswer(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	state := strings.Split(req.State, ",")
	if len(state) != 2 {
		dialogError(w, "wrong state", nil)
		return
	}
	id := state[0]
	qID := state[1]

	answer, ok := req.Submission[DialogSubmissionFieldGameAnswer].(string)
	answer = strings.TrimSpace(answer)
	if !ok || answer == "" {
		errors := map[string]string{
			DialogSubmissionFieldGameAnswer: "Could not get the answer",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	g, err := p.store.GetGame(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	if g == nil {
		dialogError(w, "game not found", nil)
		return
	}

	if qID != g.RemainingQuestions[0].ID {
		dialogError(w, "this question has been already passed", nil)
		return
	}

	user, err := p.mm.User.Get(actingUserID)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	if g.AlreadyAnswered[user.Username] {
		dialogError(w, "you already tried to answer this question", nil)
		return
	}

	g.AlreadyAnswered[user.Username] = true

	responseMessage := "Your answer is incorrect."
	if answer == g.RemainingQuestions[0].CorrectAnswer {
		responseMessage = "You are correct!"
		g.Score[user.Username] += 1
		if g.ScoringType == ScoringTypeFirst && len(g.AlreadyAnswered) == 1 {
			g.Score[user.Username] += 2
		}
		g.RightPlayers = append(g.RightPlayers, user.Username)
	}

	responsePost := &model.Post{
		UserId:    p.BotUserID,
		Message:   responseMessage,
		ChannelId: req.ChannelId,
	}

	post, err := p.mm.Post.GetPost(g.CurrentPostID)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	if g.Type == GameTypeParty {
		model.ParseSlackAttachment(post, p.GameAttachment(g))
		err = p.mm.Post.UpdatePost(post)
		if err != nil {
			dialogError(w, err.Error(), nil)
			return
		}
		err = p.store.StoreGame(g)
		if err != nil {
			dialogError(w, err.Error(), nil)
			return
		}
		p.mm.Post.SendEphemeralPost(actingUserID, responsePost)
		dialogOK(w)
		return
	}

	model.ParseSlackAttachment(post, p.GameSolutionAttachment(g))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	err = p.handleNextQuestion(g, req.ChannelId, actingUserID)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	p.mm.Post.SendEphemeralPost(actingUserID, responsePost)
	dialogOK(w)
}

func (p *Plugin) dialogNameCourse(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	name, ok := req.Submission[DialogSubmissionFieldName].(string)
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldName: "Invalid name",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	c.Name = name
	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromCourse(c))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogCourseDescription(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	description, ok := req.Submission[DialogSubmissionFieldDescription].(string)
	description = strings.TrimSpace(description)
	if !ok || description == "" {
		errors := map[string]string{
			DialogSubmissionFieldDescription: "Invalid name",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	c.Description = description
	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromCourse(c))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogCourseDelete(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	err := p.mm.Post.DeletePost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	err = p.store.DeleteCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: req.ChannelId,
		Message:   "Course creation cancelled",
	}

	p.mm.Post.SendEphemeralPost(actingUserID, post)
	dialogOK(w)
}

func (p *Plugin) dialogAddLesson(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	name, ok := req.Submission[DialogSubmissionFieldName].(string)
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldName: "Invalid name",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	description, ok := req.Submission[DialogSubmissionFieldDescription].(string)
	description = strings.TrimSpace(description)
	if !ok || description == "" {
		errors := map[string]string{
			DialogSubmissionFieldDescription: "Invalid name",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	c.Lessons = append(c.Lessons, &Lesson{
		Name:         name,
		Introduction: description,
	})
	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateLessonAttachmentFromCourse(c, len(c.Lessons)-1))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogEditLesson(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id := req.State
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	indexStr, ok := req.Submission[DialogSubmissionFieldLesson].(string)
	indexStr = strings.TrimSpace(indexStr)
	if !ok || indexStr == "" {
		errors := map[string]string{
			DialogSubmissionFieldLesson: "Invalid lesson",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		errors := map[string]string{
			DialogSubmissionFieldLesson: "Invalid lesson",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		errors := map[string]string{
			DialogSubmissionFieldLesson: "Lesson not found",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	model.ParseSlackAttachment(post, p.CreateLessonAttachmentFromCourse(c, index))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogNameLesson(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id, index := getLessonIndexAndIDFromState(req.State)
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	name, ok := req.Submission[DialogSubmissionFieldName].(string)
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldName: "Invalid name",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		dialogError(w, "Cannot find this lesson. Please hit the back button.", nil)
		return
	}

	c.Lessons[index].Name = name
	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateLessonAttachmentFromCourse(c, index))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogLessonIntroduction(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id, index := getLessonIndexAndIDFromState(req.State)
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	introduction, ok := req.Submission[DialogSubmissionFieldDescription].(string)
	introduction = strings.TrimSpace(introduction)
	if !ok || introduction == "" {
		errors := map[string]string{
			DialogSubmissionFieldDescription: "Invalid introduction",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		dialogError(w, "Cannot find this lesson. Please hit the back button.", nil)
		return
	}

	c.Lessons[index].Introduction = introduction
	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateLessonAttachmentFromCourse(c, index))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogAddResource(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id, index := getLessonIndexAndIDFromState(req.State)
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	name, ok := req.Submission[DialogSubmissionFieldName].(string)
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldName: "Invalid introduction",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	pretext, ok := req.Submission[DialogSubmissionFieldDescription].(string)
	pretext = strings.TrimSpace(pretext)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldDescription: "Invalid pretext",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	resourceType, ok := req.Submission[DialogSubmissionFieldType].(string)
	resourceType = strings.TrimSpace(resourceType)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldType: "Invalid resource type",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	content, ok := req.Submission[DialogSubmissionFieldContent].(string)
	content = strings.TrimSpace(content)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldContent: "Invalid content",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		dialogError(w, "Cannot find this lesson. Please hit the back button.", nil)
		return
	}

	lesson := c.Lessons[index]

	lesson.Resources = append(lesson.Resources, &Resource{
		Name:    name,
		Type:    resourceType,
		Content: content,
		Pretext: pretext,
	})

	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateLessonAttachmentFromCourse(c, index))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogAddQuizResource(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id, index := getLessonIndexAndIDFromState(req.State)
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	name, ok := req.Submission[DialogSubmissionFieldName].(string)
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldName: "Invalid introduction",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	pretext, ok := req.Submission[DialogSubmissionFieldDescription].(string)
	pretext = strings.TrimSpace(pretext)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldDescription: "Invalid pretext",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	quizID, ok := req.Submission[DialogSubmissionFieldQuiz].(string)
	quizID = strings.TrimSpace(quizID)
	if !ok || name == "" {
		errors := map[string]string{
			DialogSubmissionFieldQuiz: "Invalid quiz",
		}
		dialogError(w, "Missing some value", errors)
		return
	}

	quiz, err := p.store.GetQuiz(quizID)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if quiz == nil {
		errors := map[string]string{
			DialogSubmissionFieldQuiz: "Quiz not found",
		}
		dialogError(w, "quiz not found", errors)
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		dialogError(w, "Cannot find this lesson. Please hit the back button.", nil)
		return
	}

	lesson := c.Lessons[index]

	lesson.Resources = append(lesson.Resources, &Resource{
		Name:    name,
		Type:    string(ResourceTypeQuiz),
		Content: quizID,
		Pretext: pretext,
	})

	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateLessonAttachmentFromCourse(c, index))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	dialogOK(w)
}

func (p *Plugin) dialogRemoveResources(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id, index := getLessonIndexAndIDFromState(req.State)
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		dialogError(w, "Cannot find this lesson. Please hit the back button.", nil)
		return
	}

	lesson := c.Lessons[index]

	deleteIndexes := []int{}
	for toDeleteID, value := range req.Submission {
		v, ok := value.(bool)
		if !ok || !v {
			continue
		}

		rIndex, err := strconv.Atoi(toDeleteID)
		if err != nil {
			p.mm.Log.Debug("Index not convertible to int", "err", err)
			continue
		}

		deleteIndexes = append(deleteIndexes, rIndex)
	}
	sort.Slice(deleteIndexes, func(i, j int) bool { return i > j })

	for _, i := range deleteIndexes {
		lesson.Resources = append(lesson.Resources[:i], lesson.Resources[i+1:]...)
	}

	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateLessonAttachmentFromCourse(c, index))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	dialogOK(w)
}

func (p *Plugin) dialogLessonDelete(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.SubmitDialogRequestFromJson(r.Body)
	id, index := getLessonIndexAndIDFromState(req.State)
	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	c, err := p.store.GetCourse(id)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}
	if c == nil {
		dialogError(w, "course not found", nil)
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		dialogError(w, "Cannot find this lesson. Please hit the back button.", nil)
		return
	}

	c.Lessons = append(c.Lessons[:index], c.Lessons[index+1:]...)

	err = p.store.StoreCourse(c)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromCourse(c))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		dialogError(w, err.Error(), nil)
		return
	}

	dialogOK(w)
}

func (p *Plugin) attachmentNameQuiz(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)

	q, err := p.store.GetQuiz(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	defaultName := ""
	if q != nil {
		defaultName = q.Name
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathNameQuiz,
		Dialog: model.Dialog{
			Title:            "Name quiz",
			IntroductionText: "Write the name of your Quiz",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Quiz name",
					Name:        DialogSubmissionFieldName,
					Type:        DialogTypeText,
					Default:     defaultName,
				},
			},
			State: id,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentChangeType(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)

	q, err := p.store.GetQuiz(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if q == nil {
		attachmentError(w, "quiz not found")
		return
	}

	dr := model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathChangeType,
		Dialog: model.Dialog{
			Title:            "Select type",
			IntroductionText: "Select the type of quiz you want to create",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Quiz type",
					Name:        DialogSubmissionFieldType,
					Type:        DialogTypeSelect,
					Default:     string(q.Type),
					Options: []*model.PostActionOptions{
						{
							Text:  "Single answer",
							Value: string(QuizTypeSingleAnswer),
						},
						{
							Text:  "Multiple choice",
							Value: string(QuizTypeMultipleChoice),
						},
					},
				},
			},
			State: id,
		},
	}

	if q.Type != "" {
		dr.Dialog.IntroductionText += "\n\nWARNING: Changing the type may render some questions invalid. The questions will not be lost, but will not appear in the quiz."
	}
	err = p.mm.Frontend.OpenInteractiveDialog(dr)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentDelete(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)
	err := p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathDelete,
		Dialog: model.Dialog{
			Title:            "Cancel creation",
			IntroductionText: "Are you sure you want to cancel the creation of this quiz? All changes will be lost.",
			SubmitLabel:      "Delete",
			State:            id,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentAddQuestion(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)

	q, err := p.store.GetQuiz(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if q == nil {
		attachmentError(w, "quiz not found")
		return
	}

	wrongAnswerElements := []model.DialogElement{}
	if q.Type == QuizTypeMultipleChoice {
		for i := 0; i < IncorrectAnswerCount; i++ {
			e := model.DialogElement{
				DisplayName: "Incorrect Answer",
				Name:        DialogSubmissionFieldWrongAnswer + strconv.Itoa(i),
				Type:        DialogTypeText,
			}
			wrongAnswerElements = append(wrongAnswerElements, e)
		}
	}

	dr := model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathAddQuestion,
		Dialog: model.Dialog{
			Title:            "Add Question",
			IntroductionText: "Write the question to add",
			SubmitLabel:      "Add",
			Elements: []model.DialogElement{
				{
					DisplayName: "Question",
					Name:        DialogSubmissionFieldQuestion,
					Type:        DialogTypeText,
				},
				{
					DisplayName: "Answer",
					Name:        DialogSubmissionFieldAnswer,
					Type:        DialogTypeText,
				},
			},
			State: id,
		},
	}

	dr.Dialog.Elements = append(dr.Dialog.Elements, wrongAnswerElements...)
	err = p.mm.Frontend.OpenInteractiveDialog(dr)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentReviewQuestions(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)

	q, err := p.store.GetQuiz(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if q == nil {
		attachmentError(w, "quiz not found")
		return
	}

	questions := ""
	for _, question := range q.Questions {
		questions += "\n\n" + Separator
		if q.Type == QuizTypeMultipleChoice && len(question.IncorrectAnswers) < IncorrectAnswerCount {
			questions += "\nWARNING: Invalid question\n"
		}
		questions += fmt.Sprintf("\nQuestion: %s\n\nCorrect Answer: %s", question.Question, question.CorrectAnswer)
		if q.Type == QuizTypeMultipleChoice {
			questions += "\n\nIncorrect Answers:"
			for _, answer := range question.IncorrectAnswers {
				questions += "\n\n" + answer
			}
		}
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathReviewQuestions,
		Dialog: model.Dialog{
			Title:            "Review questions",
			IntroductionText: "These are all the questions registered:\n" + questions,
			SubmitLabel:      "OK",
			State:            id,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentRemoveQuestions(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)

	q, err := p.store.GetQuiz(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if q == nil {
		attachmentError(w, "quiz not found")
		return
	}

	elements := []model.DialogElement{}
	for _, question := range q.Questions {
		element := model.DialogElement{
			DisplayName: question.Question,
			Name:        question.ID,
			Type:        DialogTypeBool,
			Optional:    true,
		}
		text := ""
		if q.Type == QuizTypeMultipleChoice && len(question.IncorrectAnswers) < IncorrectAnswerCount {
			text += "WARNING: Invalid question; "
		}
		text += fmt.Sprintf("Correct Answer: %s", question.CorrectAnswer)
		if q.Type == QuizTypeMultipleChoice {
			text += "; Incorrect Answers:"
			firstRun := true
			for _, answer := range question.IncorrectAnswers {
				if !firstRun {
					text += ";"
				}
				text += " " + answer
				firstRun = false
			}
		}
		element.HelpText = text
		elements = append(elements, element)
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathRemoveQuestion,
		Dialog: model.Dialog{
			Title:            "Remove questions",
			IntroductionText: "Select the questions to delete.",
			SubmitLabel:      "Remove selected",
			State:            id,
			Elements:         elements,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentSave(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)

	q, err := p.store.GetQuiz(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if q == nil {
		attachmentError(w, "quiz not found")
		return
	}

	if q.ValidQuestions() == 0 {
		attachmentError(w, "cannot save a quiz with no valid questions")
		return
	}

	err = p.store.AddAvailableQuiz(q)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	p.GrantBadge(AchievementNameContentCreator, actingUserID)

	resp := model.PostActionIntegrationResponse{
		Update: &model.Post{
			Message: fmt.Sprintf("Quiz `%s` saved and ready to use.", q.Name),
			Props:   model.StringInterface{},
		},
	}
	_, _ = w.Write(resp.ToJson())
}

func (p *Plugin) attachmentAnswer(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getGameIDFromPostActionRequest(req)

	g, err := p.store.GetGame(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if g == nil {
		attachmentError(w, "game not found")
		return
	}

	qID, ok := req.Context[AttachmentContextFieldQuestionID].(string)
	if !ok || id == "" {
		attachmentError(w, "cannot find question ID")
		return
	}

	if qID != g.RemainingQuestions[0].ID {
		attachmentError(w, "mismatch between post and internal state")
		return
	}

	user, err := p.mm.User.Get(actingUserID)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if g.AlreadyAnswered[user.Username] {
		attachmentError(w, "you already tried to answer this question")
		return
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathAnswer,
		Dialog: model.Dialog{
			Title:            "Answer",
			IntroductionText: g.RemainingQuestions[0].Question,
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Your answer",
					Name:        DialogSubmissionFieldGameAnswer,
					Type:        DialogTypeText,
				},
			},
			State: strings.Join([]string{id, qID}, ","),
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentSelectAnswer(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getGameIDFromPostActionRequest(req)

	g, err := p.store.GetGame(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if g == nil {
		attachmentError(w, "game not found")
		return
	}

	qID, ok := req.Context[AttachmentContextFieldQuestionID].(string)
	if !ok || id == "" {
		attachmentError(w, "cannot find question ID")
		return
	}

	if qID != g.RemainingQuestions[0].ID {
		attachmentError(w, "mismatch between post and internal state")
		return
	}

	user, err := p.mm.User.Get(actingUserID)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if g.AlreadyAnswered[user.Username] {
		attachmentError(w, "you already tried to answer this question")
		return
	}

	g.AlreadyAnswered[user.Username] = true

	correctAnswer, ok := req.Context[AttachmentContextFieldCorrect].(bool)
	if !ok {
		attachmentError(w, "cannot find whether is the correct answer")
		return
	}

	responseMessage := "Your answer is incorrect."
	if correctAnswer {

		responseMessage = "You are correct!"
		g.Score[user.Username] += 1
		if g.ScoringType == ScoringTypeFirst && len(g.AlreadyAnswered) == 1 {
			g.Score[user.Username] += 2
		}
		g.RightPlayers = append(g.RightPlayers, user.Username)
	}

	post, err := p.mm.Post.GetPost(req.PostId)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if g.Type == GameTypeParty {
		model.ParseSlackAttachment(post, p.GameAttachment(g))
		err = p.mm.Post.UpdatePost(post)
		if err != nil {
			attachmentError(w, err.Error())
			return
		}
		err = p.store.StoreGame(g)
		if err != nil {
			attachmentError(w, err.Error())
			return
		}
		attachmentOK(w, responseMessage)
		return
	}

	model.ParseSlackAttachment(post, p.GameSolutionAttachment(g))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	err = p.handleNextQuestion(g, req.ChannelId, actingUserID)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	attachmentOK(w, responseMessage)
}

func (p *Plugin) attachmentScore(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getGameIDFromPostActionRequest(req)

	g, err := p.store.GetGame(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if g == nil {
		attachmentError(w, "game not found")
		return
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathScore,
		Dialog: model.Dialog{
			Title:            "Score",
			IntroductionText: getScores(g),
			SubmitLabel:      "OK",
			State:            id,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentNext(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getGameIDFromPostActionRequest(req)

	g, err := p.store.GetGame(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if g == nil {
		attachmentError(w, "game not found")
		return
	}

	if g.GM != actingUserID {
		attachmentError(w, "only the person who created the quiz can pass to the next question")
		return
	}

	qID, ok := req.Context[AttachmentContextFieldQuestionID].(string)
	if !ok || id == "" {
		attachmentError(w, "cannot find question ID")
		return
	}

	if qID != g.RemainingQuestions[0].ID {
		attachmentError(w, "mismatch between post and internal state")
		return
	}

	post, err := p.mm.Post.GetPost(req.PostId)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	model.ParseSlackAttachment(post, p.GameSolutionAttachment(g))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	err = p.handleNextQuestion(g, req.ChannelId, actingUserID)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	attachmentOK(w, "")
}

func (p *Plugin) handleNextQuestion(g *Game, channelID, actingUserID string) error {
	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: channelID,
		RootId:    g.RootPostID,
		Message:   "Next question!",
	}

	g.RemainingQuestions = g.RemainingQuestions[1:]

	if len(g.RemainingQuestions) == 0 {
		post.Message = "Quiz finished!"
		model.ParseSlackAttachment(post, p.GameEndAttachment(g))
		err := p.mm.Post.CreatePost(post)
		if err != nil {
			return err
		}

		if g.Type == GameTypeSolo {
			p.GrantBadge(AchievementNameHardWorker, actingUserID)
		}

		if g.Type == GameTypeParty {
			rows := getScoreRows(g)
			if len(rows) > 0 {
				user, err := p.mm.User.GetByUsername(rows[0].name)
				if err == nil {
					p.GrantBadge(AchievementNameWinner, user.Id)
				}
			}
		}
		return p.store.DeleteGame(g.RootPostID)

	}

	g.CurrentAnswers, g.CorrectAnswer = getRandomAnswers(g.RemainingQuestions[0])
	g.AlreadyAnswered = map[string]bool{}
	g.RightPlayers = []string{}

	model.ParseSlackAttachment(post, p.GameAttachment(g))
	err := p.mm.Post.CreatePost(post)
	if err != nil {
		return err
	}

	g.CurrentPostID = post.Id

	return p.store.StoreGame(g)
}

func (p *Plugin) attachmentNameCourse(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	defaultName := ""
	if c != nil {
		defaultName = c.Name
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathNameCourse,
		Dialog: model.Dialog{
			Title:            "Name course",
			IntroductionText: "Write the name of your course",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Course name",
					Name:        DialogSubmissionFieldName,
					Type:        DialogTypeText,
					Default:     defaultName,
				},
			},
			State: id,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentCourseDescription(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	defaultDescription := ""
	if c != nil {
		defaultDescription = c.Description
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathCourseDescription,
		Dialog: model.Dialog{
			Title:            "Add course description",
			IntroductionText: "Write the description of your course",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Course description",
					Name:        DialogSubmissionFieldDescription,
					Type:        DialogTypeText,
					Default:     defaultDescription,
				},
			},
			State: id,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentCourseDelete(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)
	err := p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathCourseDelete,
		Dialog: model.Dialog{
			Title:            "Cancel creation",
			IntroductionText: "Are you sure you want to cancel the creation of this course? All changes will be lost.",
			SubmitLabel:      "Delete",
			State:            id,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentSaveCourse(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if c == nil {
		attachmentError(w, "course not found")
		return
	}

	if len(c.Lessons) == 0 {
		attachmentError(w, "cannot save a course with no lessons")
		return
	}

	err = p.store.AddAvailableCourse(c)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	resp := model.PostActionIntegrationResponse{
		Update: &model.Post{
			Message: fmt.Sprintf("Course `%s` saved and ready to use.", c.Name),
			Props:   model.StringInterface{},
		},
	}
	_, _ = w.Write(resp.ToJson())
}

func (p *Plugin) attachmentAddLesson(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if c == nil {
		attachmentError(w, "course not found")
		return
	}

	dr := model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathAddLesson,
		Dialog: model.Dialog{
			Title:            "Add a lesson",
			IntroductionText: "Write the lesson basic information to add",
			SubmitLabel:      "Add",
			Elements: []model.DialogElement{
				{
					DisplayName: "Name",
					Name:        DialogSubmissionFieldName,
					Type:        DialogTypeText,
				},
				{
					DisplayName: "Description",
					Name:        DialogSubmissionFieldDescription,
					Type:        DialogTypeText,
				},
			},
			State: id,
		},
	}

	err = p.mm.Frontend.OpenInteractiveDialog(dr)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentEditLesson(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if c == nil {
		attachmentError(w, "quiz not found")
		return
	}

	element := model.DialogElement{
		DisplayName: "Lesson",
		Name:        DialogSubmissionFieldLesson,
		Type:        "select",
		Options:     []*model.PostActionOptions{},
	}
	for i, lesson := range c.Lessons {
		option := &model.PostActionOptions{
			Text:  lesson.Name,
			Value: strconv.Itoa(i),
		}
		element.Options = append(element.Options, option)
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathEditLesson,
		Dialog: model.Dialog{
			Title:            "Edit lessons",
			IntroductionText: "Select the lesson to edit.",
			SubmitLabel:      "Edit",
			State:            id,
			Elements:         []model.DialogElement{element},
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentNameLesson(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)
	index := getLessonIndexFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		attachmentError(w, "Cannot find this lesson. Please hit the back button.")
		return
	}

	lesson := c.Lessons[index]

	defaultName := ""
	if c != nil {
		defaultName = lesson.Name
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathNameLesson,
		Dialog: model.Dialog{
			Title:            "Name lesson",
			IntroductionText: "Write the name of your lesson",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Lesson name",
					Name:        DialogSubmissionFieldName,
					Type:        DialogTypeText,
					Default:     defaultName,
				},
			},
			State: getLessonDialogState(id, index),
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentLessonIntroduction(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)
	index := getLessonIndexFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		attachmentError(w, "Cannot find this lesson. Please hit the back button.")
		return
	}

	lesson := c.Lessons[index]

	defaultIntroduction := ""
	if c != nil {
		defaultIntroduction = lesson.Introduction
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathLessonIntroduction,
		Dialog: model.Dialog{
			Title:            "Name lesson",
			IntroductionText: "Write the introduction of your lesson",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Lesson introduction",
					Name:        DialogSubmissionFieldDescription,
					Type:        DialogTypeText,
					Default:     defaultIntroduction,
				},
			},
			State: getLessonDialogState(id, index),
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentAddResource(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)
	index := getLessonIndexFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		attachmentError(w, "Cannot find this lesson. Please hit the back button.")
		return
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathAddResource,
		Dialog: model.Dialog{
			Title:            "Add resource",
			IntroductionText: "Fill the resource information",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Resource name",
					Name:        DialogSubmissionFieldName,
					Type:        DialogTypeText,
				},
				{
					DisplayName: "Resource pretext",
					Name:        DialogSubmissionFieldDescription,
					Type:        DialogTypeText,
				},
				{
					DisplayName: "Resource type",
					Name:        DialogSubmissionFieldType,
					Type:        DialogTypeSelect,
					Options: []*model.PostActionOptions{
						{
							Text:  "Video",
							Value: string(ResourceTypeVideo),
						},
						{
							Text:  "Text",
							Value: string(ResourceTypeText),
						},
						{
							Text:  "Link",
							Value: string(ResourceTypeLink),
						},
					},
				},
				{
					DisplayName: "Resource content",
					Name:        DialogSubmissionFieldContent,
					Type:        DialogTypeText,
				},
			},
			State: getLessonDialogState(id, index),
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentAddQuizResource(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)
	index := getLessonIndexFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		attachmentError(w, "Cannot find this lesson. Please hit the back button.")
		return
	}

	quizOptions := []*model.PostActionOptions{}
	qq := p.store.GetAvailableQuizes()

	if len(qq) == 0 {
		attachmentError(w, "No quizzes available to add.")
		return
	}

	for _, q := range qq {
		quizOptions = append(quizOptions, &model.PostActionOptions{Text: q.Name, Value: q.ID})
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathAddQuizResource,
		Dialog: model.Dialog{
			Title:            "Add resource",
			IntroductionText: "Fill the resource information",
			SubmitLabel:      "Submit",
			Elements: []model.DialogElement{
				{
					DisplayName: "Resource name",
					Name:        DialogSubmissionFieldName,
					Type:        DialogTypeText,
				},
				{
					DisplayName: "Resource pretext",
					Name:        DialogSubmissionFieldDescription,
					Type:        DialogTypeText,
				},
				{
					DisplayName: "Quiz",
					Name:        DialogSubmissionFieldQuiz,
					Type:        DialogTypeSelect,
					Options:     quizOptions,
				},
			},
			State: getLessonDialogState(id, index),
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentRemoveResources(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getQuizIDFromPostActionRequest(req)
	index := getLessonIndexFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	if c == nil {
		attachmentError(w, "course not found")
		return
	}

	if index < 0 || index >= len(c.Lessons) {
		attachmentError(w, "Cannot find this lesson. Please hit the back button.")
		return
	}

	lesson := c.Lessons[index]

	elements := []model.DialogElement{}
	for i, resource := range lesson.Resources {
		element := model.DialogElement{
			DisplayName: resource.Name,
			Name:        strconv.Itoa(i),
			Type:        DialogTypeBool,
			Optional:    true,
		}
		elements = append(elements, element)
	}

	err = p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathRemoveResources,
		Dialog: model.Dialog{
			Title:            "Remove resources",
			IntroductionText: "Select the resources to delete.",
			SubmitLabel:      "Remove selected",
			State:            getLessonDialogState(id, index),
			Elements:         elements,
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentLessonBack(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)

	c, err := p.store.GetCourse(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	post, err := p.mm.Post.GetPost(id)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}

	model.ParseSlackAttachment(post, p.CreateAttachmentFromCourse(c))
	err = p.mm.Post.UpdatePost(post)
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func (p *Plugin) attachmentLessonDelete(w http.ResponseWriter, r *http.Request, actingUserID string) {
	req := model.PostActionIntegrationRequestFromJson(r.Body)
	id := getCourseIDFromPostActionRequest(req)
	index := getLessonIndexFromPostActionRequest(req)
	err := p.mm.Frontend.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: req.TriggerId,
		URL:       p.getDialogURL() + DialogPathLessonDelete,
		Dialog: model.Dialog{
			Title:            "Delete lesson",
			IntroductionText: "Are you sure you want to delete this lesson? All changes will be lost.",
			SubmitLabel:      "Delete",
			State:            getLessonDialogState(id, index),
		},
	})
	if err != nil {
		attachmentError(w, err.Error())
		return
	}
	attachmentOK(w, "")
}

func getQuizIDFromPostActionRequest(req *model.PostActionIntegrationRequest) string {
	id, ok := req.Context[AttachmentContextFieldID].(string)
	if !ok || id == "" {
		id = req.PostId
	}
	return id
}

func getCourseIDFromPostActionRequest(req *model.PostActionIntegrationRequest) string {
	id, ok := req.Context[AttachmentContextFieldID].(string)
	if !ok || id == "" {
		id = req.PostId
	}
	return id
}

func getLessonIndexFromPostActionRequest(req *model.PostActionIntegrationRequest) int {
	id, ok := req.Context[AttachmentContextFieldLessonIndex].(float64)
	if !ok {
		id = -1
	}
	return int(id)
}

func getGameIDFromPostActionRequest(req *model.PostActionIntegrationRequest) string {
	id, ok := req.Context[AttachmentContextFieldGameID].(string)
	if !ok || id == "" {
		id = req.PostId
	}
	return id
}

func getLessonDialogState(cID string, index int) string {
	return fmt.Sprintf("%s,%d", cID, index)
}

func getLessonIndexAndIDFromState(state string) (string, int) {
	parts := strings.Split(state, ",")
	if len(parts) != 2 {
		return "", -1
	}

	cID := parts[0]
	indexStr := parts[1]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", -1
	}

	return cID, index
}

// func (p *Plugin) getUserBadges(w http.ResponseWriter, r *http.Request, actingUserID string) {
// 	userID, ok := mux.Vars(r)["userID"]
// 	if !ok {
// 		userID = actingUserID
// 	}

// 	badges, err := p.store.GetUserBadges(userID)
// 	if err != nil {
// 		p.mm.Log.Debug("Error getting the badges for user", "error", err, "user", userID)
// 	}

// 	b, _ := json.Marshal(badges)
// 	_, _ = w.Write(b)
// }

func (p *Plugin) extractUserMiddleWare(handler HTTPHandlerFuncWithUser, responseType ResponseType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			errorMessage := "Not authorized."
			switch responseType {
			case ResponseTypeJSON:
				p.writeAPIError(w, &APIErrorResponse{ID: "", Message: errorMessage, StatusCode: http.StatusUnauthorized})
			case ResponseTypePlain:
				http.Error(w, errorMessage, http.StatusUnauthorized)
			case ResponseTypeDialog:
				dialogError(w, errorMessage, nil)
			case ResponseTypeAttachment:
				dialogError(w, errorMessage, nil)
			default:
				p.mm.Log.Error("Unknown ResponseType detected")
			}
			return
		}

		handler(w, r, userID)
	}
}

func (p *Plugin) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.mm.Log.Error("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func checkPluginRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// All other plugins are allowed
		pluginID := r.Header.Get("Mattermost-Plugin-ID")
		if pluginID == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (p *Plugin) writeAPIError(w http.ResponseWriter, apiErr *APIErrorResponse) {
	b, err := json.Marshal(apiErr)
	if err != nil {
		p.mm.Log.Warn("Failed to marshal API error", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(apiErr.StatusCode)

	_, err = w.Write(b)
	if err != nil {
		p.mm.Log.Warn("Failed to write JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) getPluginURL() string {
	urlP := p.mm.Configuration.GetConfig().ServiceSettings.SiteURL
	url := "/"
	if urlP != nil {
		url = *urlP
	}
	if url[len(url)-1] == '/' {
		url = url[0 : len(url)-1]
	}
	return url + "/plugins/" + manifest.Id
}

func (p *Plugin) getAPIURL() string {
	return p.getPluginURL() + APIPath
}

func (p *Plugin) getStaticURL() string {
	return p.getPluginURL() + StaticPath
}

func (p *Plugin) getDialogURL() string {
	return p.getPluginURL() + DialogPath
}

func (p *Plugin) getAttachmentURL() string {
	return p.getPluginURL() + AttachmentPath
}

func attachmentError(w http.ResponseWriter, text string) {
	r := model.PostActionIntegrationResponse{
		EphemeralText: "Error: " + text,
	}
	_, _ = w.Write(r.ToJson())
}

func dialogError(w http.ResponseWriter, text string, errors map[string]string) {
	r := model.SubmitDialogResponse{
		Error:  "Error: " + text,
		Errors: errors,
	}
	_, _ = w.Write(r.ToJson())
}

func attachmentOK(w http.ResponseWriter, text string) {
	r := model.PostActionIntegrationResponse{
		EphemeralText: text,
	}
	_, _ = w.Write(r.ToJson())
}

func dialogOK(w http.ResponseWriter) {
	r := model.SubmitDialogResponse{}
	_, _ = w.Write(r.ToJson())
}
