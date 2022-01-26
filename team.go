package main

import (
	"encoding/json"
	"fmt"
	s "github.com/jlehtimaki/toornament-csgo/pkg/structs"
	log "github.com/sirupsen/logrus"
)

func (d *discordBot) team(teamName string) {
	team := s.Team{}
	data, err := d.getTeam(teamName)
	if err != nil {
		d.s.ChannelMessageSend(d.mc.ChannelID, "sorry could not get that for you")
		log.Error(err)
		return
	}

	err = json.Unmarshal([]byte(data), &team)
	if err != nil {
		d.s.ChannelMessageSend(d.mc.ChannelID, "sorry could not get that for you")
		log.Error(err)
		return
	}

	// Players string
	var playerString string
	for _, p := range team.Players {
		pTmpString := fmt.Sprintf("\t**%s (%s)** \n"+
			"\t\t Faceit:\n"+
			"\t\t\t Rank/Elo/Most Player/Least Played/KD:\n"+
			"\t\t\t\t %d/%d/%s/%s/%s \n"+
			"\t\t Esportal:\n"+
			"\t\t\t Rank %d:\n"+
			"\t\t MM Rank: %d\n"+
			"\t\t Kanaliiga:\n"+
			"\t\t\t Rating/ADR/Kills/Assists/Deaths\n"+
			"\t\t\t\t %.2f/%.2f/%d/%d/%d \n",
			p.Name, p.CustomFields.SteamId,
			p.Faceit.Rank, p.Faceit.Elo, p.Faceit.MostPlayedMap.Name, p.Faceit.LeastPlayedMap.Name, p.Faceit.KD,
			p.Esportal.Rank,
			p.MM.Rank,
			p.Kanaliiga.KanaRating, p.Kanaliiga.ADR, p.Kanaliiga.Kills, p.Kanaliiga.Assists, p.Kanaliiga.Deaths)

		playerString = playerString + pTmpString
	}

	// Set the team information string
	teamString := fmt.Sprintf("***Team***: %s \n\n"+
		"***Captain:***: %s\n"+
		"***Best Map***: %s (%s%%)\n"+
		"***Worst Map***: %s (%s%%)\n\n"+
		"***Lineup:***: \n"+
		"%s"+
		"",
		team.Name,
		team.CustomFields.CaptainDiscord,
		team.BestMap.Name, team.BestMap.WinRate,
		team.WorstMap.Name, team.WorstMap.WinRate,
		playerString)

	_, err = d.s.ChannelMessageSend(d.mc.ChannelID, teamString)
	if err != nil {
		d.s.ChannelMessageSend(d.mc.ChannelID, "sorry could not get that for you")
		log.Error(err)
		return
	}
}