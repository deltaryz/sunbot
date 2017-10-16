package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	version = "0.1 Dev"
)

// Environment variables
type config struct {
	DiscordAuthToken     string `env:"DISCORD_AUTH_TOKEN,required"`      // environment variable DISCORD_AUTH_TOKEN
	DefaultPrefix        string `env:"COMMAND_PREFIX" envDefault:"."`    // environment variable COMMAND_PREFIX
	DebugEnabled         bool   `env:"DEBUG_OUTPUT" envDefault:"true"`   // environment variable DEBUG_OUTPUT
	SillyCommandsEnabled bool   `env:"SILLY_COMMANDS" envDefault:"true"` // environment variable SILLY_COMMANDS
}

// Global variables
var (
	commands map[string]*command // verb string -> command object (see commands.go)
	cfg      config
)

func init() {
	rand.Seed(time.Now().Unix())
}

// randomRange gives a random whole integer between the given integers [min, max)
func RandomRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// println, except only outputs if DEBUG_OUTPUT is true
func DebugPrint(output string) {
	if cfg.DebugEnabled {
		fmt.Println(output)
	}
}

func main() {

	cfg = config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Error processing environment variables.\nPlease check https://github.com/techniponi/sunbot for details.\n\n" + err.Error())
		return
	}

	DebugPrint("Command prefix: " + cfg.DefaultPrefix)

	// Initialize commands
	commands = initCommands()

	// Initialize discordgo
	discord, err := discordgo.New("Bot " + cfg.DiscordAuthToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// message handler
	discord.AddHandler(parseChatMessage)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Sunbot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()

}

// Called any time a message is sent
func parseChatMessage(discordSession *discordgo.Session, msgEvent *discordgo.MessageCreate) {

	if len(msgEvent.Content) == 0 {
		DebugPrint("Message received; did not contain text.")
		return
	}

	// Ignore all messages created by the bot itself (or Doritobot) // TODO: separate ignored IDs into database
	if msgEvent.Author.ID == discordSession.State.User.ID || msgEvent.Author.ID == "311737429608628224" {
		return
	}

	// Make it easier to reference message text
	msg := msgEvent.Content

	// Make sure text is actually present to avoid crashing
	if len(msg) > 0 {
		DebugPrint("\nMessage received:\n" + msgEvent.Author.Username + ": " + msg)
	} else {
		DebugPrint("\nMessage received:\n" + msgEvent.Author.Username + ": " + "(file)")

	}

	// Did the message start with the command prefix?
	if msg[:1] == cfg.DefaultPrefix {

		DebugPrint("Message is a command.")

		// prepare variables to parse command
		args := strings.Split(msg[1:], " ")
		cmdInput := args[0]
		args = append(args[:0], args[1:]...)

		if cmd, ok := commands[cmdInput]; ok {
			DebugPrint("Command is valid.")
			discordSession.ChannelMessageSend(msgEvent.ChannelID, cmd.function(args, discordSession).response) // TODO: account for the possibility of a file embed
		} else {
			DebugPrint("Command is not valid.")
			discordSession.ChannelMessageSend(msgEvent.ChannelID, "I don't understand that command.")
		}

		// TODO: implement command usage metrics
	} else {
		DebugPrint("Message is not a command.")

		if cfg.SillyCommandsEnabled {

			switch msg {
			case "h":
				// This is an inside joke.
				// Don't ask, for there isn't an answer.
				discordSession.ChannelMessageSend(msgEvent.ChannelID, "h")
				break
			}

			// since this one dynamically accepts different numbers of letters, it can't be in the switch statement
			if len(msg) >= 3 && strings.ToLower(msg[:3]) == "eee" {
				discordSession.ChannelMessageSend(msgEvent.ChannelID, msg)
			}

		}

		// TODO: implement metrics of standard chat messages
	}

}
