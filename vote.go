package main

import (
	"fmt"
	_ "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (d *discordBot) vote(s string) {
	voteDayMessage := fmt.Sprintf("**%s**: \n"+
		":one: - Monday \n"+
		":two: - Tuesday \n"+
		":three: - Wednesday \n"+
		":four: - Thursday \n"+
		":five: - Friday \n"+
		":six: - Saturday \n"+
		":seven: - Sunday \n\n", s)

	rMsg, err := d.s.ChannelMessageSend(d.mc.ChannelID, voteDayMessage)
	if err != nil {
		log.Error(err)
	}
	log.Info(rMsg)

	days := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7⃣"}

	for _, day := range days {
		err = d.s.MessageReactionAdd(d.mc.ChannelID, rMsg.ID, day)
		if err != nil {
			log.Error(err)
		}
	}
}

func (d *discordBot) nextMatch(team string) {
	data, err := d.voteDayInfo(team)
	if err != nil {
		_, err = d.s.ChannelMessageSend(d.mc.ChannelID, "Sorry, could not get that for ya!")
		if err != nil {
			log.Error(err)
		}
		return
	}

	voteDayMessage := fmt.Sprintf("**Next game is going to be played against: *%s*, vote for the day**: \n"+
		":one: - Monday \n"+
		":two: - Tuesday \n"+
		":three: - Wednesday \n"+
		":four: - Thursday \n"+
		":five: - Friday \n"+
		":six: - Saturday \n"+
		":seven: - Sunday \n\n"+
		"Captain Discord: %s", data.Name, data.CustomFields.CaptainDiscord)

	rMsg, err := d.s.ChannelMessageSend(d.mc.ChannelID, voteDayMessage)
	if err != nil {
		log.Error(err)
	}
	log.Info(rMsg)

	days := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7⃣"}

	for _, day := range days {
		err = d.s.MessageReactionAdd(d.mc.ChannelID, rMsg.ID, day)
		if err != nil {
			log.Error(err)
		}
	}
}
