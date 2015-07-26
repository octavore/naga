package service

import (
	"fmt"
	"os"
)

// Command represents a command-line keyword for the app.
// This is then typically invoked as follows:
//   ./myapp <keyword>
type Command struct {
	Keyword    string
	Run        func(*CommandContext)
	ShortUsage string
	Usage      string
}

// AddCommand adds a command to the service via its Config.
func (c *Config) AddCommand(cmd *Command) {
	c.service.registerCommand(cmd)
}

// CommandContext is passed to the command when it is run,
// containing an array of parsed arguments.
type CommandContext struct {
	cmd  *Command
	Args []string
}

// UsageExit prints the usage for the executed command and exits.
func (c *CommandContext) UsageExit() {
	fmt.Println(c.cmd.Keyword)
	fmt.Println(c.cmd.Usage)
	os.Exit(0)
}

// RequireAtLeastNArgs is a helper function to ensure we have at least n args.
func (c *CommandContext) RequireAtLeastNArgs(n int) {
	if len(c.Args) < n {
		c.UsageExit()
	}
}

// RequireExactlyNArgs is a helper function to ensure we have at exactly n args.
func (c *CommandContext) RequireExactlyNArgs(n int) {
	if len(c.Args) != n {
		c.UsageExit()
	}
}
