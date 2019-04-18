package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bennyz/example-finder/util"

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

var repos map[int64]*backends.Result

// Search searches code using the github api
func (c *client) Search(query, lang string) map[int64]*backends.Result {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: c.resultsPerPage},
	}

	results, _, err := c.Client.Search.Code(c.ctx, query, opt)
	if err != nil {
		log.Fatal("Failed to fetch results", err)
	}

	var repoIDs []int64
	for _, codeResult := range results.CodeResults {
		repoIDs = append(repoIDs, codeResult.GetRepository().GetID())
	}

	values, err := storage.Get(repoIDs)
	if err != nil {
		fmt.Println(err)
	}

	cachedRepoIds := initializeCache(values)
	missingRepoIDs := util.Difference(cachedRepoIds, repoIDs)
	var missingRepos []*backends.Result
	for _, repoID := range missingRepoIDs {
		if _, ok := repos[repoID]; !ok {
			repo, _, err := c.Client.Repositories.GetByID(c.ctx, repoID)
			if err != nil {
				log.Printf("Could not fetch repo %v \n", repoID)
			}
			fmt.Printf("Fetched missing repo %v\n", repo.GetName())

			value := handleMissingRepo(repo)
			missingRepoResult := backends.Result{}
			err = json.Unmarshal(value, &missingRepoResult)
			missingRepos = append(missingRepos, &missingRepoResult)
			repos[repoID] = &missingRepoResult
		}
	}

	return repos
}

func handleMissingRepo(repo *github.Repository) []byte {
	repoData := backends.Result{
		RepoID:   repo.GetID(),
		RepoName: repo.GetName(),
		RepoURL:  repo.GetHTMLURL(),
		Stars:    repo.GetStargazersCount(),
	}

	bytes, err := json.Marshal(&repoData)
	if err != nil {
		log.Fatal(err)
	}

	storage.Save(repo.GetID(), bytes)

	return bytes
}

func initializeCache(values []persistence.JSONValue) []int64 {
	repos = make(map[int64]*backends.Result)

	var results []int64
	for _, value := range values {
		result := backends.Result{}
		err := json.Unmarshal(value, &result)
		if err != nil {
			log.Println(err)
		}

		results = append(results, result.RepoID)
		repos[result.RepoID] = &result
	}

	return results
}
