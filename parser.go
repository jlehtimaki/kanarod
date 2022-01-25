package main

import (
	"encoding/json"
	"github.com/jlehtimaki/toornament-csgo/pkg/structs"
)

func (d *discordBot) voteDayInfo(t string) (structs.Team, error) {
	team := structs.Team{}
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
