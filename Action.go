package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
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

func (a *Action) checkLegal() bool {
	switch a.actionType.name {
	case "channeldelete":
		channels, _ := client.GuildChannels(SERVERID)
		var selected *discordgo.Channel
		for _, channel := range channels {
			if channel.Name == a.info {
				selected = channel
			}
		}

		if Has(UNTOUCHABLETOPICS, selected.ParentID) {
			// don't break the core functions of the server mk thx
			return false
		} else if Has(UNTOUCHABLECHANNELS, selected.ID) {
			return false
		}
		break
	case "kickmember":
		if Has(UNTOUCHABLEUSERS, a.info) {
			return false
		}
		break
	case "banmember":
		if Has(UNTOUCHABLEUSERS, a.info) {
			return false
		}
		break
	case "applyrole":
		roles, _ := client.GuildRoles(SERVERID)
		var selected *discordgo.Role

		selected = nil

		for _, role := range roles {
			if role.Name == a.info {
				selected = role
			}
		}

		if selected == nil {
			return false
		}

		if Has(UNTOUCHABLEROLES, selected.ID) {
			// no retards no admin for you
			return false
		}
	case "removerole":
		roles, _ := client.GuildRoles(SERVERID)
		var selected *discordgo.Role

		selected = nil

		for _, role := range roles {
			if role.Name == a.info {
				selected = role
			}
		}

		if selected == nil {
			return false
		}

		if Has(UNTOUCHABLEROLES, selected.ID) {
			// no retards no admin for you
			return false
		}
	}

	return true
}

func (a *Action) reactStatus(status bool) {
	if status {
		client.MessageReactionAdd(ACTIONSUBMISSIONID, a.msgID, "✅")
	} else {
		client.MessageReactionAdd(ACTIONSUBMISSIONID, a.msgID, "❌")
	}
}
