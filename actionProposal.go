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
	name            string
	votingTimeHours int
}

var (
	actionTypes []ActionType
)

type Action struct {
	actionType ActionType
	info       string
	authorID   string
	time       time.Time
	msgID      string
	votes      int
}

func initActionTypes() {
	actionTypes = append(actionTypes, ActionType{name: "textchannelcreate", votingTimeHours: 12})
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


	client.ChannelMessageSendEmbed(ACTIONVOTINGID, embed)

}

func formatTime(action Action) string {
	hour := action.time.Hour() + action.actionType.votingTimeHours
	var timeMeridian string
	hour12 := (hour % 12)

	if (hour % 24)/12 >= 1 {
		timeMeridian = "PM"
	} else {
		timeMeridian = "AM"
	}

	return fmt.Sprintf("%d hours, at %d:%d %s", action.actionType.votingTimeHours, hour12, action.time.Minute(), timeMeridian)
}
