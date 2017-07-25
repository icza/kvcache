// This file contains the interface of the kvcache.

package kvcache

// Cache describes the operations of the key-value cache.
type Cache interface {
	// Get returns the value associated with the given key.
	// nil slice and nil error is returned if the key is not found.
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

	// Len returns the number of key-value pairs in the Cache.
	Len() int

	// Stat returns some basic statistics about the cache,
	// including the occupied storage size of the cache in bytes.
	Stat() (*Stat, error)

	// Close closes the cache, releases any associated resources.
	Close() error
}

// Stat wraps basic statistics about the cache.
type Stat struct {
	// Len is the number of key-value pairs in the Cache (same as Cache.Len()).
	Len int

	// StorageSize is the total storage size of the cache in bytes,
	// it is the sum of IndexSize and DataSize.
	StorageSize int64

	// IndexSize is the total size of the keys (plus the version plus
	// 8 bytes of metadata per key) in bytes.
	IndexSize int64

	// DataSize is the total size of the values in bytes.
	DataSize int64
}
