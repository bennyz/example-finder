package backends

// Result hold the result to amount of stars mapping
type Result struct {
	RepoID   int64
	RepoName string
	RepoURL  string
	Stars    int
}
