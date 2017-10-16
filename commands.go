package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// generic command struct which contains name, description, and a function
type command struct {
	name        string                                            // human-readable name of the command
	description string                                            // description of command's function
	usage       string                                            // example of how to correctly use command - [] for optional arguments, <> for required arguments
	verbs       []string                                          // all verbs which are mapped to the same command
	function    func([]string, *discordgo.Session) *commandOutput // function which receives a slice of arguments and returns a string to display to the user
}

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
			name:        "Test command",
			description: "A simple command for testing Sunbot.",
			usage:       "test [message]",
			verbs:       []string{"test", "test2"},
			function: func(args []string, discordSession *discordgo.Session) *commandOutput {

				DebugPrint("Running test command.")

				output := "Pong!"
				if len(args) > 0 {
					DebugPrint("Message was included.")
					output += "\nAnd you included a message! Thanks <3"
				}

				return &commandOutput{response: output}
			},
		},

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
						output += cmd.name + "\n`" + DefaultPrefix + cmd.usage + "`\n"
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
						output += "Usage:\n`" + DefaultPrefix + cmd.usage + "`\n"
						output += "Verbs:\n"
						// for each verb
						for index, verb := range cmd.verbs {
							// don't add a comma if it's the last one
							if index == (len(cmd.verbs) - 1) {
								output += "`" + DefaultPrefix + verb + "`"
							} else {
								output += "`" + DefaultPrefix + verb + "`, "
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
			description: "Searches Derpibooru with the given tags as the query, chooses a random result to display.\nUse commas to separate like you would on the website.",
			usage:       "derpi <tags>",
			verbs:       []string{"derpi", "db", "derpibooru"},
			function: func(args []string, discordSession *discordgo.Session) *commandOutput {
				if len(args) < 1 {
					DebugPrint("User ran derpibooru command with no tags given.")
					return &commandOutput{response: "Error: no tags specified"}
				} else {
					DebugPrint("User is running derpibooru command...")

					// format for URL query
					derpiTags := strings.Replace(args[0], " ", "+", -1)

					// make URL query
					resp, err := http.Get("https://derpibooru.org/search.json?q=safe," + derpiTags)
					if err != nil {
						DebugPrint("Failed with HTTP error.")
						log.Fatal(err)
						return &commandOutput{response: "Failed with HTTP error."}
					}

					// read response body
					defer resp.Body.Close()
					respBody, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						DebugPrint("Failed with error reading response body.")
						log.Fatal(err)
						return &commandOutput{response: "Failed with error reading response body."}
					}

					// parse json
					results := DerpiResults{}
					err = json.Unmarshal(respBody, &results)
					if err != nil {
						DebugPrint("Failed with JSON parsing error.")
						log.Fatal(err)
						return &commandOutput{response: "Failed with JSON parsing error."}
					}

					// check for results
					if len(results.Search) <= 0 {
						DebugPrint("Derpibooru returned no results.")
						return &commandOutput{response: "Error: no results."}
					} else {

						// pick one randomly
						output := "http:" + results.Search[RandomRange(0, len(results.Search))].Image

						return &commandOutput{response: output}

					}

				}
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
