package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"os/exec"
)

// generic command struct which contains name, description, and a function
type command struct {
	name             string                                                                                          // human-readable name of the command
	description      string                                                                                          // description of command's function
	usage            string                                                                                          // example of how to correctly use command - [] for optional arguments, <> for required arguments
	verbs            []string                                                                                        // all verbs which are mapped to the same command
	requiresDatabase bool                                                                                            // does this command require database access?
	function         func([]string, *discordgo.Channel, *discordgo.MessageCreate, *discordgo.Session) *commandOutput // function which receives a slice of arguments and returns a string to display to the user
}

// output returned by all command functions, can contain a file to be uploaded
type commandOutput struct {
	response string
	file     io.Reader
	embed    *discordgo.MessageEmbed
}

func initCommands() map[string]*command {
	commandList := []*command{}

	commandList = append(commandList,

		// Define all commands here in the order they will be displayed by the help command
		// The 'usage' field should use the default verb
		// Do not include the command prefix

		&command{
			name:             "Display help",
			description:      "Lists all commands and their purposes.\nCan also display detailed info about a given command.",
			usage:            "help [verb]",
			verbs:            []string{"help", "commands"},
			requiresDatabase: false,
			function: func(args []string, channel *discordgo.Channel, msgEvent *discordgo.MessageCreate, discordSession *discordgo.Session) *commandOutput {

				DebugPrint("Running help command.")

				if len(args) <= 0 {

					DebugPrint("No arguments; listing commands.")

					embed := NewEmbed().
						SetTitle("Source").
						SetAuthor("Sunbot " + version).
						//SetDescription("Database enabled: " + strconv.FormatBool(redisEnabled)).
						SetURL("https://github.com/techniponi/sunbot").
						SetImage(discordSession.State.User.AvatarURL("128"))

					for _, cmd := range commandList {
						/*
						if cmd.requiresDatabase && !redisEnabled {
							// Database is not enabled, this command needs it
						} else {
						*/
							embed.AddField(cmd.name, "`"+cfg.DefaultPrefix+cmd.usage+"`")
						//}
					}

					return &commandOutput{embed: embed.MessageEmbed}
				}

				DebugPrint("Verb was given...")

				// check if command exists
				if cmd, ok := commands[args[0]]; ok {

					embed := NewEmbed().
						SetTitle(cmd.name).
						SetDescription(cmd.description).
						AddField("Usage", "`"+cfg.DefaultPrefix+cmd.usage+"`")

					DebugPrint("Providing help for given verb.")

					// compile verbs
					verbOutput := ""
					for index, verb := range cmd.verbs {
						// don't add a comma if it's the last one
						if index == (len(cmd.verbs) - 1) {
							verbOutput += "`" + cfg.DefaultPrefix + verb + "`"
						} else {
							verbOutput += "`" + cfg.DefaultPrefix + verb + "`, "
						}
					}
					embed.AddField("Verbs", verbOutput)
					return &commandOutput{embed: embed.MessageEmbed}
				}
				DebugPrint("Given verb was not found.")
				return &commandOutput{response: "That isn't a valid command."}
			},
		},

		&command{
			name:             "Derpibooru search",
			description:      "Searches Derpibooru with the given tags as the query, chooses a random result to display.\nUse commas to separate tags like you would on the website.",
			usage:            "derpi <tags>",
			verbs:            []string{"derpi", "db", "derpibooru"},
			requiresDatabase: false,
			function: func(args []string, channel *discordgo.Channel, msgEvent *discordgo.MessageCreate, discordSession *discordgo.Session) *commandOutput {
				if len(args) < 1 {
					DebugPrint("User ran derpibooru command with no tags given.")
					return &commandOutput{response: "Error: no tags specified"}
				}
				DebugPrint("User is running derpibooru command...")

				searchQuery := ""

				for _, arg := range args {
					searchQuery += arg + " "
				}

				// enforce 'safe' tag if channel is not nsfw
				if !channel.NSFW {
					DebugPrint("Channel #" + channel.Name + " is SFW, adding safe tag...")
					searchQuery += ",safe"
				}

				DebugPrint("Searching with tags:\n" + searchQuery)

				// use derpibooru.go to perform search
				results, err := DerpiSearchWithTags(searchQuery, cfg.DerpiApiKey)
				if err != nil {
					fmt.Println(err)
					return &commandOutput{response: "Error: " + err.Error()}
				}

				// check for results
				if len(results.Search) <= 0 {
					DebugPrint("Derpibooru returned no results.")
					return &commandOutput{response: "Error: no results."}
				}
				DebugPrint("Derpibooru returned results; parsed successfully.")
				// pick one randomly
				output := "http:" + results.Search[RandomRange(0, len(results.Search))].Image

				return &commandOutput{response: output}
			},
		},

		&command{
			name: "Exec",
			description: "Execute a shell command on my server.\nRequires admin permissions.",
			usage: "exec <command>",
			verbs: []string{"exec"},
			requiresDatabase: false,
			function: func(args []string, channel *discordgo.Channel, msgEvent *discordgo.MessageCreate, discordSession *discordgo.Session) *commandOutput {

				hasPermission := false

				// get user object
				DebugPrint("Getting user object")
				user, err := discordSession.State.Member(channel.GuildID, msgEvent.Author.ID)
				if err != nil {
					fmt.Println(err)
					return &commandOutput{response: "Error getting user from exec command"}
				}

				// get roles from that user
				DebugPrint("Getting user's roles")
				for _, roleID := range user.Roles {
					role, err := discordSession.State.Role(channel.GuildID, roleID)
					if err != nil {
						fmt.Println(err)
						return &commandOutput{response: "Error getting roles from user"}
					}

					DebugPrint("Checking for admin permission")
					if role.Permissions&discordgo.PermissionAdministrator != 0 {
						hasPermission = true
					}
				}

				if hasPermission {

					// convert slice to single string
					fullCommand := ""
					for _, arg := range args {
						fullCommand += arg + " "
					}

					cmd := exec.Command("/bin/bash", "-c", fullCommand)
					stdout, err := cmd.Output()

					if err != nil {
						fmt.Println(err)
						return &commandOutput{response: "Error running command"}
					}

					return &commandOutput{
						response: "```sh\n" + string(stdout) + "\n```",
					}
				}else{
					return &commandOutput{response: "Sorry, but only administrators can use that command."}
				}
			},
		},

		&command{
			name:             "Join",
			description:      "I will join the voice channel of the sender.",
			usage:            "join",
			verbs:            []string{"join"},
			requiresDatabase: false,
			function: func(args []string, channel *discordgo.Channel, msgEvent *discordgo.MessageCreate, discordSession *discordgo.Session) *commandOutput {

				voiceConnection, err := JoinUserVoiceChannel(discordSession, msgEvent.Author.ID)
				if err != nil {
					fmt.Println(err)
					return &commandOutput{response: "Error joining voice channel - are you in one in this server?"}
				}
				fmt.Println("Joined channel:")
				fmt.Println(voiceConnection)

				// TODO: fix audio
/*
				DebugPrint("Setting options")
				options := dca.StdEncodeOptions
				options.RawOutput = true
				options.Bitrate = 96
				options.Application = "lowdelay"

				DebugPrint("Getting info")
				videoInfo, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=gWBZJkfzhNY")
				if err != nil {
					DebugPrint("Error getting info")
					fmt.Println(err)
				}

				DebugPrint("Setting format")
				format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
				downloadURL, err := videoInfo.GetDownloadURL(format)
				if err != nil {
					DebugPrint("Error setting format")
					fmt.Println(err)
				}

				DebugPrint("Encoding file")
				encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
				if err != nil {
					DebugPrint("Error encoding file")
					fmt.Println(err)
				}
				defer encodingSession.Cleanup()

				DebugPrint("Initiating stream")
				done := make(chan error)
				dca.NewStream(encodingSession, voiceConnection, done)
				playErr := <- done
				if playErr != nil && playErr != io.EOF {
					fmt.Println(playErr)
				}
*/
				return &commandOutput{response: "Joining your channel!"}
			},
		},

		&command{
			name:             "Leave",
			description:      "I will leave the voice channel.",
			usage:            "leave",
			verbs:            []string{"leave", "disconnect", "quit", "exit"},
			requiresDatabase: false,
			function: func(args []string, channel *discordgo.Channel, msgEvent *discordgo.MessageCreate, discordSession *discordgo.Session) *commandOutput {

				for _, voiceChannel := range discordSession.VoiceConnections {
					DebugPrint("Looking for voice channel in this guild...")
					if voiceChannel.GuildID == channel.GuildID {
						DebugPrint("Found a channel in this guild!")
						voiceChannel.Disconnect()
						return &commandOutput{response: "Bye!"}
					}
				}
				return &commandOutput{response: ""}
			},
		},

		&command{
			name:             "Gay",
			description:      "Posts a very gay image.",
			usage:            "gay",
			verbs:            []string{"gay"},
			requiresDatabase: false,
			function: func(args []string, channel *discordgo.Channel, msgEvent *discordgo.MessageCreate, discordSession *discordgo.Session) *commandOutput {
				file, err := os.Open("img/gaybats.png") // TODO: move this to database; allow users to add images (permission system?)
				if err != nil {
					fmt.Println(err)
					return &commandOutput{response: "Error opening file"}

				}
				return &commandOutput{file: file}
			},
		},

	/*
		&command{
			name:             "User stats",
			description:      "Displays the statistics of the user.",
			usage:            "stats [user]", // TODO: implement pinging users
			verbs:            []string{"stats"},
			requiresDatabase: true,
			function: func(args []string, channel *discordgo.Channel, msgEvent *discordgo.MessageCreate, discordSession *discordgo.Session) *commandOutput {

				if len(args) > 0 {
					if len(msgEvent.Mentions) > 0 {
						// User tagged someone else
						taggedUser := msgEvent.Mentions[0] // only the first one

						userDb, err := GetUser(taggedUser, false)
						if err != nil {
							return &commandOutput{response: "That user doesn't exist in the database yet. They need to chat some!"}
						}

						posts := userDb.Val()["posts"]
						return &commandOutput{response: taggedUser.Username + " has made " + posts + " posts!"} // TODO: format as embed, show more values
					}
					// user didn't tag anyone
					// TODO: accept aliases as well as mentions
					return &commandOutput{response: "To see someone's stats, tag the person directly!"}
				}
				// User's own stats
				userDb, err := GetUser(msgEvent.Author, false)
				if err != nil {
					return &commandOutput{response: "You don't exist in the database yet. You need to chat some!"}
				}
				posts := userDb.Val()["posts"]
				return &commandOutput{response: "You have made " + posts + " posts!"} // TODO: format as embed, show more values
			},
		},
	*/
	)

	// Map for matching verbs to commands
	commandMap := make(map[string]*command)

	// Loop through commandList to get each verb
	for _, cmd := range commandList {
		for _, verb := range cmd.verbs {
			commandMap[verb] = cmd
			DebugPrint("Mapped '" + verb + "' to '" + cmd.name + "'")
		}
	}

	return commandMap
}
