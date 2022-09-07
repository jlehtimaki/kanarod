package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func voteDayMessage(teamName, captain, channelId string, s *discordgo.Session) {
	voteDayMessage := fmt.Sprintf("**Next game is going to be played against: *%s*, vote for the day**: \n"+
		":zero: - No can do \n"+
		":one: - Monday \n"+
		":two: - Tuesday \n"+
		":three: - Wednesday \n"+
		":four: - Thursday \n"+
		":five: - Friday \n"+
		":six: - Saturday \n"+
		":seven: - Sunday \n\n"+
		"Captain Discord: %s", teamName, captain)

	rMsg, err := s.ChannelMessageSend(channelId, voteDayMessage)
	if err != nil {
		log.Error(err)
	}
	log.Info(rMsg)

	days := []string{"0⃣", "1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7⃣"}

	for _, day := range days {
		err = s.MessageReactionAdd(channelId, rMsg.ID, day)
		if err != nil {
			log.Error(err)
		}
	}
}
