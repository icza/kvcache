// This file holds the implementation of the Cache interface.

package kvcache

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	// KeySizeLimit is the max allowed length for (individual) keys
	// and the version string
	KeySizeLimit = 1<<16 - 1 // 64 KB

	// DataSizeLimit is the max allowed total data size
	DataSizeLimit = 1<<32 - 1 // 4 GB
)

const (
	// indexName is the name of the index file
	indexName = "index"

	// dataName is the name of the data file
	dataName = "data"
)

var (
	// ErrKeyExists is returned when attempting to put an existing key into the cache
	ErrKeyExists = errors.New("key already in cache")

	// ErrKeySize is returned when attempting to put a too long key into the cache
	ErrKeySize = errors.New("key too long")

	// ErrDataSize is returned when attempting to put a value into the cache
	// which would raise total data size over the limit (DataSizeLimit)
	ErrDataSize = errors.New("total data too big")
)

// valueInfo describes a value in the index map.
type valueInfo struct {
	// Pos is the byte position of the data
	Pos uint32

	// Size is the byte-size of the data
	Size uint32
}

// cache is the implementation of the Cache interface.
type cache struct {
	// mutex to protect concurrent access
	sync.Mutex

	// folder of the cache where data is persisted
	folder string

	// version of the data in the cache
	version string

	// indexf is the index file
	indexf *os.File

	// dataf is the data file
	dataf *os.File

	// indexMap is the in-memory index
	indexMap map[string]valueInfo
}

// New creates a new Cache, using the given folder for persisting the data.
// You also have to provide the version of the data. If a persisted cache already
// exists, and its version is different from this, it will be cleared
// before return.
//
// ErrKeySize is returned if version is too long (>KeySizeLimit).
func New(folder, version string) (cch Cache, err error) {
	if len(version) > KeySizeLimit {
		return nil, ErrKeySize
	}

	// Make sure folder exists:
	if err = os.MkdirAll(folder, 0775); err != nil {
		return
	}

	c := &cache{
		Mutex:    sync.Mutex{},
		folder:   folder,
		version:  version,
		indexMap: map[string]valueInfo{},
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
		var vlen uint16
		if err = binary.Read(c.indexf, binary.LittleEndian, &vlen); err != nil {
			return
		}
		existingVer := make([]byte, vlen)
		if _, err = io.ReadFull(c.indexf, existingVer); err != nil {
			return
		}
		doClear = version != string(existingVer)
	}

	if doClear {
		if err = c.Clear(); err != nil {
			return
		}
	}

	// Read index into memory:
	for {
		var keyLen uint16
		if err = binary.Read(c.indexf, binary.LittleEndian, &keyLen); err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		key := make([]byte, keyLen)
		if _, err = io.ReadFull(c.indexf, key); err != nil {
			return
		}
		vi := valueInfo{}
		if err = binary.Read(c.indexf, binary.LittleEndian, &vi.Pos); err != nil {
			return
		}
		if err = binary.Read(c.indexf, binary.LittleEndian, &vi.Size); err != nil {
			return
		}
		c.indexMap[string(key)] = vi
	}

	return c, nil
}

// Get implements Cache.Get().
func (c *cache) Get(key string) ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	vi, ok := c.indexMap[key]
	if !ok {
		return nil, nil
	}

	value := make([]byte, vi.Size)
	if _, err := c.dataf.ReadAt(value, int64(vi.Pos)); err != nil {
		return nil, err
	}

	return value, nil
}

// Put implements Cache.Put().
func (c *cache) Put(key string, value []byte) error {
	if len(key) > KeySizeLimit {
		return ErrKeySize
	}

	c.Lock()
	defer c.Unlock()

	vi, ok := c.indexMap[key]
	if ok {
		return ErrKeyExists
	}

	// Write value into data file, at the end:
	pos, err := c.dataf.Seek(0, 2)
	if err != nil {
		return err
	}
	if pos+int64(len(value)) > DataSizeLimit {
		return ErrDataSize
	}

	vi.Pos = uint32(pos)
	vi.Size = uint32(len(value))

	if _, err = c.dataf.Write(value); err != nil {
		return err
	}

	// Write index entry
	// Index file position is always at the end
	if err = binary.Write(c.indexf, binary.LittleEndian, uint16(len(key))); err != nil {
		return err
	}
	if _, err = c.indexf.WriteString(key); err != nil {
		return err
	}
	if err = binary.Write(c.indexf, binary.LittleEndian, &vi.Pos); err != nil {
		return err
	}
	if err = binary.Write(c.indexf, binary.LittleEndian, &vi.Size); err != nil {
		return err
	}

	c.indexMap[key] = vi

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
	if err := binary.Write(c.indexf, binary.LittleEndian, uint16(len(c.version))); err != nil {
		return err
	}
	if _, err := c.indexf.WriteString(c.version); err != nil {
		return err
	}

	c.indexMap = map[string]valueInfo{}

	return nil
}

// Folder implements Cache.Folder()
func (c *cache) Folder() string {
	return c.folder
}

// Len implements Cache.Len().
func (c *cache) Len() int {
	c.Lock()
	defer c.Unlock()

	return len(c.indexMap)
}

// Stat implements Cache.Stat().
func (c *cache) Stat() (*Stat, error) {
	c.Lock()
	defer c.Unlock()

	s := &Stat{Len: len(c.indexMap)}

	var err error
	if s.IndexSize, err = c.indexf.Seek(0, 1); err != nil { // Stay at current offset (always the end)
		return nil, err
	}

	if s.DataSize, err = c.dataf.Seek(0, 2); err != nil { // Seek to the end
		return nil, err
	}

	s.StorageSize = s.IndexSize + s.DataSize

	return s, nil
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
