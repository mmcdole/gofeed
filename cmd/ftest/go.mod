module github.com/mmcdole/gofeed/cmd/ftest

go 1.23.0

toolchain go1.24.3

require (
	github.com/mmcdole/gofeed/v2 v2.0.0
	github.com/urfave/cli v1.22.16
)

require (
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/mmcdole/goxpp v1.1.1-0.20240225020742-a0c311522b23 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

replace github.com/mmcdole/gofeed/v2 => ../..
