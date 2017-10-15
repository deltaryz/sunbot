package main

// generic command struct which contains name, description, and a function
type command struct {
	name        string                // human-readable name of the command
	description string                // description of command's function
	usage       string                // example of how to correctly use command - [] for optional arguments, <> for required arguments
	verbs       []string              // all verbs which are mapped to the same command
	function    func([]string) string // function which receives a slice of arguments and returns a string to display to the user
	// TODO: create 'commandResponse' struct to allow for image embeds, instead of simply printing strings
}

func initCommands() map[string]*command {
	commandList := []*command{}

	// Define all commands here in the order they will be displayed
	commandList = append(commandList,
		&command{
			name:        "Test command",
			description: "A simple command for testing Sunbot.",
			usage:       "test [message]",
			verbs:       []string{"test", "test2"},
			function: func(args []string) string {
				output := "Pong!"
				if len(args) > 0 {
					output += "\nAnd you included a message! Thanks <3"
				}
				return output
			},
		},

		&command{
			name:        "Display help",
			description: "Lists all commands and their purposes.\nCan also display detailed info about a given command.",
			usage:       "help [verb]",
			verbs:       []string{"help", "commands"},
			function: func(args []string) string {
				if len(args) <= 0 {
					output := "**Sunbot " + version + "**\n<https://github.com/techniponi/sunbot>\n\n__Commands:__\n\n"
					for _, cmd := range commandList {
						output += cmd.name + "\n`" + DefaultPrefix + cmd.usage + "`\n"
					}
					return output
				} else {
					// check if command exists
					if cmd, ok := commands[args[0]]; ok {
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
						return output
					} else {
						return "That isn't a valid command."
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
		}
	}

	return commandMap
}
