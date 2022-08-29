package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	s "github.com/jlehtimaki/toornament-csgo/pkg/structs"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Automode struct {
	Redis          Redis
	Rest           Rest
	SyncTime       int
	Teams          []Team
	DiscordSession *discordgo.Session
}

type Team struct {
	Name             string
	Scheduled        []s.ScheduledMatch
	NextMatch        s.Team
	ChannelId        string
	NotificationSend bool
	VoteSend         bool
}

func automodeInit(session *discordgo.Session) (Automode, error) {
	redis, err := newRedis()
	if err != nil {
		return Automode{}, err
	}
	syncTime, err := strconv.Atoi(os.Getenv("AUTOMODE_SYNCTIME"))
	if err != nil {
		syncTime = 5
	}
	a := Automode{
		Redis:          redis,
		SyncTime:       syncTime,
		DiscordSession: session,
		Rest:           initRest(),
	}
	return a, nil
}

func (a *Automode) start() {
	log.Info("getting the teams from database")
	a.Teams = a.getTeams()
	ticker := time.NewTicker(time.Second * time.Duration(a.SyncTime))
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			fmt.Println("it's automating")
			a.autoLogic()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (a *Automode) getTeams() []Team {
	var teams []Team
	keys, err := a.Redis.getKeys("*")
	if err != nil {
		log.Error(err)
	}
	for _, k := range keys {
		v, err := a.Redis.getValue(k)
		if err != nil {
			log.Error(err)
			return []Team{}
		}
		team := Team{}
		err = json.Unmarshal([]byte(v), &team)
		if err != nil {
			log.Error()
			return []Team{}
		}
		teams = append(teams, team)
	}

	return teams
}

func (a *Automode) addTeam(teamName string, channelID string) {
	log.Infof("adding team %s to automation process", teamName)
	if a.alreadyInAuto(teamName) {
		log.Infof("team already in automation process")
		a.DiscordSession.ChannelMessageSend(channelID, "sorry team is already in automode")
		return
	}
	team := Team{Name: teamName, ChannelId: channelID}

	log.Infof("getting next match information for team %s", teamName)
	nextMatch, err := a.Rest.getNextOpponent(teamName)
	if err != nil {
		log.Error(err)
		a.DiscordSession.ChannelMessageSend(channelID, "sorry could get information for the team in automode")
		return
	}
	err = json.Unmarshal([]byte(nextMatch), &team.NextMatch)
	if err != nil {
		log.Error(err)
		a.DiscordSession.ChannelMessageSend(channelID, "sorry could get information for the team in automode")
		return
	}

	log.Infof("getting scheduled match information for team %s", teamName)
	scheduledData, err := a.Rest.getScheduled(teamName)
	if err != nil {
		log.Error(err)
		a.DiscordSession.ChannelMessageSend(channelID, "sorry could get information for the team in automode")
		return
	}
	err = json.Unmarshal([]byte(scheduledData), &team.Scheduled)
	if err != nil {
		log.Error(err)
		a.DiscordSession.ChannelMessageSend(channelID, "sorry could get information for the team in automode")
		return
	}
	err = a.Redis.addKey(teamName, team)
	if err != nil {
		log.Error(err)
		a.DiscordSession.ChannelMessageSend(channelID, "sorry could not put the team in automode")
		return
	}
	a.Teams = append(a.Teams, team)
	a.DiscordSession.ChannelMessageSend(channelID, fmt.Sprintf("team %s added to the automode! You can "+
		"disable me by ?disable command", teamName))
}

func (a *Automode) removeTeam(teamName string, channelId string) {
	log.Infof("removing team %s to automation process", teamName)
	for i, team := range a.Teams {
		if team.Name == teamName {
			a.Teams = remove(a.Teams, i)
			err := a.Redis.removeKey(teamName)
			if err != nil {
				log.Error(err)
				a.DiscordSession.ChannelMessageSend(channelId, "could not remove the team from automode")
				return
			}
			a.DiscordSession.ChannelMessageSend(channelId, "team removed")
			return
		}
	}
	log.Errorf("could not remove the team %s, not found", teamName)
	a.DiscordSession.ChannelMessageSend(channelId, "did not find the correct team to remove")
}

func (a *Automode) alreadyInAuto(teamName string) bool {
	for _, x := range a.Teams {
		if x.Name == teamName {
			return true
		}
	}
	return false
}

func (a *Automode) autoLogic() {
	today := time.Now()
	for i, team := range a.Teams {
		// Remind of scheduled match in the morning
		if !team.NotificationSend {
			gameDayText := "🎮 🔥 **GAME DAY** 🔥 🎮 \n"
			for _, match := range team.Scheduled {
				gameDayText = fmt.Sprintf("%s**%s** - **%s** \n%s	|	Server Number: %d	|	📺: %s	\n",
					gameDayText, match.Team1Name, match.Team2Name, match.Date.Format(time.Kitchen), match.ServerID, match.Stream)
			}
			if checkDate(today, team.Scheduled[0].Date) {
				a.DiscordSession.ChannelMessageSend(team.ChannelId, gameDayText)
				a.Teams[i].NotificationSend = true
			}
			a.saveData(a.Teams[i])
		}

		// Next Match notification
		if !team.VoteSend {
			if today.Weekday().String() == "Monday" {
				voteDayMessage(team.NextMatch.Name, team.NextMatch.CustomFields.CaptainDiscord,
					team.ChannelId, a.DiscordSession)
			}
			a.Teams[i].VoteSend = true
			a.saveData(a.Teams[i])
		}

		// TODO: Game day checks and scores
		// Game Day Check scores
		//if !gameDay(team.Scheduled, today) {
		//	for _, match := range team.Scheduled {
		//		fmt.Println(match)
		//	}
		//}

		// Update if necessary, if updates --> save to redis
		scheduledData, err := a.Rest.getScheduled(team.Name)
		if err != nil {
			log.Error(err)
			break
		}
		var scheduledMatches []s.ScheduledMatch
		err = json.Unmarshal([]byte(scheduledData), &scheduledMatches)
		if err != nil {
			log.Error(err)
			break
		}
		if (scheduledMatches[0].Team2Name != team.Scheduled[0].Team2Name) && scheduledMatches[1].Team1Name != team.Scheduled[1].Team1Name {
			a.Teams[i].Scheduled = scheduledMatches
			a.Teams[i].NotificationSend = false

			nextMatchData, err := a.Rest.getNextOpponent(team.Name)
			if err != nil {
				log.Error(err)
				break
			}
			var nextMatch s.Team
			err = json.Unmarshal([]byte(nextMatchData), &nextMatch)
			if err != nil {
				log.Error(err)
				break
			}
			a.Teams[i].NextMatch = nextMatch
			a.Teams[i].VoteSend = false

			a.saveData(a.Teams[i])
		}
	}
}

func remove(slice []Team, s int) []Team {
	return append(slice[:s], slice[s+1:]...)
}

func checkDate(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	notificationDate := time.Date(y2, m2, d2, 8, 0, 0, 0, date2.Location())
	return y1 == y2 && m1 == m2 && d1 == d2 && date1.After(notificationDate)
}

func gameDay(matches []s.ScheduledMatch, today time.Time) bool {
	var match s.ScheduledMatch
	for _, m := range matches {
		if match.ID == 0 {
			match = m
			continue
		}
		if match.Date.After(m.Date) {
			match = m
			continue
		}
	}
	if today.Before(match.Date) {
		return false
	}
	return true
}

func (a *Automode) saveData(team Team) {
	err := a.Redis.addKey(team.Name, team)
	if err != nil {
		log.Error("could not save notification send information")
	}
}
