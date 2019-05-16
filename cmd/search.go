package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/bennyz/example-finder/backend/rest"
	"github.com/bennyz/example-finder/persistence"
	"github.com/bennyz/example-finder/persistence/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	token          string
	lang           string
	mode           string
	dbPath         string
	resultsPerPage int
	refreshDB      bool
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "this is the main command",
	// TODO: add a better description
	Long: `No description at the moment`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not enough arguments")
		}

		if args[0] == "" {
			return errors.New("Please add a string to search for")
		}

		if token == "" && readTokenFromFile() == "" {
			return errors.New("Please provide a token, either with -t or put it in a .token file")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch mode {
		case "rest":
			co := rest.ClientOptions{
				ResultsPerPage: resultsPerPage,
				Lang:           lang,
				RefreshDB:      refreshDB,
			}
			// TODO: handle errors and stuff
			client, _ := rest.New(token, &co, initDb())
			results := client.Search(args[0], lang)

			fmt.Println()
			for _, r := range results {
				fmt.Printf("repo: %s \n\t%s\n\tstars: %d\n", r.RepoName, r.RepoURL, r.Stars)
				fmt.Printf("\tfiles:\n")
				// TODO: extract this
				for _, path := range r.FilePaths {
					fmt.Printf("\t %v\n", path)
				}
			}
		}
	},
}

func init() {
	viper.SetDefault("token_file", ".token")

	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&token, "token", "t", "", "string to search")
	searchCmd.Flags().StringVarP(&lang, "lang", "l", "", "language")
	searchCmd.Flags().StringVarP(&mode, "mode", "m", "rest", "search backend")
	searchCmd.Flags().StringVarP(&dbPath, "db", "", "db.sqlite", "database file")
	searchCmd.Flags().IntVarP(&resultsPerPage, "results", "r", 30, "results per page")
	searchCmd.Flags().BoolVarP(&refreshDB, "refresh-db", "", false, "refresh database")
}

func readTokenFromFile() string {
	s, err := ioutil.ReadFile(viper.GetString("token_file"))
	if err != nil {
		log.Fatal(err)
	}

	token = strings.TrimSuffix(string(s), "\n")

	return token
}

func initDb() persistence.Storage {
	storage, err := sqlite.New("sqlite.db")
	if err != nil {
		log.Fatal(err)
	}

	return storage
}
