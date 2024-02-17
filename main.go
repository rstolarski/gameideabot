package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rstolarski/gameideabot/command"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	log.Println("Adding handlers")
	s.AddHandler(command.Add)
	s.AddHandler(command.Del)
	s.AddHandler(command.List)
	s.AddHandler(command.Random)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err := s.Open()
	if err != nil {
		log.Fatalln(err)
	}

	defer s.Close()

	log.Println("Bot is online")

	// Prevent bot from turning off by listening to channel for specific message

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
