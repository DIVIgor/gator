package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/DIVIgor/gator/internal/config"
	"github.com/DIVIgor/gator/internal/database"
	"github.com/DIVIgor/gator/internal/requests"
	_ "github.com/lib/pq"
)

// App state
type state struct {
    cfg *config.Config
	client *requests.Client
	db *database.Queries
}


func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config: %w", err)
	}

	// connect to the database
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	webClient := requests.NewClient(15 * time.Second)
	dbQueries := database.New(db)

	appState := &state{
		cfg: &cfg,
		client: &webClient,
		db: dbQueries,
	}

	cmds := commands{cmdList: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	// clearing table command for tests
	cmds.register("reset", handlerReset)

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