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
	case "addemoji":
		client.GuildEmojiCreate(SERVERID, action.info, urlToDataScheme(action.info2), nil)
	case "servericonchange":
		client.GuildEdit(SERVERID, discordgo.GuildParams{
			Icon: urlToDataScheme(action.info),
		})
	case "voicechannelcreate":
		client.GuildChannelCreate(SERVERID, action.info, discordgo.ChannelTypeGuildVoice)
		break
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
	case "unbanmember":
		client.GuildBanDelete(SERVERID, action.info)
		break
	case "banmember":
		client.GuildBanCreate(SERVERID, action.info, 0)
		break
	case "applyrole":
		roles, _ := client.GuildRoles(SERVERID)
		var selected *discordgo.Role

		for _, role := range roles {
			if role.Name == action.info {
				selected = role
			}
		}

		client.GuildMemberRoleAdd(SERVERID, action.info2, selected.ID)

		break
	case "removerole":
		roles, _ := client.GuildRoles(SERVERID)
		var selected *discordgo.Role

		for _, role := range roles {
			if role.Name == action.info {
				selected = role
			}
		}

		client.GuildMemberRoleRemove(SERVERID, action.info2, selected.ID)
		break
	case "renameserver":
		client.GuildEdit(SERVERID, discordgo.GuildParams{
			Name: action.info,
		})
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
				{Name: "Action Information: ", Value: action.prettyPrintInfo()},
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
				{Name: "Action Information: ", Value: action.prettyPrintInfo()},
			},
		}
	}

	client.ChannelMessageSendEmbed(ACTIONRESULTSID, embed)
}
