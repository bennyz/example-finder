package backends

// Result hold the result to amount of stars mapping
type Result struct {
	RepoID    int64
	RepoName  string
	RepoURL   string
	Stars     int
	FilePaths []string
}

// ByStars a type to sort results by stars
type ByStars []*Result

func (s ByStars) Len() int {
	return len(s)
}

func (s ByStars) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByStars) Less(i, j int) bool {
	return s[i].Stars > s[j].Stars
}
