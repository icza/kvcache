# kvcache

![Build Status](https://github.com/icza/kvcache/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/icza/kvcache.svg)](https://pkg.go.dev/github.com/icza/kvcache)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/kvcache)](https://goreportcard.com/report/github.com/icza/kvcache)

Simple, optimized, embedded, persistent (file-based) key-value cache.

## Intended use case

kvcache is intended for caching data

- that can be reproduced at any time
- but reproduction is time or resource consuming
- not feasible to store all in memory
- desirable to persist (remains between application restarts)

## Features

- Very simple interface. Basically just a `Get` and a `Put` operation. Does not support
changing or removing elements, but it supports removing all elements with the
`Clear` method.
- Optimized. Keys are kept in memory for fast lookups.
- Embedded. You init / create your cache from within your app. No external services
need to run.
- Persistent. Data (key-value pairs) are written to files in a folder given at
creation time.
- Simple data versioning (global to the cache). If version supplied at
create / init time differs, cache will be auto-cleared.
- Supports concurrent access.

## Notes

Since element removal and changing is not supported, the usability of this cache
implementation is limited, but in exchange it provides fast access and very
compact storage: basically the required storage size equals to the total size
of values and total size of keys (plus a tiny overhead per key).

## Implementation restrictions

- Length of the data version and (individual) keys must be less than 64 KB (`1<<16 - 1`),
exposed as `KeySizeLimit`.
- Total data size (total size of values) must not exceed 4 GB (`1<<32 - 1`),
exposed as `DataSizeLimit`.
