package main

import (
	"fmt"
)

// set current user using an argument from CLI
func handlerLogin(s *state, cmd command) (err error) {
    if len(cmd.args) == 0 {
        return fmt.Errorf("[Usage Error] %s has not enough arguments", cmd.name)
    }

    err = s.cfg.SetUser(cmd.args[0])
    if err != nil {
        return fmt.Errorf("[Usage Error] couldn't set a user %w", err)
    }

    fmt.Println("The user has been successfully set to", cmd.args[0])
    return err
}