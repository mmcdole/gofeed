package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/rss"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("ftest", flag.ContinueOnError)
	// main reports errors once; help is written to the normal output stream.
	flags.SetOutput(io.Discard)
	flags.Usage = func() {
		fmt.Fprintln(out, "Usage: ftest [--type atom|rss|universal] <feed file or URL>")
		fmt.Fprintln(out, "  --type, -t  type of parser (default universal)")
		fmt.Fprintln(out, "  --help, -h  show help")
	}
	var feedType string
	flags.StringVar(&feedType, "type", "universal", "type of parser")
	flags.StringVar(&feedType, "t", "universal", "type of parser")
	if len(args) == 1 && args[0] == "help" {
		flags.Usage()
		return nil
	}
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() == 0 {
		return errors.New("missing feed path or URL")
	}

	fc, err := fetchFeed(flags.Arg(0))
	if err != nil {
		return err
	}

	var feed interface{}
	switch strings.ToLower(feedType) {
	case "rss", "r":
		feed, err = (&rss.Parser{}).Parse(strings.NewReader(fc))
	case "atom", "a":
		feed, err = (&atom.Parser{}).Parse(strings.NewReader(fc))
	default:
		feed, err = gofeed.NewParser().ParseString(fc)
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, feed)
	return err
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
