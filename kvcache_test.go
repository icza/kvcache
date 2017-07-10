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
	eq(0, c.Len())
	eq(nil, c.Put("a", []byte("Aa")))
	eq(nil, c.Put("b", []byte("Bb")))
	eq(2, c.Len())
	eq(nil, c.Close())

	c, err = New(folder, "v1.0")
	eq(nil, err)
	eq(2, c.Len())

	expDeq([]byte("Aa"))(c.Get("a"))
	expDeq([]byte("Bb"))(c.Get("b"))

	eq(nil, c.Close())
	os.RemoveAll(folder)
}

func TestPutGet(t *testing.T) {
	eq, expDeq := mighty.Eq(t), mighty.ExpDeq(t)

	folder := filepath.Join(baseFolder, t.Name())
	c, err := New(folder, "v1.0")
	eq(nil, err)

	eq(0, c.Len())
	expDeq([]byte(nil))(c.Get("a"))
	eq(nil, c.Put("a", []byte("A")))
	eq(1, c.Len())
	expDeq([]byte("A"))(c.Get("a"))

	eq(ErrKeyExists, c.Put("a", []byte("A")))
	eq(1, c.Len())
	expDeq([]byte("A"))(c.Get("a"))

	longKey := make([]byte, KeySizeLimit+1)
	eq(ErrKeySize, c.Put(string(longKey), []byte("A")))
	eq(1, c.Len())

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
	eq(0, c.Len())
	eq(nil, c.Put("a", []byte("Aa")))
	eq(nil, c.Put("b", []byte("Bb")))
	eq(2, c.Len())

	eq(nil, c.Clear())

	eq(0, c.Len())
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
	eq(0, c.Len())
	eq(nil, c.Put("a", []byte("Aa")))
	eq(1, c.Len())
	eq(nil, c.Close())

	c, err = New(folder, "v1.1")
	eq(nil, err)

	eq(0, c.Len())
	expDeq([]byte(nil))(c.Get("a"))

	eq(nil, c.Close())
	os.RemoveAll(folder)
}
