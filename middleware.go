package main

import (
	"context"

	"github.com/DIVIgor/gator/internal/database"
)

// Get the current user for handlers
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) (err error) {
		user, err := s.db.GetUser(context.Background(), s.cfg.User)
		if err != nil {return}

		return handler(s, cmd, user)
	}
}