package main

import (
	"github.com/bwmarrin/discordgo"

	"os"
)

var SERVERID string
var ACTIONRESULTSID string
var ACTIONSUBMISSIONID string
var ACTIONVOTINGID string

func loadChannelIDs() {
	SERVERID = os.Getenv("SERVERID")
	ACTIONSUBMISSIONID = os.Getenv("ACTIONSUBMISSION")
	ACTIONVOTINGID = os.Getenv("ACTIONVOTING")
	ACTIONRESULTSID = os.Getenv("ACTIONRESULTS")
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
		%propose applyrole "role name" @username
		%propose removerole "role name" @username
		`)
	}
}
