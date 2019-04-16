package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	backends "github.com/bennyz/example-finder/backend"

	"github.com/bennyz/example-finder/persistence"
	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

type client struct {
	*github.Client
	ctx            context.Context
	resultsPerPage int
}

var storage persistence.Storage

// New creats an new instance of the rest client
func New(token string, resultsPerPage int, storageProvider persistence.Storage) (backends.Backend, error) {
	storage = storageProvider
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient := github.NewClient(tc)

	return &client{githubClient, ctx, resultsPerPage}, nil
}

// Search searches code using the github api
func (c *client) Search(query string, lang string) []string {

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: c.resultsPerPage},
	}

	var repos []string
	results, _, err := c.Client.Search.Code(c.ctx, query, opt)
	if err != nil {
		log.Fatal("Failed to fetch results", err)
	}

	for _, result := range results.CodeResults {
		repos = append(repos, result.GetRepository().GetHTMLURL())
		repoID := result.GetRepository().GetID()
		fmt.Printf("Fetching value for: %v", repoID)
		value, err := storage.Get(repoID)
		fmt.Printf("Found value: %s for repoID: %v\n", value, repoID)

		if err != nil || value == "" {
			log.Printf("Value not found in DB, requesting amount of stars...\n")
			repo, _, err := c.Client.Repositories.GetByID(c.ctx, repoID)
			log.Printf("repoID %v has %d\n", repoID, repo.GetStargazersCount())

			if err != nil {
				log.Printf("Failed to fetch repo with ID: %v\n", repoID)
			}

			repoData := backends.Result{
				RepoName: repo.GetName(),
				RepoURL:  repo.GetURL(),
				FileURL:  result.GetPath(),
				Stars:    repo.GetStargazersCount(),
			}

			bytes, err := json.Marshal(&repoData)
			if err != nil {
				log.Fatal(err)
			}
			storage.Save(repoID, bytes)
		}

	}

	return repos
}
