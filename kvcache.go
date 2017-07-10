// This file holds the implementation of the Cache interface.

package kvcache

// cache is the implementation of the Cache interface.
type cache struct {
}

// New creates a new Cache, using the given folder for persisting the data.
func New(folder string) (Cache, error) {
	// TODO
	return nil, nil
}

// Get implements Cache.Get().
func (c *cache) Get(key string) ([]byte, error) {
	// TODO
	return nil, nil
}

// Get implements Cache.Get().
func (c *cache) Put(key string, value []byte) error {
	// TODO
	return nil
}

// Get implements Cache.Get().
func (c *cache) Clear() error {
	// TODO
	return nil
}

// Get implements Cache.Get().
func (c *cache) Close() error {
	// TODO
	return nil
}
