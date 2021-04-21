package main

import (
	"fmt"

	commandparser "github.com/larkox/mattermost-plugin-quiz/server/command_parser"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func getHelp() string {
	return `Available Commands:
`
}

func (p *Plugin) getCommand() *model.Command {
	return &model.Command{
		Trigger:          CommandTrigger,
		DisplayName:      CommandDisplayName,
		Description:      CommandDescription,
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands:",
		AutoCompleteHint: "[command]",
		AutocompleteData: p.getAutocompleteData(),
	}
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	p.mm.Post.SendEphemeralPost(args.UserId, post)
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	stringArgs := commandparser.Parse(args.Command)
	lengthOfArgs := len(stringArgs)
	restOfArgs := []string{}

	var handler func([]string, *model.CommandArgs) (bool, *model.CommandResponse, error)
	if lengthOfArgs == 1 {
		p.postCommandResponse(args, getHelp())
		return &model.CommandResponse{}, nil
	}
	command := stringArgs[1]
	if lengthOfArgs > 2 {
		restOfArgs = stringArgs[2:]
	}
	switch command {
	case "test-clean":
		handler = p.runClean
	case "create":
		handler = p.runCreate
	default:
		p.postCommandResponse(args, getHelp())
		return &model.CommandResponse{}, nil
	}
	isUserError, resp, err := handler(restOfArgs, args)
	if err != nil {
		if isUserError {
			p.postCommandResponse(args, fmt.Sprintf("__Error: %s.__\n\nRun `/todo help` for usage instructions.", err.Error()))
		} else {
			p.mm.Log.Error(err.Error())
			p.postCommandResponse(args, "An unknown error occurred. Please talk to your system administrator for help.")
		}
	}

	if resp != nil {
		return resp, nil
	}

	return &model.CommandResponse{}, nil
}

func (p *Plugin) runClean(args []string, extra *model.CommandArgs) (bool, *model.CommandResponse, error) {
	_ = p.mm.KV.DeleteAll()
	p.postCommandResponse(extra, "Database cleaned")
	return emptyCommandResponse()
}

func (p *Plugin) runCreate(args []string, extra *model.CommandArgs) (bool, *model.CommandResponse, error) {
	post := &model.Post{
		Message: "Creating quiz",
	}
	q := &Quiz{}
	model.ParseSlackAttachment(post, p.CreateAttachmentFromQuiz(q))

	err := p.mm.Post.DM(p.BotUserID, extra.UserId, post)
	if err != nil {
		p.postCommandResponse(extra, "Error: "+err.Error())
		return emptyCommandResponse()
	}

	q.ID = post.Id
	err = p.store.StoreQuiz(q)
	if err != nil {
		p.postCommandResponse(extra, "Error: "+err.Error())
		return emptyCommandResponse()
	}

	p.postCommandResponse(extra, "The bot will contact you soon and guide you through the creation process.")
	return emptyCommandResponse()
}

func emptyCommandResponse() (bool, *model.CommandResponse, error) {
	return false, &model.CommandResponse{}, nil
}

func (p *Plugin) getAutocompleteData() *model.AutocompleteData {
	// badges := model.NewAutocompleteData("badges", "[command]", "Available commands: grant")

	// grant := model.NewAutocompleteData("grant", "--user @username --badge id", "Grant a badge to a user")
	// grant.AddNamedDynamicListArgument("badge", "--badge badgeID", getAutocompletePath(AutocompletePathBadgeSuggestions), true)
	// grant.AddNamedTextArgument("user", "User to grant the badge to", "--user @username", "", true)
	// badges.AddCommand(grant)

	// create := model.NewAutocompleteData("create", "badge | type", "Create a badge or a type")

	// badge := model.NewAutocompleteData(
	// 	"badge",
	// 	"--name badgeName --description badgeDescription --image :image: --type typeID --multiple true|false",
	// 	"Create a badge",
	// )
	// badge.AddNamedTextArgument("name", "Name of the badge", "--name badgeName", "", true)
	// badge.AddNamedTextArgument("description", "Description of the badge", "--description description", "", true)
	// badge.AddNamedTextArgument("image", "Image of the badge", "--image :image:", "", true)
	// badge.AddNamedDynamicListArgument("type", "Type of the badge", getAutocompletePath(AutocompletePathTypeSuggestions), true)
	// badge.AddNamedStaticListArgument("multiple", "Whether the badge can be granted multiple times", true, []model.AutocompleteListItem{
	// 	{Item: TrueString},
	// 	{Item: FalseString},
	// })
	// create.AddCommand(badge)

	// createType := model.NewAutocompleteData(
	// 	"type",
	// 	"--name typeName --everyoneCanCreate true|false --everyoneCanGrant true|false",
	// 	"Create a badge type",
	// )
	// createType.AddNamedTextArgument("name", "Name of the type", "--name typeName", "", true)
	// createType.AddNamedStaticListArgument("everyoneCanCreate", "Whether the badge can be granted by everyone", true, []model.AutocompleteListItem{
	// 	{Item: TrueString},
	// 	{Item: FalseString},
	// })
	// createType.AddNamedStaticListArgument("everyoneCanGrant", "Whether the badge can be created by everyone", true, []model.AutocompleteListItem{
	// 	{Item: TrueString},
	// 	{Item: FalseString},
	// })
	// create.AddCommand(createType)

	// badges.AddCommand(create)
	// return badges
	return nil
}
