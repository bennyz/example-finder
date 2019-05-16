package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/bennyz/example-finder/util"

	backends "github.com/bennyz/example-finder/backend"

	"github.com/bennyz/example-finder/persistence"
	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

// ClientOptions holds the options with which to create the REST client
type ClientOptions struct {
	ResultsPerPage int
	Lang           string
	RefreshDB      bool
}

type client struct {
	*github.Client
	ctx context.Context
	*ClientOptions
}

var (
	storage   persistence.Storage
	repos     = make(map[int64]*backends.Result)
	repoPaths = make(map[int64][]string)
)

// New creats an new instance of the rest client
func New(token string, co *ClientOptions, storageProvider persistence.Storage) (backends.Backend, error) {
	storage = storageProvider
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient := github.NewClient(tc)

	return &client{githubClient, ctx, co}, nil
}

// Search searches code using the github api
func (c *client) Search(query, lang string) []*backends.Result {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: c.ResultsPerPage},
	}

	if c.RefreshDB {
		storage.Clear()
	}

	if lang != "" {
		query = fmt.Sprint("language:", lang, " ", query)
	}

	results, _, err := c.Client.Search.Code(c.ctx, query, opt)
	if err != nil {
		log.Fatal("Failed to fetch results", err)
	}

	// Collect all repo IDs from the results
	var repoIDs []int64
	for _, codeResult := range results.CodeResults {
		repoID := codeResult.GetRepository().GetID()
		repoIDs = append(repoIDs, repoID)
		repoPaths[repoID] = append(repoPaths[repoID], codeResult.GetHTMLURL())
	}

	cachedRepoIds := make([]int64, 0, 0)
	if !c.RefreshDB {
		values, err := storage.Get(repoIDs)
		if err != nil {
			fmt.Println(err)
		}

		cachedRepoIds = initializeCache(values)
	}

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

	reposSlice := util.MapToSlice(repos)
	sort.Sort(backends.ByStars(reposSlice))
	return reposSlice
}

func handleMissingRepo(repo *github.Repository) []byte {
	repoData := backends.Result{
		RepoID:    repo.GetID(),
		RepoName:  repo.GetName(),
		RepoURL:   repo.GetHTMLURL(),
		Stars:     repo.GetStargazersCount(),
		FilePaths: repoPaths[repo.GetID()],
	}

	bytes, err := json.Marshal(&repoData)
	if err != nil {
		log.Fatal(err)
	}

	storage.Save(repo.GetID(), bytes)

	return bytes
}

func initializeCache(values []persistence.JSONValue) []int64 {
	var results []int64
	for _, value := range values {
		result := backends.Result{}
		err := json.Unmarshal(value, &result)
		if err != nil {
			log.Println(err)
		}

		result.FilePaths = repoPaths[result.RepoID]
		results = append(results, result.RepoID)
		repos[result.RepoID] = &result
	}

	return results
}
