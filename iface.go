// This file contains the interface of the kvcache.

package kvcache

// Cache describes the key-value cache.
type Cache interface {
	// Get returns value associated with the given key.
	// nil is returned if the key is not found.
	Get(key string) ([]byte, error)

	// Put puts a new key-value pair into the cache.
	//
	// ErrKeyExists is returned if key is already in the case.
	// ErrKeySize is returned if key is too long (> KeySizeLimit).
	// ErrDataSize is returned if putting the value would increase total
	// data size over the limit (DataSizeLimit).
	Put(key string, value []byte) error

	// Clear removes all key-value pairs from the cache.
	Clear() error

	// Close closes the cache, releases any associated resources.
	Close() error
}
