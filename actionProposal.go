package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	WARNINGEMOJI = "❔"
)

var (
	parseActionRegex = regexp.MustCompile(`%(?:[A-Za-z]+) ([A-Za-z]+) (?:"|“)([A-Za-z0-9\./:?,\- ]+)(?:"|”)`)
	parseNoInfoRegex = regexp.MustCompile(`%(?:[A-Za-z]+) ([A-Za-z]+)`)
)

func initActionTypes() {
	actionTypes = append(actionTypes, ActionType{name: "textchannelcreate", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "voicechannelcreate", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "channeldelete", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "kickmember", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: false})
	actionTypes = append(actionTypes, ActionType{name: "banmember", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: false})
	actionTypes = append(actionTypes, ActionType{name: "unbanmember", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "applyrole", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "removerole", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "addemoji", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "renameserver", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
	actionTypes = append(actionTypes, ActionType{name: "servericonchange", votingTimeMinutes: 30, approvalPercentage: 51, infoNeeded: true})
}

func parseActionProposal(msg *discordgo.MessageCreate, client *discordgo.Session) {
	commandWhole := parseActionRegex.FindAllStringSubmatch(msg.Content, -1)

	var actionType ActionType

	if len(commandWhole) != 0 {
		if len(commandWhole[0]) == 0 {
			commandWhole = parseNoInfoRegex.FindAllStringSubmatch(msg.Content, -1)
			if len(commandWhole) == 0 {
				client.MessageReactionAdd(msg.ChannelID, msg.ID, WARNINGEMOJI)
				return
			} else if len(commandWhole[0]) == 0 {
				client.MessageReactionAdd(msg.ChannelID, msg.ID, WARNINGEMOJI)
				return
			}
			actionType = findActionType(strings.ToLower(commandWhole[0][1]))

			if actionType.infoNeeded {
				client.MessageReactionAdd(msg.ChannelID, msg.ID, WARNINGEMOJI)
				return
			}
		} else {
			actionType = findActionType(strings.ToLower(commandWhole[0][1]))
		}

	} else {
		commandWhole = parseNoInfoRegex.FindAllStringSubmatch(msg.Content, -1)
		if len(commandWhole) == 0 {
			client.MessageReactionAdd(msg.ChannelID, msg.ID, WARNINGEMOJI)
			return
		} else if len(commandWhole[0]) == 0 {
			client.MessageReactionAdd(msg.ChannelID, msg.ID, WARNINGEMOJI)
			return
		}
		actionType = findActionType(strings.ToLower(commandWhole[0][1]))

		if actionType.infoNeeded {
			client.MessageReactionAdd(msg.ChannelID, msg.ID, WARNINGEMOJI)
			return
		}
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
		if len(msg.Mentions) == 0 {
			return
		}

		action.info = commandWhole[0][2]
		action.info2 = msg.Mentions[0].ID

		break
	case "removerole":
		if len(msg.Mentions) == 0 {
			return
		}

		action.info = commandWhole[0][2]
		action.info2 = msg.Mentions[0].ID

		break
	default:
		action.info = commandWhole[0][2]
	}

	legal := action.checkLegal()

	action.reactStatus(legal)

	if !legal {
		return
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
