package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	parseActionRegex  = regexp.MustCompile(`%(?:[A-Za-z]+) ([A-Za-z]+) "([A-Za-z0-9\./:?, ]+)"`)
	parseNoInfoRegex  = regexp.MustCompile(`%(?:[A-Za-z]+) ([A-Za-z]+)`)
	parseMentionRegex = regexp.MustCompile(`%(?:[A-Za-z]+) (?:[A-Za-z]+) "(?:[A-Za-z0-9\./:?, ]+)" <@!(\d+)>`)
)

type ActionType struct {
	name               string
	approvalPercentage int
	votingTimeMinutes  int
}

var (
	actionTypes []ActionType
)

type Action struct {
	actionType  ActionType
	info        string
	info2       string
	authorID    string
	time        time.Time
	msgID       string
	votingMsgID string
	votesUp     float32
	votesDown   float32
}

func (a *Action) prettyPrintInfo() string {
	var response string
	switch a.actionType.name {
	case "kickmember":
		user, _ := client.User(a.info)
		response = user.Username
		break
	case "unbanmember":
		user, _ := client.User(a.info)
		response = user.Username
		break
	case "banmember":
		user, _ := client.User(a.info)
		response = user.Username
		break
	case "applyrole":
		user, _ := client.User(a.info2)
		response = fmt.Sprintf("Role name = %s, Username = %s", a.info, user.Username)
		break
	case "removerole":
		user, _ := client.User(a.info2)
		response = fmt.Sprintf("Role name = %s, Username = %s", a.info, user.Username)
		break
	default:
		response = a.info
	}
	return response
}

func initActionTypes() {
	actionTypes = append(actionTypes, ActionType{name: "textchannelcreate", votingTimeMinutes: 30, approvalPercentage: 51})
	actionTypes = append(actionTypes, ActionType{name: "channeldelete", votingTimeMinutes: 30, approvalPercentage: 51})
	actionTypes = append(actionTypes, ActionType{name: "kickmember", votingTimeMinutes: 30, approvalPercentage: 51})
	actionTypes = append(actionTypes, ActionType{name: "banmember", votingTimeMinutes: 30, approvalPercentage: 51})
	actionTypes = append(actionTypes, ActionType{name: "unbanmember", votingTimeMinutes: 30, approvalPercentage: 51})
	actionTypes = append(actionTypes, ActionType{name: "applyrole", votingTimeMinutes: 30, approvalPercentage: 51})
	actionTypes = append(actionTypes, ActionType{name: "removerole", votingTimeMinutes: 30, approvalPercentage: 51})
	actionTypes = append(actionTypes, ActionType{name: "addemoji", votingTimeMinutes: 30, approvalPercentage: 51})
}

func findActionType(actionType string) ActionType {
	for _, action := range actionTypes {
		if action.name == actionType {
			return action
		}
	}
	return ActionType{}
}

func parseActionProposal(msg *discordgo.MessageCreate, client *discordgo.Session) {
	commandWhole := parseActionRegex.FindAllStringSubmatch(msg.Content, -1)

	var actionType ActionType

	if len(commandWhole) != 0 {
		actionType = findActionType(strings.ToLower(commandWhole[0][1]))
	} else {
		commandWhole = parseNoInfoRegex.FindAllStringSubmatch(msg.Content, -1)
		if len(commandWhole) == 0 {
			return
		}
		actionType = findActionType(strings.ToLower(commandWhole[0][1]))
	}

	action := Action{
		actionType: actionType,
		authorID:   msg.Author.ID,
		msgID:      msg.ID,
		time:       time.Now(),
	}

	switch actionType.name {
	case "addemoji":
		info := strings.Split(commandWhole[0][2], ",")
		action.info = info[0]
		action.info2 = info[1]
	case "kickmember":
		action.info = msg.Mentions[0].ID
		if len(msg.Mentions) == 0 {
			return
		}
		break
	case "unbanmember":
		action.info = commandWhole[0][2]
		break
	case "banmember":
		if len(msg.Mentions) == 0 {
			return
		}
		action.info = msg.Mentions[0].ID
		break
	case "applyrole":
		secondParsing := parseMentionRegex.FindAllStringSubmatch(msg.Content, -1)

		action.info = commandWhole[0][2]
		action.info2 = secondParsing[0][1]

		break
	case "removerole":
		secondParsing := parseMentionRegex.FindAllStringSubmatch(msg.Content, -1)

		action.info = commandWhole[0][2]
		action.info2 = secondParsing[0][1]

		break
	default:
		action.info = commandWhole[0][2]
	}

	embed := &discordgo.MessageEmbed{
		Title: "Action",

		Color: 3108255,

		Author: &discordgo.MessageEmbedAuthor{
			Name:    msg.Author.Username,
			IconURL: msg.Author.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Type: ", Value: action.actionType.name},
			{Name: "Info: ", Value: action.prettyPrintInfo()},
			{Name: "Ends in: ", Value: formatTime(action)},
		},

		Footer: &discordgo.MessageEmbedFooter{
			Text: "Written by [REDACTED]#4242",
		},
	}

	message, _ := client.ChannelMessageSendEmbed(ACTIONVOTINGID, embed)

	client.MessageReactionAdd(message.ChannelID, message.ID, "👍")
	client.MessageReactionAdd(message.ChannelID, message.ID, "👎")

	action.votingMsgID = message.ID

	time.AfterFunc(time.Minute*time.Duration(action.actionType.votingTimeMinutes), func() {
		actionResult(action, client)
	})

}

func formatTime(action Action) string {
	// hour := action.time.Hour() + action.actionType.votingTimeMinutes
	// var timeMeridian string
	// hour12 := (hour % 12)

	// if (hour%24)/12 >= 1 {
	// 	timeMeridian = "PM"
	// } else {
	// 	timeMeridian = "AM"
	// }

	// if hour12 == 0 {
	// 	hour12 = 12
	// }

	// 06 != 6 | 6:30 != 6:3

	return fmt.Sprintf("%d minutes", action.actionType.votingTimeMinutes)
}
