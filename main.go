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
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("could not find TOKEN environment variable")
	}
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	db := discordBot{
		token: token,
		s:     session,
		mc:    nil,
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

	if strings.Contains(m.Content, "?help") {
		log.Infof("user %s asked %s", m.Author, m.Content)
		helpMessage := "**Help:** \n\n" +
			"**?next-match** *team-name*		-		Vote for the day of teams next match\n" +
			"**?vote** *custom-string*		-		Vote for the day of custom string\n" +
			"**?team** *team-name*			-		Shows information about that team\n"
		_, err := s.ChannelMessageSend(m.ChannelID, helpMessage)
		if err != nil {
			log.Error(err)
		}
	}

	if strings.Contains(m.Content, "?next-match") {
		log.Infof("user %s asked %s", m.Author, m.Content)
		teamName := strings.Split(m.Content, "?next-match ")[1]
		d.nextMatch(teamName)
	}

	if strings.Contains(m.Content, "?vote") {
		log.Infof("user %s asked %s", m.Author, m.Content)
		customString := strings.Split(m.Content, "?vote ")[1]
		d.vote(customString)
	}

	if strings.Contains(m.Content, "?team") {
		log.Infof("user %s asked %s", m.Author, m.Content)
		teamName := strings.Split(m.Content, "?team ")[1]
		d.team(teamName)
	}
}
