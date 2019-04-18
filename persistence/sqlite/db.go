package sqlite

import (
	"database/sql"
	"log"

	"github.com/bennyz/example-finder/util"

	"github.com/bennyz/example-finder/persistence"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS repo_data(repo_id INTEGER PRIMARY KEY, data JSON1)`
	insertKey   = `INSERT INTO repo_data(repo_id, data) values(?, ?)`
	getKeys     = `SELECT data FROM repo_data WHERE repo_id IN (?)`
)

type sqlite struct {
	*sql.DB
}

// New initializes sqlite database
func New(path string) (persistence.Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTable)

	if err != nil {
		log.Fatal(err)
	}

	return &sqlite{db}, nil
}

func (s *sqlite) Save(key int64, value persistence.JSONValue) (int64, error) {
	tx, err := s.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(insertKey)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(key, value)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	defer stmt.Close()

	return key, nil
}

func (s *sqlite) Get(keys []int64) ([]persistence.JSONValue, error) {
	stmt, err := s.Prepare(getKeys)
	if err != nil {
		log.Fatal(err)
	}

	var value []persistence.JSONValue
	err = stmt.QueryRow(util.SliceToString(keys)).Scan(&value)
	if err == nil {
		log.Printf("Fetched for keys %v", keys)
	}

	defer stmt.Close()

	return value, err
}
