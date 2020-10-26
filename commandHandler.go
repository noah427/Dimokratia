package main

import (
	"github.com/bwmarrin/discordgo"

	"os"
	"strings"
)

var SERVERID string
var ACTIONRESULTSID string
var ACTIONSUBMISSIONID string
var ACTIONVOTINGID string

//
var UNTOUCHABLEROLES []string
var UNTOUCHABLETOPICS []string
var UNTOUCHABLECHANNELS []string
var UNTOUCHABLEUSERS []string

func loadChannelIDs() {
	SERVERID = os.Getenv("SERVERID")
	ACTIONSUBMISSIONID = os.Getenv("ACTIONSUBMISSION")
	ACTIONVOTINGID = os.Getenv("ACTIONVOTING")
	ACTIONRESULTSID = os.Getenv("ACTIONRESULTS")
	//
	UNTOUCHABLEROLES = strings.Split(os.Getenv("UNTOUCHABLEROLES"), ",")
	UNTOUCHABLETOPICS = strings.Split(os.Getenv("UNTOUCHABLETOPICS"), ",")
	UNTOUCHABLECHANNELS = strings.Split(os.Getenv("UNTOUCHABLECHANNELS"), ",")
	UNTOUCHABLEUSERS = strings.Split(os.Getenv("UNTOUCHABLEUSERS"), ",")
}

func onMessage(client *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.ChannelID == ACTIONSUBMISSIONID {
		parseActionProposal(msg, client)
		return
	}

	if msg.Content == "%help" {
		client.ChannelMessageSend(msg.ChannelID, `
		commands:
		%propose textchannelcreate "channel name"
		%propose channeldelete "channel name"
		%propose kickmember @username
		%propose banmember @username
		%propose unbanmember "userID"
		%propose applyrole "role name" @username
		%propose removerole "role name" @username
		%propose addemoji "emojiName,discordfileurl"
		`)
	}

	if msg.Content == "%roles" {
		client.ChannelMessageSend(msg.ChannelID, `

		Text-Enforcer
		Voice-Enforcer
		Muted
		VC-Muted
		`)
	}
}
