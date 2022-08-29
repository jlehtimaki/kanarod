package main

import (
	"encoding/json"
	"fmt"
	s "github.com/jlehtimaki/toornament-csgo/pkg/structs"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (d *discordBot) team(teamName string) {
	team := s.Team{}

	// Create the userchannel
	channel, err := d.s.UserChannelCreate(d.mc.Author.ID)
	if err != nil {
		log.Error(err)
		return
	}
	data, err := d.rest.getTeam(teamName)
	if err != nil {
		d.s.ChannelMessageSend(channel.ID, "sorry could not get that for you")
		log.Error(err)
		return
	}

	err = json.Unmarshal([]byte(data), &team)
	if err != nil {
		d.s.ChannelMessageSend(channel.ID, "sorry could not get that for you")
		log.Error(err)
		return
	}

	// Players string
	var playerString string
	for _, p := range team.Players {
		pTmpString := fmt.Sprintf("> %s (%s)\n"+
			"\t\t Faceit:\n"+
			"\t\t\t Rank/Elo/Most Played/Least Played/KD:\n"+
			"\t\t\t\t %d/%d/%s/%s/%s \n"+
			"\t\t Esportal:\n"+
			"\t\t\t Rank: %s\n"+
			"\t\t MM Rank: %s\n"+
			"\t\t Kanaliiga:\n"+
			"\t\t\t Rating/ADR/Kills/Assists/Deaths\n"+
			"\t\t\t %.2f/%.2f/%d/%d/%d \n",
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
		"```%s```"+
		"",
		team.Name,
		team.CustomFields.CaptainDiscord,
		team.BestMap.Name, team.BestMap.WinRate,
		team.WorstMap.Name, team.WorstMap.WinRate,
		playerString)

	if len(teamString) >= 2000 {
		teamStringSplit := strings.Split(teamString, "```")
		for _, ts := range teamStringSplit {
			if ts == "" {
				continue
			}
			if len(ts) >= 2000 {
				playerSplit := strings.Split(ts, ">")
				for n, ps := range playerSplit {
					if n == 0 {
						continue
					}
					_, err = d.s.ChannelMessageSend(channel.ID, ps)
					if err != nil {
						log.Info("fooo")
						d.s.ChannelMessageSend(channel.ID, "sorry could not get that for you")
						log.Error(err)
						return
					}
				}
			} else {
				_, err = d.s.ChannelMessageSend(channel.ID, ts)
				if err != nil {
					d.s.ChannelMessageSend(channel.ID, "sorry could not get that for you")
					log.Error(err)
					return
				}
			}
		}
		return
	}
	_, err = d.s.ChannelMessageSend(channel.ID, teamString)
	if err != nil {
		d.s.ChannelMessageSend(channel.ID, "sorry could not get that for you")
		log.Error(err)
		return
	}
}
