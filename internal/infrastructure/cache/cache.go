package cache

// Cache interface is used in
// controller and cache-middleware.
type Cache interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, value []byte) error
}
