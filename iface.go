// This file contains the interface of the kvcache.

package kvcache

// Cache describes the key-value cache.
type Cache interface {
	// Get returns value associated with the given key.
	// nil is returned if the key is not found.
	Get(key string) ([]byte, error)

	// Put puts a new key-value pair into the cache.
	Put(key string, value []byte) error

	// Clear removes all key-value pairs from the cache.
	Clear() error

	// Close closes the cache, releases any associated resources.
	// Close is idempotent (may be called multiple times).
	Close() error
}
