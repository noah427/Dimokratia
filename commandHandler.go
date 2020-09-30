package main

import (
	"github.com/bwmarrin/discordgo"

	"os"
)

var ACTIONSUBMISSIONID string
var ACTIONVOTINGID string

func loadChannelIDs() {
	ACTIONSUBMISSIONID = os.Getenv("ACTIONSUBMISSION")
	ACTIONVOTINGID = os.Getenv("ACTIONVOTING")
}

func onMessage(client *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.ChannelID == ACTIONSUBMISSIONID {
		parseActionProposal(msg, client)
	}
}
