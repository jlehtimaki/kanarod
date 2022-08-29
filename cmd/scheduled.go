package main

import (
	"encoding/json"
	"fmt"
	s "github.com/jlehtimaki/toornament-csgo/pkg/structs"
	log "github.com/sirupsen/logrus"
)

func (d *discordBot) scheduledMatches(teamName string) {
	var scheduledMatches []s.ScheduledMatch
	returnString := fmt.Sprintf("**Scheduled matches for team** *%s*\n\n", teamName)

	data, err := d.rest.getScheduled(teamName)
	if err != nil {
		d.s.ChannelMessageSend(d.mc.ChannelID, "sorry could not get that for you")
		log.Error(err)
		return
	}
	err = json.Unmarshal([]byte(data), &scheduledMatches)
	if err != nil {
		d.s.ChannelMessageSend(d.mc.ChannelID, "sorry could not get that for you")
		log.Error(err)
		return
	}

	for _, match := range scheduledMatches {
		var opponent string
		if match.Team1Name == teamName {
			opponent = match.Team2Name
		} else {
			opponent = match.Team1Name
		}
		s := fmt.Sprintf("\t**Opponent**: %s \n"+
			"\t**Date**: %s \n"+
			"\t**Stream**: %s \n\n", opponent, match.Date, match.Stream)
		returnString = returnString + s
	}
	d.s.ChannelMessageSend(d.mc.ChannelID, returnString)
}
