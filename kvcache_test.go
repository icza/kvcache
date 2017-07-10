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
	eq(nil, c.Put("b", []byte("Bb")))
	eq(nil, c.Close())

	c, err = New(folder, "v1.0")
	eq(nil, err)

	expDeq([]byte("Aa"))(c.Get("a"))
	expDeq([]byte("Bb"))(c.Get("b"))

	eq(nil, c.Close())
	os.RemoveAll(folder)
}

func TestPut(t *testing.T) {
	eq := mighty.Eq(t)

	folder := filepath.Join(baseFolder, t.Name())
	c, err := New(folder, "v1.0")
	eq(nil, err)

	eq(nil, c.Put("a", []byte("A")))
	eq(ErrKeyExists, c.Put("a", []byte("A")))

	longKey := make([]byte, KeySizeLimit+1)
	eq(ErrKeySize, c.Put(string(longKey), []byte("A")))

	eq(nil, c.Close())
	os.RemoveAll(folder)
}

func TestLongVersion(t *testing.T) {
	eq := mighty.Eq(t)

	folder := filepath.Join(baseFolder, t.Name())

	longVer := make([]byte, KeySizeLimit+1)
	_, err := New(folder, string(longVer))
	eq(ErrKeySize, err)
}

func TestClear(t *testing.T) {
	eq, expDeq := mighty.Eq(t), mighty.ExpDeq(t)

	folder := filepath.Join(baseFolder, t.Name())

	c, err := New(folder, "v1.0")
	eq(nil, err)
	eq(nil, c.Put("a", []byte("Aa")))
	eq(nil, c.Put("b", []byte("Bb")))

	eq(nil, c.Clear())

	expDeq([]byte(nil))(c.Get("a"))
	expDeq([]byte(nil))(c.Get("b"))

	eq(nil, c.Close())
	os.RemoveAll(folder)
}

func TestVersionMismatch(t *testing.T) {
	eq, expDeq := mighty.Eq(t), mighty.ExpDeq(t)

	folder := filepath.Join(baseFolder, t.Name())

	c, err := New(folder, "v1.0")
	eq(nil, err)
	eq(nil, c.Put("a", []byte("Aa")))
	eq(nil, c.Close())

	c, err = New(folder, "v1.1")
	eq(nil, err)

	expDeq([]byte(nil))(c.Get("a"))

	eq(nil, c.Close())
	os.RemoveAll(folder)
}
