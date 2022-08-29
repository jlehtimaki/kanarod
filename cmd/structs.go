package main

import "github.com/bwmarrin/discordgo"

type discordBot struct {
	token       string
	accessToken string
	s           *discordgo.Session
	mc          *discordgo.MessageCreate
	apiPath     string
	autoMode    Automode
	rest        Rest
}
