package cache

type Cache interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, value []byte) error
}
