package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DIVIgor/gator/internal/database"
	"github.com/google/uuid"
)


const outputTimeFormat string = "02-Jan-2006 at 15:04"
const printDelimiter string = "=============================================================================================================="


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
func handlerAgg(s *state, cmd command, user database.User) (err error) {
    if len(cmd.args) < 1 {
        return fmt.Errorf("%s has not enough arguments", cmd.name)
    }

    requestDelay, err := time.ParseDuration(cmd.args[0])
    if err != nil {
        return fmt.Errorf("invalid duration: %w", err)
    }

    log.Printf("Collecting feeds every %s", requestDelay)

    ticker := time.NewTicker(requestDelay)
    for ; ; <-ticker.C {
        scrapeFeeds(s, user)
    }
}

// Parse scraped timestamp from feed entries
func parseTime(timeStr string) (parsedTime time.Time, err error) {
    // timestamps
    tsLayouts := []string{
        time.RFC1123,
        time.RFC822,
        time.RFC3339,
        "2006-01-02T15:04:05",
    }

    for _, ts := range tsLayouts {
        parsedTime, err = time.Parse(ts, timeStr)
        if err == nil {
            return parsedTime, err
        }
    }

    return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// Scrape a feed and print its details
func scrapeFeed(s *state, nextFeed database.GetNextToFetchRow) {
    // mark the feed as fetched or update fetched time
    err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
    if err != nil {
        log.Printf("Couldn't mark feed %s fetched: %v", nextFeed.Name, err)
        return
    }

    // fetch the feed
    feed, err := fetchFeed(s.client, context.Background(), nextFeed.Url)
    if err != nil {
        log.Printf("Couldn't collect feed %s: %v", nextFeed.Name, err)
        return
    }

    log.Printf("Feed %s collected. Found %v posts.", feed.Channel.Title, len(feed.Channel.Item))
    for _, el := range feed.Channel.Item {
        parsedTime, err := parseTime(el.PubDate)
        if err != nil {return}

        _, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
            Title: el.Title,
            Url: el.Link,
            Description: sql.NullString{
                String: el.Description,
                Valid: true,
            },
            PublishedAt: parsedTime,  // probably may be NULL
            FeedID: nextFeed.ID,
            CreatedAt: time.Now().UTC(),
            UpdatedAt: time.Now().UTC(),
        })

        if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
            continue
        }
        if err != nil {
            log.Println(err)
            return
        }
    }
}

// Browse saved posts
func handlerBrowsePosts(s *state, cmd command, user database.User) (err error) {
    postLimit := 2
    if len(cmd.args) > 0 {
        postLimit, err = strconv.Atoi(cmd.args[0])
        if err != nil {
            log.Println(err)
            return fmt.Errorf("invalid limit: %w", err)
        }
    }

    posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
        UserID: user.ID,
        Limit: int32(postLimit),
    })
    if err != nil {
        return fmt.Errorf("couldn't get posts for user: %w", err)
    }

    for _, post := range posts {
        fmt.Println("Title:", post.Title, "|", "Published at:", post.PublishedAt.Format(outputTimeFormat))
        if post.Description.Valid {
            fmt.Println("Description")
            fmt.Println(post.Description.String)
        }
        fmt.Println(printDelimiter)
    }

    return err
}

// Feed aggregation
func scrapeFeeds(s *state, user database.User) {
    // get the latest/unfetched feed
    nextFeed, err := s.db.GetNextToFetch(context.Background(), user.ID)
    if err != nil {
        log.Println("Couldn't get next feeds to fetch", err)
        return
    }

    scrapeFeed(s, nextFeed)
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

    printFeed(feed)
    fmt.Println(printDelimiter)

    return err
}

// Print general info on a given feed
func printFeed(feed database.Feed) {
    fmt.Println("* ID:", feed.ID)
	fmt.Println("* Created:", feed.CreatedAt.Format(outputTimeFormat))
	fmt.Println("* Updated:", feed.UpdatedAt.Format(outputTimeFormat))
	fmt.Println("* Name:", feed.Name)
	fmt.Println("* URL:", feed.Url)
    fmt.Println("* Last Fetched At:", feed.LastFetchedAt.Time.Format(outputTimeFormat))
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
        printFeed(feed)
        fmt.Println(printDelimiter)
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

func handlerUnfollow(s *state, cmd command, user database.User) (err error) {
    if len(cmd.args) < 1 {
        return fmt.Errorf("%s has not enough arguments", cmd.name)
    }

    return s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
        UserID: user.ID,
        Url: cmd.args[0],
    })
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