package main

import "github.com/bwmarrin/discordgo"

type discordBot struct {
	s       *discordgo.Session
	mc      *discordgo.MessageCreate
	apiPath string
}
