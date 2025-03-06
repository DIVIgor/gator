package main

import "errors"


type command struct {
    name string
    args []string
}

type commands struct {
	cmdList map[string]func(*state, command) error
}


// Register a new handler method for CLI
func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdList[name] = f
}

// Run a CLI method if it exists
func (c *commands) run(s *state, cmd command) (err error) {
	command, exists := c.cmdList[cmd.name]
	if !exists {
		return errors.New("command not found")
	}

	return command(s, cmd)
}