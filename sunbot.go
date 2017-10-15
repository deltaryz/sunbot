package main

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"strconv"
)

// Global variables
var (
	commands         map[string]*command // verb string -> command object (see commands.go)
	DiscordAuthToken string              // environment variable DISCORD_AUTH_TOKEN
	DefaultPrefix    string              // environment variable COMMAND_PREFIX
	DebugEnabled     bool                // environment variable DEBUG_OUTPUT
)

// println, except only outputs if DEBUG_OUTPUT is true
func DebugPrint(output string) {
	if DebugEnabled {
		fmt.Println(output)
	}
}

func main() {

	// Initialize env configs
	DiscordAuthToken = os.Getenv("DISCORD_AUTH_TOKEN")
	DefaultPrefix = os.Getenv("COMMAND_PREFIX")
	DebugEnabled, _ = strconv.ParseBool(os.Getenv("DEBUG_OUTPUT"))

	// Initialize commands
	commands = initCommands()

	// Initialize discordgo
	discord, err := discordgo.New("Bot " + DiscordAuthToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// 	message handler
	discord.AddHandler(parseChatMessage)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
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
	discord.Close()

}

// Called any time a message is sent
func parseChatMessage(discordSession *discordgo.Session, msgEvent *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if msgEvent.Author.ID == discordSession.State.User.ID {
		return
	}

	// Make it easier to reference message text
	msg := msgEvent.Content

	DebugPrint("Message received.\n" + msgEvent.Author.Username + ": " + msg)

	// Did the message start with the command prefix?
	if msg[:1] == DefaultPrefix {

		DebugPrint("Message is a command.")

		// prepare variables to parse command
		args := strings.Split(msg[1:], " ")
		cmdInput := args[0]
		args = append(args[:0], args[1:]...)

		if cmd, ok := commands[cmdInput]; ok {
			discordSession.ChannelMessageSend(msgEvent.ChannelID, cmd.function(args)) // TODO: accomodate a message response struct rather than string
		} else {
			discordSession.ChannelMessageSend(msgEvent.ChannelID, "I don't understand that command.")
		}

		// TODO: implement command usage metrics
	} else {
		// TODO: implement metrics of standard chat messages
	}

}
