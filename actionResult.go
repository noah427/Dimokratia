package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func actionResult(action Action, client *discordgo.Session) {

	votesUp, _ := client.MessageReactions(ACTIONVOTINGID, action.votingMsgID, "👍", -1, "", "")
	votesDown, _ := client.MessageReactions(ACTIONVOTINGID, action.votingMsgID, "👎", -1, "", "")

	votesUpCount := float32(len(votesUp))
	votesDownCount := float32(len(votesDown))

	result := action.actionType.approvalPercentage < int(((votesUpCount) / (votesUpCount + votesDownCount) * 100))

	action.votesUp = votesUpCount
	action.votesDown = votesDownCount

	postActionResults(action, client, result)

	if !result {
		return
	}

	switch action.actionType.name {
	case "textchannelcreate":
		client.GuildChannelCreate(SERVERID, action.info, discordgo.ChannelTypeGuildText)
		break
	case "channeldelete":
		channels, _ := client.GuildChannels(SERVERID)
		var selected *discordgo.Channel
		for _, channel := range channels {
			if channel.Name == action.info {
				selected = channel
			}
		}

		client.ChannelDelete(selected.ID)
		break
	case "kickmember":
		client.GuildMemberDelete(SERVERID, action.info)
		break
	case "banmember":
		client.GuildBanCreate(SERVERID, action.info, 10)
		break
	}

}

func postActionResults(action Action, client *discordgo.Session, result bool) {
	var embed *discordgo.MessageEmbed
	if result {
		embed = &discordgo.MessageEmbed{
			Title: "Vote Succeeded",
			Color: 3079834,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Upvotes: ", Value: strconv.Itoa(int(action.votesUp))},
				{Name: "Downvotes: ", Value: strconv.Itoa(int(action.votesDown))},
				{Name: "Percentage: ", Value: strconv.Itoa(int(((action.votesUp) / (action.votesUp + action.votesDown) * 100)))},
				{Name: "Action Type: ", Value: action.actionType.name},
				{Name: "Action Information: ", Value: action.info},
			},
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title: "Vote Failed",
			Color: 11797508,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Upvotes: ", Value: strconv.Itoa(int(action.votesUp))},
				{Name: "Downvotes: ", Value: strconv.Itoa(int(action.votesDown))},
				{Name: "Percentage: ", Value: strconv.Itoa(int(((action.votesUp) / (action.votesUp + action.votesDown) * 100)))},
				{Name: "Action Type: ", Value: action.actionType.name},
				{Name: "Action Information: ", Value: action.info},
			},
		}
	}

	client.ChannelMessageSendEmbed(ACTIONRESULTSID, embed)
}
