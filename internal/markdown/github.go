package markdown

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/go-github/v47/github"
	"github.com/gookit/config/v2"
	"github.com/jomei/notionapi"
	"github.com/joshmenden/n2h/internal/log"

	"golang.org/x/oauth2"
)

func getClient() *github.Client {
	ctx := context.Background()
	githubPAT := config.String("githubPAT")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubPAT},
	)

	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func authenticatedGithubUsername() (*string, error) {
	user, _, err := getClient().Users.Get(context.Background(), "")
	if err != nil {
		return nil, err
	}

	return user.Login, nil
}

func createOrUpdateGist(filename string, codeblock *notionapi.CodeBlock) (gist *github.Gist, err error) {
	foundGist, err := findGist(filename)
	if err != nil {
		return
	}

	if foundGist == nil {
		gist, err = postGist(filename, codeblock)
		log.SubStatus(fmt.Sprintf("%s %s", log.Colorify("creating", "green"), log.Linkify("new gist", *gist.HTMLURL)))
	} else {
		gist, err = editGist(*foundGist.ID, filename, codeblock)
		log.SubStatus(fmt.Sprintf("%s %s", log.Colorify("updating", "blue"), log.Linkify("existing gist", *gist.HTMLURL)))
	}

	return
}

func editGist(gistID, filename string, codeblock *notionapi.CodeBlock) (*github.Gist, error) {
	gist, _, err := getClient().Gists.Edit(context.Background(), gistID, createGist(filename, codeblock))
	if err != nil {
		return nil, err
	}

	return gist, nil
}

func findGist(filename string) (foundGist *github.Gist, err error) {
	// calling with an empty string in the user param returns gists for authenticated user
	// see https://github.com/google/go-github/blob/af69917d404934f24ae81ffaf83fbc1db074967d/github/gists.go#L18
	gists, _, err := getClient().Gists.List(context.Background(), "", &github.GistListOptions{})
	for _, gist := range gists {
		if _, ok := gist.Files[github.GistFilename(filename)]; ok {
			if foundGist == nil {
				foundGist = gist
			}
		}
	}

	return
}

func postGist(filename string, codeblock *notionapi.CodeBlock) (*github.Gist, error) {
	gist, _, err := getClient().Gists.Create(context.Background(), createGist(filename, codeblock))

	return gist, err
}

func createGist(filename string, codeblock *notionapi.CodeBlock) *github.Gist {
	return &github.Gist{
		Public: aws.Bool(true),
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(filename): {
				Filename: &filename,
				Language: &codeblock.Code.Language,
				Content:  &codeblock.Code.RichText[0].Text.Content,
			},
		},
	}
}
