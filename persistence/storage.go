package persistence

// Storage backend to be used by the persistence layer
type Storage interface {
	Save(key int64, value []byte) (int64, error)
	Get(key int64) ([]byte, error)
}
