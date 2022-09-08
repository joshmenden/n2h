package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/joshmenden/n2h/internal/draft"
	"github.com/joshmenden/n2h/internal/log"
	"github.com/joshmenden/n2h/internal/markdown"
	"github.com/joshmenden/n2h/internal/notionfind"
)

var PWD string

func main() {
	pageTitle := flag.String("title", "", "The complete title of the Notion article you're ready to publish")
	statusProperty := flag.String("status-property", "Status", "The name of the `status` property to look for a specific article in")
	status := flag.String("status", "", "The value of the status in which the desired article is found")

	flag.Parse()

	if *pageTitle == "" || *status == "" {
		log.Error(errors.New("missing required flags `status` or `pageTitle`"))
	}

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yamlv3.Driver)

	err := config.LoadFiles(PWD + "/secrets.yml")
	if err != nil {
		log.Error(err)
	}

	log.Status("fetching content from Notion", "🏇")
	page, blocks, err := notionfind.Content(pageTitle, statusProperty, status)
	if err != nil {
		log.Error(err)
	}

	log.Status("building markdown from blocks", "🔨")
	content, err := markdown.Build(page, blocks)
	if err != nil {
		log.Error(err)
	}

	log.Status("creating hugo blog template with markdown", "🌍")
	ok, err := draft.GenerateAndSaveDraftFile(&draft.Info{
		Title:       markdown.TitleFromPage(page),
		Filename:    fmt.Sprintf("%s.md", markdown.ToSnakeCase(markdown.TitleFromPage(page))),
		Description: nil,
		Draft:       aws.Bool(true),
		Content:     *content,
	}, PWD)

	if !ok {
		log.Error(err)
	}

	log.Status("draft created", "✅")
}
