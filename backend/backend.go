package backends

// Backend specifies which backend should handle the request
// Currently only rest is supported
type Backend interface {
	Search(string, string) map[int64]*Result
}
