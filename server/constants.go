package main

const (
	CommandTrigger     = "quiz"
	CommandDisplayName = "Quiz commands"
	CommandDescription = "Created by Quiz plugin"

	BotUserName    = "quiz"
	BotDisplayName = "Quiz Bot"
	BotDescription = "Created by the Quiz plugin."

	APIPath = "/api/v1"

	DialogPath                = "/dialog"
	DialogPathNameQuiz        = "/name"
	DialogPathChangeType      = "/type"
	DialogPathDelete          = "/delete"
	DialogPathAddQuestion     = "/add"
	DialogPathReviewQuestions = "/review"
	DialogPathRemoveQuestion  = "/remove"
	DialogPathGameStart       = "/start"
	DialogPathScore           = "/score"
	DialogPathAnswer          = "/answer"

	AttachmentPath                = "/attachment"
	AttachmentPathNameQuiz        = "/name"
	AttachmentPathChangeType      = "/type"
	AttachmentPathDelete          = "/delete"
	AttachmentPathAddQuestion     = "/add"
	AttachmentPathReviewQuestions = "/review"
	AttachmentPathRemoveQuestion  = "/remove"
	AttachmentPathSave            = "/save"
	AttachmentPathSelectAnswer    = "/selectAnswer"
	AttachmentPathAnswer          = "/answer"
	AttachmentPathNext            = "/next"
	AttachmentPathScore           = "/score"

	DialogTypeSelect    = "select"
	DialogTypeBool      = "bool"
	DialogTypeText      = "text"
	DialogSubtypeNumber = "number"

	AttachmentContextFieldID         = "ID"
	AttachmentContextFieldCorrect    = "correct"
	AttachmentContextFieldGameID     = "gameID"
	AttachmentContextFieldQuestionID = "questionID"

	DialogSubmissionFieldName              = "name"
	DialogSubmissionFieldType              = "type"
	DialogSubmissionFieldQuestion          = "question"
	DialogSubmissionFieldAnswer            = "answer"
	DialogSubmissionFieldWrongAnswer       = "wrong_"
	DialogSubmissionFieldGameQuiz          = "quiz"
	DialogSubmissionFieldGameType          = "type"
	DialogSubmissionFieldGameScoring       = "scoring"
	DialogSubmissionFieldNumberOfQuestions = "nquestions"
	DialogSubmissionFieldGameAnswer        = "game_answer"

	IncorrectAnswerCount = 3
	Separator            = "-------------"
)
