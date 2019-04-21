package persistence

// JSONValue holds a json representation of a struct
type JSONValue []byte

// Storage backend to be used by the persistence layer
type Storage interface {
	Save(key int64, value JSONValue) (int64, error)

	// Get JSON values for keys
	Get(keys []int64) ([]JSONValue, error)
}
