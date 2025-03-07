package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DIVIgor/gator/internal/database"
	"github.com/google/uuid"
)

// set current user using an argument from CLI
func handlerLogin(s *state, cmd command) (err error) {
    if len(cmd.args) == 0 {
        return fmt.Errorf("%s has not enough arguments", cmd.name)
    }

    usr, err := s.db.GetUser(context.Background(), cmd.args[0])
    if err != nil {
        return errors.New("user does not exist")
    }


    err = s.cfg.SetUser(usr.Name)
    if err != nil {
        return fmt.Errorf("couldn't set current user %w", err)
    }

    fmt.Println("The user has been successfully set to", cmd.args[0])
    return err
}

// create a user
func handlerRegister(s *state, cmd command) (err error) {
    if len(cmd.args) == 0 {
        return fmt.Errorf("%s has not enough arguments", cmd.name)
    }

    user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: cmd.args[0],
    })
    if err != nil {
        return fmt.Errorf("couldn't create user: %w", err)
    }

    err = s.cfg.SetUser(user.Name)
    if err != nil {
        return fmt.Errorf("couldn't set current user %w", err)
    }

    fmt.Println("User", user.Name, "successfully created")

    return err
}

// get and print all the registered users from the DB
func handlerUsers(s *state, cmd command) (err error) {
    userlist, err := s.db.GetUsers(context.Background())
    if err != nil {return}

    if len(userlist) == 0 {
        fmt.Println("No registered users")
        return
    }

    fmt.Println("Registered users:")
    for _, user := range userlist {
        if user.Name == s.cfg.User {
            fmt.Println("*", user.Name, "(current)")
            continue
        }
        fmt.Println("*", user.Name)
    }

    return err
}

// reset users table FOR TEST PURPOSES
func handlerReset(s *state, cmd command) (err error) {
    err = s.db.ClearUsers(context.Background())
    if err != nil {return}

    fmt.Println("Database was successfully reset")

    return err
}