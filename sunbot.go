package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
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
	RedisURL             string `env:"REDIS_URL" envDefault:""`          // environment variable REDIS_URL
	RedisPassword        string `env:"REDIS_PASSWORD" envDefault:""`     // environment variable REDIS_PASSWORD
}

// Global variables
var (
	commands     map[string]*command // verb string -> command object (see commands.go)
	cfg          config
	client       *redis.Client
	redisEnabled bool
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

	// init environment cariables
	cfg = config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Error processing environment variables.\nPlease check https://github.com/techniponi/sunbot for details.\n\n" + err.Error())
		return
	}

	if cfg.RedisURL != "" {
		// init redis
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisURL,
			Password: cfg.RedisPassword,
			DB:       0, // use default DB
		})

		pong, err := client.Ping().Result()
		fmt.Println("Connecting to Redis..."+pong, err) // Output: PONG <nil>

		// tell user if db didn't connect
		if err != nil {
			fmt.Println("Error connecting to Redis.")
		}

		// for later checks; less typing than checking if cfg.RedisURL is empty
		redisEnabled = true
	} else {
		redisEnabled = false
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
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

	// Ignore all messages created by any bot (including itself)
	if msgEvent.Author.Bot {
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
			output := cmd.function(args, discordSession)
			if output.file == nil {
				discordSession.ChannelMessageSend(msgEvent.ChannelID, output.response)
			} else {
				DebugPrint("Response contains image, uploading now")
				discordSession.ChannelFileSend(msgEvent.ChannelID, "image.png", output.file)
			}
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

		if redisEnabled {

			DebugPrint("Checking if user exists in database...")
			// check if user is in database yet
			userDb := client.HGetAll("user:" + msgEvent.Author.ID)
			_, err := userDb.Result()
			if err != nil {
				DebugPrint("Redis error:")
				fmt.Println(err)

				if err == redis.Nil {
					DebugPrint("User does not exist. Adding user to database.")
					client.HMSet("user:"+msgEvent.Author.ID, map[string]interface{}{
						"username": msgEvent.Author.Username,
						"isBot":    msgEvent.Author.Bot,
					})
				} else {
					fmt.Println("Database error, see log.")
				}

			} else {
				DebugPrint("User does exist.")
				DebugPrint(userDb.String())
				DebugPrint(userDb.Val()["username"]) // TODO: remove this (reference for later)

				// TODO: separate all database transactions to a separate go file/api which handles missing values
			}

			// TODO: implement metrics of standard chat messages
		}

	}

}
