package crawler

import (
	"os"
	cr "xsky_crawler/cmd/crawler/campus_recruitment"
)

func Run() error {
	cmd := cr.NewCrawlerCommand(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
