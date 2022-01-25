package main

import "github.com/bwmarrin/discordgo"

type discordBot struct {
	token   string
	s       *discordgo.Session
	mc      *discordgo.MessageCreate
	apiPath string
}
