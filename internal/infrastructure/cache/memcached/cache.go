package memcached

import (
	"errors"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/andrsj/feedback-service/internal/infrastructure/cache"
	"github.com/andrsj/feedback-service/pkg/logger"
)

type Memcached struct {
	logger        logger.Logger
	client        *memcache.Client
	secondsToLive int32
}

// Check that actual implementation fits the interface.
var _ cache.Cache = (*Memcached)(nil)

func New(host string, secondsToLive int32, logger logger.Logger) *Memcached {
	return &Memcached{
		logger:        logger.Named("memcached"),
		client:        memcache.New(host),
		secondsToLive: secondsToLive,
	}
}

func (c *Memcached) Get(key string) ([]byte, bool, error) {
	item, err := c.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			c.logger.Info("Cache miss", logger.M{"key": key})

			return nil, false, nil
		}

		c.logger.Error("Error Memcached", logger.M{"err": err})

		return nil, false, fmt.Errorf("getting cache: %w", err)
	}

	c.logger.Info("Successfully get", logger.M{"value": string(item.Value)})

	return item.Value, true, nil
}

func (c *Memcached) Set(key string, value []byte) error {
	//nolint:exhaustivestruct,exhaustruct
	err := c.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: c.secondsToLive,
	})
	if err != nil {
		c.logger.Error("setting error", logger.M{"err": err})

		return fmt.Errorf("setting cache: %w", err)
	}

	c.logger.Info("Successfully set", logger.M{"key": key, "value": string(value)})

	return nil
}
