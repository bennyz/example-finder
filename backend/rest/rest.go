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
func (c *client) Search(query string, lang string) map[int64]*backends.Result {

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: c.resultsPerPage},
	}

	results, _, err := c.Client.Search.Code(c.ctx, query, opt)
	repos := make(map[int64]*backends.Result)

	if err != nil {
		log.Fatal("Failed to fetch results", err)
	}

	for _, codeResult := range results.CodeResults {
		repoID := codeResult.GetRepository().GetID()
		fmt.Printf("Fetching value for: %v\n", repoID)
		value, err := storage.Get(repoID)
		fmt.Printf("Found value: %s for repoID: %v\n", value, repoID)

		if err != nil || value == nil {
			repo, _, err := c.Client.Repositories.GetByID(c.ctx, repoID)

			if err != nil {
				log.Printf("Could not fetch repo %v \n", repoID)
			}

			log.Printf("Value not found in DB, requesting amount of stars...\n")
			value = handleMissingRepo(repo)
		}
		result := backends.Result{}
		err = json.Unmarshal(value, &result)

		if err != nil {
			log.Fatal(err)
		}

		repos[repoID] = &result
	}

	return repos
}

func handleMissingRepo(repo *github.Repository) []byte {
	repoData := backends.Result{
		RepoName: repo.GetName(),
		RepoURL:  repo.GetURL(),
		Stars:    repo.GetStargazersCount(),
	}

	bytes, err := json.Marshal(&repoData)
	if err != nil {
		log.Fatal(err)
	}

	storage.Save(repo.GetID(), bytes)

	return bytes
}
