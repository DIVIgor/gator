# Blog AggreGATOR 🐊

**Gator**🐊 is a simple Go-based CLI blog aggregator that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

*Note that there's no secure authentication system available (at least for now). This means that user registration is implemented by simply setting a username.*

## Table of contents

- [Requirements](#requirements)
- [Configuration](#configuration)
- [Installation Guide](#installation-guide)
- [Usage Guide](#usage-guide)
  - [General Info](#general-info)
  - [Command List](#command-list)

## Requirements

Any unix-based OS: Linux (or WSL for Windows) / MacOS.
To build/install the app you need to have on your computer:

- **Go** version 1.23+
- **PostgreSQL** version 16.8+

## Configuration

Gator requires `.gatorconfig.json` - a JSON configuration file at your home directory with the following structure:

`
{
  "db_url": "db_connection_string"
}
`

## Installation Guide

You can either build or install Gator on your computer:

- to build the app, navigate to the gator root folder and use `go build` command. This will compile the app to a single executable file.
- to install the app use `go install` command from within the gator root folder. This will compile and install the app globally on your system. Now Gator will be accessible by `gator` command in your CLI.

## Usage Guide

Once you've set a config file, you can run or build/install the app.

### General Info

*Note that the usage depends on whether the app is installed or not.*

- If the app is not installed, navigate to the gator root folder and run:
`go run . <command> [argument]`
- If the app is installed/builded, use the following tmeplate
`gator <command> [argument]`

### Command List

- `register <user name>` - register a user by user name
- `login <user name>` - login as a user by user name (should be registered)
- `users` - show a list of registered users
- `agg` - aggregate data for the feeds followed by the current user
- `addfeed <feed name> <URL>` - add a new RSS feed (automatically marks as following by the current user)
- `feeds` - show a full list of saved feeds
- `follow <URL>` - follow feed by URL for the current user
- `unfollow <URL>` - unfollow feed for the current user
- `following` - show a list of following feeds for the current user
- `browse [number of entries]` - show a list of following feeds (2 by default) for the current user, starting from the most recently updated entries
