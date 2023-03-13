package memory

import (
	"sync"

	c "github.com/andrsj/feedback-service/internal/infrastructure/cache"
	"github.com/andrsj/feedback-service/pkg/logger"
)

type cache struct {
	mu     sync.RWMutex
	items  map[string][]byte
	logger logger.Logger
}

var _ c.Cache = (*cache)(nil)

func New(logger logger.Logger) *cache {
	return &cache{
		mu:     sync.RWMutex{},
		items:  make(map[string][]byte),
		logger: logger.Named("cache"),
	}
}

// Set adds a new item to the cache.
func (c *cache) Set(key string, value []byte) error {
	c.logger.Info("Setting values", logger.M{
		"key":   key,
		"value": string(value),
	})

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value

	return nil
}

// Get retrieves an item from the cache.
func (c *cache) Get(key string) ([]byte, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, keyExists := c.items[key]
	c.logger.Info("Getting values", logger.M{
		"key":   key,
		"value": string(value),
		"exist": keyExists,
	})

	return value, keyExists, nil
}
