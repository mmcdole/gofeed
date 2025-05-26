package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mmcdole/gofeed/v2"
	"github.com/mmcdole/gofeed/v2/atom"
	"github.com/mmcdole/gofeed/v2/rss"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ftest"
	app.Usage = "provide a feed file path or url to parse and print"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "type,t",
			Value: "universal",
			Usage: "type of parser (atom, rss, universal)",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			fmt.Println("Missing feed path or url")
			os.Exit(1)
		}

		feedType := c.String("type")
		feedLoc := c.Args()[0]

		fc, err := fetchFeed(feedLoc)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var feed interface{}

		if strings.EqualFold(feedType, "rss") ||
			strings.EqualFold(feedType, "r") {
			p := rss.Parser{}
			feed, err = p.Parse(strings.NewReader(fc), nil)
		} else if strings.EqualFold(feedType, "atom") ||
			strings.EqualFold(feedType, "a") {
			p := atom.Parser{}
			feed, err = p.Parse(strings.NewReader(fc), nil)
		} else {
			p := gofeed.NewParser()
			feed, err = p.ParseString(fc, nil)
		}

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(feed)
	}
	app.Run(os.Args)
}

func fetchFeed(feedLoc string) (string, error) {
	if strings.HasPrefix(feedLoc, "http") {
		return fetchURL(feedLoc)
	}
	file, err := fetchFile(feedLoc)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func fetchFile(path string) (string, error) {
	f, err := os.ReadFile(path)
	return string(f), err
}

func fetchURL(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
