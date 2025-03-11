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

// Fetch feed by URL
func handlerAgg(s *state, cmd command) (err error) {
    feed, err := fetchFeed(s.client, context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {return}
    fmt.Print(feed)

    return err
}

// Add feed to DB
func handlerAddFeed(s *state, cmd command, user database.User) (err error) {
    if len(cmd.args) < 2 {
        return fmt.Errorf("%s has not enough arguments", cmd.name)
    }

    feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
        Name: cmd.args[0],
        Url: cmd.args[1],
        UserID: user.ID,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    })
    if err != nil {
        return fmt.Errorf("couldn't create feed: %w", err)
    }

    _, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        UserID: user.ID,
        FeedID: feed.ID,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    })
    if err != nil {
        return fmt.Errorf("couldn't create feed follow: %w", err)
    }

    printFeed(feed.ID, feed.Name, feed.Url, user.Name, feed.CreatedAt, feed.UpdatedAt)
    fmt.Println("==============================")

    return err
}

// Takes specific values instead of struct to reduce number of DB calls
func printFeed(id int32, name, url, username string, created, updated time.Time) {
    fmt.Println("* ID:", id)
	fmt.Println("* Created:", created)
	fmt.Println("* Updated:", updated)
	fmt.Println("* Name:", name)
	fmt.Println("* URL:", url)
    fmt.Println("* User:", username)
}

// Get all the feeds from DB
func handlerGetFeeds(s *state, cmd command) (err error) {
    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {
        return fmt.Errorf("couldn't get feeds: %w", err)
    }

    if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return err
	}

    fmt.Printf("Found %d feeds:\n", len(feeds))
    for _, feed := range feeds {
        printFeed(feed.ID, feed.Name, feed.Url, feed.Username, feed.CreatedAt, feed.UpdatedAt)
        fmt.Println("==============================")
    }

    return err
}

// Create a new feed follow record for the current user
func handlerFollow(s *state, cmd command, user database.User) (err error) {
    if len(cmd.args) < 1 {
        return fmt.Errorf("%s has not enough arguments", cmd.name)
    }

    feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
    if err != nil {return}

    followedFeed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        UserID: user.ID,
        FeedID: feed.ID,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    })

    fmt.Println("You are now following:")
    fmt.Println("Feed:", followedFeed.FeedName)
    fmt.Println("User:", followedFeed.UserName)

    return err
}

// Print all the names of the feeds the current user is following
func handlerFollowing(s *state, cmd command, user database.User) (err error) {
    followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {return}

    if len(followedFeeds) == 0 {
        fmt.Println("No following feeds.")
        return err
    }

    fmt.Println("Your followed feeds:")
    for _, feed := range followedFeeds {
        fmt.Println("*", feed.FeedName)
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