package main

const (
	CommandTrigger     = "quiz"
	CommandDisplayName = "Quiz commands"
	CommandDescription = "Created by Quiz plugin"

	BotUserName    = "quiz"
	BotDisplayName = "Quiz Bot"
	BotDescription = "Created by the Quiz plugin."

	APIPath                       = "/api/v1"
	DialogPath                    = "/dialog"
	DialogPathNameQuiz            = "/name"
	DialogPathChangeType          = "/type"
	DialogPathDelete              = "/delete"
	DialogPathAddQuestion         = "/add"
	DialogPathReviewQuestions     = "/review"
	DialogPathRemoveQuestion      = "/remove"
	AttachmentPath                = "/attachment"
	AttachmentPathNameQuiz        = "/name"
	AttachmentPathChangeType      = "/type"
	AttachmentPathDelete          = "/delete"
	AttachmentPathAddQuestion     = "/add"
	AttachmentPathReviewQuestions = "/review"
	AttachmentPathRemoveQuestion  = "/remove"
	AttachmentPathSave            = "/save"

	DialogTypeSelect = "select"
	DialogTypeText   = "text"
	DialogTypeBool   = "bool"

	AttachmentContextFieldID = "ID"

	DialogSubmissionFieldName        = "name"
	DialogSubmissionFieldType        = "type"
	DialogSubmissionFieldQuestion    = "question"
	DialogSubmissionFieldAnswer      = "answer"
	DialogSubmissionFieldWrongAnswer = "wrong_"

	IncorrectAnswerCount = 3
	Separator            = "-------------"
)
