package sqlite

import (
	"database/sql"
	"log"

	"fmt"
	"github.com/bennyz/example-finder/util"

	"github.com/bennyz/example-finder/persistence"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS repo_data(repo_id INTEGER PRIMARY KEY, data JSON1)`
	insertKey   = `INSERT INTO repo_data(repo_id, data) values(?, ?)`
	getKeys     = `SELECT data FROM repo_data WHERE repo_id IN (%s)`
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
	var result []persistence.JSONValue

	var keysRaw []int
	for key := range keys {
		keysRaw = append(keysRaw, key)
	}

	query := fmt.Sprintf(getKeys, util.SliceToString(keys))
	rows, err := s.Query(query)
	if err != nil {
		log.Fatal("Failed querying rows. Query: %v, Keys: %v, error: %v", query, keys, err);
	}
	
	for rows.Next() {
		var value persistence.JSONValue
		err := rows.Scan(&value)
		if (err != nil) {
			log.Fatal("Failed scanning row. Rows: %v, err %v", rows, err);
		}
		result = append(result, value)
	}

	return result, err
}
