package main

import (
	"encoding/json"
	s "github.com/jlehtimaki/toornament-csgo/pkg/structs"
)

func (d *discordBot) voteDayInfo(t string) (s.Team, error) {
	team := s.Team{}
	data, err := d.getNextOpponent(t)
	if err != nil {
		return team, err
	}

	err = json.Unmarshal([]byte(data), &team)
	if err != nil {
		return team, err
	}

	return team, nil
}
