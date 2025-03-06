package main

import (
	"log"
	"os"

	"github.com/DIVIgor/gator/internal/config"
)

// App state
type state struct {
    cfg *config.Config
}


func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config: %w", err)
	}

	appState := &state{cfg: &cfg}

	cmds := commands{cmdList: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
		return
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	
	err = cmds.run(appState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}