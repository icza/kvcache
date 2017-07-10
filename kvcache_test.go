package kvcache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/icza/mighty"
)

const baseFolder = "testdata"

func TestPersisting(t *testing.T) {
	eq, expDeq := mighty.Eq(t), mighty.ExpDeq(t)

	folder := filepath.Join(baseFolder, t.Name())
	c, err := New(folder, "v1.0")
	eq(nil, err)
	eq(nil, c.Put("a", []byte("Aa")))
	eq(nil, c.Put("b", []byte("BB")))
	eq(nil, c.Close())

	c, err = New(folder, "v1.0")
	eq(nil, err)

	expDeq([]byte("Aa"))(c.Get("a"))
	expDeq([]byte("BB"))(c.Get("b"))

	os.RemoveAll(folder)
}

func TestLongVersion(t *testing.T) {
	eq := mighty.Eq(t)

	longVer := make([]byte, KeySizeLimit+1)
	_, err := New(filepath.Join(baseFolder, t.Name()), string(longVer))
	eq(ErrKeySize, err)
}
