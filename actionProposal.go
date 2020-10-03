package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	parseActionRegex = regexp.MustCompile(`%([A-Za-z]+) ([A-Za-z]+) "([A-Za-z ]+)"`)
)

type ActionType struct {
	name               string
	approvalPercentage int
	votingTimeHours    int
}

var (
	actionTypes []ActionType
)

type Action struct {
	actionType  ActionType
	info        string
	authorID    string
	time        time.Time
	msgID       string
	votingMsgID string
	votesUp     float32
	votesDown   float32
}

func initActionTypes() {
	actionTypes = append(actionTypes, ActionType{name: "textchannelcreate", votingTimeHours: 1, approvalPercentage: 51})
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

	action := Action{
		actionType: findActionType(strings.ToLower(commandWhole[0][2])),
		info:       commandWhole[0][3],
		authorID:   msg.Author.ID,
		msgID:      msg.ID,
		time:       time.Now(),
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
			{Name: "Info: ", Value: action.info},
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

	time.AfterFunc(time.Minute*time.Duration(action.actionType.votingTimeHours), func() {
		actionResult(action, client)
	})

}

func formatTime(action Action) string {
	hour := action.time.Hour() + action.actionType.votingTimeHours
	var timeMeridian string
	hour12 := (hour % 12)

	if (hour%24)/12 >= 1 {
		timeMeridian = "PM"
	} else {
		timeMeridian = "AM"
	}

	if hour12 == 0 {
		hour12 = 12
	}

	// 06 != 6 | 6:30 != 6:3

	return fmt.Sprintf("%d hours, at %d:%d %s", action.actionType.votingTimeHours, hour12, action.time.Minute(), timeMeridian)
}
