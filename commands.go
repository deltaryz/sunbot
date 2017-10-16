package main

import (
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"os"
)

// generic command struct which contains name, description, and a function
type command struct {
	name        string                                            // human-readable name of the command
	description string                                            // description of command's function
	usage       string                                            // example of how to correctly use command - [] for optional arguments, <> for required arguments
	verbs       []string                                          // all verbs which are mapped to the same command
	function    func([]string, *discordgo.Session) *commandOutput // function which receives a slice of arguments and returns a string to display to the user
}

// output returned by all command functions, can contain a file to be uploaded
type commandOutput struct {
	response string
	file     io.Reader
}

func initCommands() map[string]*command {
	commandList := []*command{}

	commandList = append(commandList,

		// Define all commands here in the order they will be displayed by the help command
		// The 'usage' field should use the default verb
		// Do not include the command prefix

		&command{
			name:        "Display help",
			description: "Lists all commands and their purposes.\nCan also display detailed info about a given command.",
			usage:       "help [verb]",
			verbs:       []string{"help", "commands"},
			function: func(args []string, discordSession *discordgo.Session) *commandOutput {

				DebugPrint("Running help command.")

				if len(args) <= 0 {

					DebugPrint("No arguments; listing commands.")

					output := "**Sunbot " + version + "**\n<https://github.com/techniponi/sunbot>\n\n__Commands:__\n\n"
					for _, cmd := range commandList {
						output += cmd.name + "\n`" + cfg.DefaultPrefix + cmd.usage + "`\n"
					}

					return &commandOutput{response: output}
				} else {

					DebugPrint("Verb was given...")

					// check if command exists
					if cmd, ok := commands[args[0]]; ok {
						DebugPrint("Providing help for given verb.")
						// separated for readability
						output := "**" + cmd.name + "**\n"
						output += cmd.description + "\n\n"
						output += "Usage:\n`" + cfg.DefaultPrefix + cmd.usage + "`\n"
						output += "Verbs:\n"
						// for each verb
						for index, verb := range cmd.verbs {
							// don't add a comma if it's the last one
							if index == (len(cmd.verbs) - 1) {
								output += "`" + cfg.DefaultPrefix + verb + "`"
							} else {
								output += "`" + cfg.DefaultPrefix + verb + "`, "
							}
						}

						return &commandOutput{response: output}
					} else {
						DebugPrint("Given verb was not found.")
						return &commandOutput{response: "That isn't a valid command."}
					}
				}
			},
		},

		&command{
			name:        "Derpibooru search",
			description: "Searches Derpibooru with the given tags as the query, chooses a random result to display.\nUse commas to separate tags like you would on the website.",
			usage:       "derpi <tags>",
			verbs:       []string{"derpi", "db", "derpibooru"},
			function: func(args []string, discordSession *discordgo.Session) *commandOutput {
				if len(args) < 1 {
					DebugPrint("User ran derpibooru command with no tags given.")
					return &commandOutput{response: "Error: no tags specified"}
				} else {
					DebugPrint("User is running derpibooru command...")

					// use derpibooru.go to perform search
					results, err := DerpiSearchWithTags(args[0])
					if err != nil {
						log.Fatal(err)
						return &commandOutput{response: "Error: " + err.Error()}
					}

					// check for results
					if len(results.Search) <= 0 {
						DebugPrint("Derpibooru returned no results.")
						return &commandOutput{response: "Error: no results."}
					} else {
						DebugPrint("Derpibooru returned results; parsed successfully.")
						// pick one randomly
						output := "http:" + results.Search[RandomRange(0, len(results.Search))].Image

						return &commandOutput{response: output}
					}
				}
			},
		},

		&command{
			name:        "Gay",
			description: "Posts a very gay image.",
			usage:       "gay",
			verbs:       []string{"gay"},
			function: func(args []string, discordSession *discordgo.Session) *commandOutput {
				file, err := os.Open("img/gaybats.png")
				if err != nil {
					return &commandOutput{response: "Error opening file"}

				}
				return &commandOutput{file: file}
			},
		},

		&command{
			name:        "Bot source",
			description: "Links to the github Sunbot is hosted on.",
			usage:       "source",
			verbs:       []string{"source", "src"},
			function: func(args []string, discordSession *discordgo.Session) *commandOutput {

				DebugPrint("Running source command.")

				output := "https://github.com/techniponi/sunbot"

				return &commandOutput{response: output}
			},
		},
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
