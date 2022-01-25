package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var token = "OTM0MDI2MzAwMzE2ODYwNDE2.YeqFxw.0mdKstmBjB1SnxWoRWCBc-sUMEw"

func main() {
	dBot := initDiscord()
	initRest(&dBot)

	// Register the messageCreate func as a callback for MessageCreate events.
	dBot.s.AddHandler(dBot.messageCreate)

	// In this example, we only care about receiving message events.
	dBot.s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := dBot.s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = dBot.s.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func initDiscord() discordBot {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	db := discordBot{
		s:  session,
		mc: nil,
	}
	return db
}

func (d *discordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	d.s = s
	d.mc = m
	// If the message is "ping" reply with "Pong!"
	//if m.Content == "ping" {
	//	s.ChannelMessageSend(m.ChannelID, "Pong!")
	//}
	//
	//// If the message is "pong" reply with "Ping!"
	//if m.Content == "pong" {
	//	s.ChannelMessageSend(m.ChannelID, "Ping!")
	//}
	if strings.Contains(m.Content, "?vote") {
		teamName := strings.Split(m.Content, "?vote ")[1]
		d.voteDay(teamName)
	}
}
