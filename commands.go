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

	// define commands here
	testCommand := &command{
		name:        "Test Command",
		description: "A simple command for testing Sunbot.",
		usage:       ".test",
		verbs:       []string{"test", "test2"},
		function: func(args []string) string {
			output := "Pong!"
			if len(args) > 0 {
				output += "\nAnd you included a message! Thanks <3"
			}
			return output
		},
	}

	// Commands must be listed here to be usable. Make sure every verb is mapped.
	commandList := map[string]*command{
		// testCommand
		"test":  testCommand,
		"test2": testCommand,
	}

	return commandList
}
