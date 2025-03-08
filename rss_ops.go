package main

import (
	"context"
	"encoding/xml"
	"errors"
	"html"
	"io"

	"github.com/DIVIgor/gator/internal/requests"
)


func parseXML(rawData []byte) (feed *RSSFeed, err error) {
	err = xml.Unmarshal(rawData, &feed)
	if err != nil {return feed, err}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for idx, item := range feed.Channel.Item {
		feed.Channel.Item[idx].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[idx].Description = html.UnescapeString(item.Description)
	}
	
	return feed, err
}

func fetchFeed(client *requests.Client, ctx context.Context, feedURL string) (feed *RSSFeed, err error) {
	if len(feedURL) == 0 {
		err = errors.New("no URL provided")
		return
	}

	resp, err := client.MakeRequest(ctx, "GET", feedURL, nil)
	if err != nil {return}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {return}

	return parseXML(respData)
}