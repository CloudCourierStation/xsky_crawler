package main

import (
	"os"
	"xsky_crawler/cmd/crawler"
)

func main() {
	if err := crawler.Run(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
