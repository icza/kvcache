// This file holds the implementation of the Cache interface.

package kvcache

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	// indexName is the name of the index file
	indexName = "index"
	// dataName is the name of the data file
	dataName = "data"
)

// cache is the implementation of the Cache interface.
type cache struct {
	// mutex to protect concurrent access
	sync.RWMutex

	// version of the data in the cache
	version string

	// indexf is the index file
	indexf *os.File

	// dataf is the data file
	dataf *os.File
}

// New creates a new Cache, using the given folder for persisting the data.
// You can also provide a version of the data. If a persisted cache already exists,
// and its version is different from the given one, it will be cleared before
// return.
// The cache will also be cleared if already exists but is invalid.
func New(folder, version string) (result Cache, err error) {
	// Safety check:
	const versionLimit = 1 << 10 // 1KB
	if len(version) > versionLimit {
		return nil, fmt.Errorf("too long version, limit is %d", versionLimit)
	}

	// Make sure folder exists:
	if err = os.MkdirAll(folder, 0775); err != nil {
		return
	}

	c := &cache{
		version: version,
	}
	defer func() {
		if err != nil {
			c.Close() // Close successfully opened files
		}
	}()

	c.indexf, err = os.OpenFile(
		filepath.Join(folder, indexName),
		os.O_CREATE|os.O_RDWR,
		0755,
	)
	if err != nil {
		return
	}

	c.dataf, err = os.OpenFile(
		filepath.Join(folder, dataName),
		os.O_CREATE|os.O_RDWR,
		0755,
	)
	if err != nil {
		return
	}

	stat, err := c.indexf.Stat()
	if err != nil {
		return
	}
	doClear := stat.Size() == 0
	if !doClear {
		// Read and check existing version:
		var vlen int32
		if err = binary.Read(c.indexf, binary.LittleEndian, &vlen); err != nil {
			return
		}
		// Safety limit:
		if vlen > versionLimit {
			vlen = versionLimit
		}
		existingVer := make([]byte, vlen)
		if _, err = io.ReadFull(c.indexf, existingVer); err != nil {
			return
		}
		if version != string(existingVer) {
			if err = c.Clear(); err != nil {
				return
			}
			doClear = true
		}
	}

	if doClear {
		if err = c.Clear(); err != nil {
			return
		}
	}

	return c, nil
}

// Get implements Cache.Get().
func (c *cache) Get(key string) ([]byte, error) {
	// TODO
	return nil, nil
}

// Put implements Cache.Put().
func (c *cache) Put(key string, value []byte) error {
	// TODO
	return nil
}

// Clear implements Cache.Clear().
func (c *cache) Clear() error {
	c.Lock()
	defer c.Unlock()

	if err := c.indexf.Truncate(0); err != nil {
		return err
	}
	if _, err := c.indexf.Seek(0, 0); err != nil {
		return err
	}

	if err := c.dataf.Truncate(0); err != nil {
		return err
	}
	if _, err := c.dataf.Seek(0, 0); err != nil {
		return err
	}

	// Write version: length + version bytes
	if err := binary.Write(c.indexf, binary.LittleEndian, len(c.version)); err != nil {
		return err
	}
	if _, err := c.indexf.WriteString(c.version); err != nil {
		return err
	}

	return nil
}

// Close implements Cache.Close().
func (c *cache) Close() error {
	var err1, err2 error
	if c.indexf != nil {
		err1 = c.indexf.Close()
	}
	if c.dataf != nil {
		err2 = c.dataf.Close()
	}

	if err1 != nil {
		return err1
	}
	return err2
}
