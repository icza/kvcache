/*
Package kvcache is a simple, optimized, embedded, persistent (file-based) key-value cache.

Features

- Very simple interface. Basically just a Get and a Put operation. Does not support
removing elements, but it supports removing all elements (reset) with the Clear method.

- Optimized. Keys are kept in memory for fast lookups.

- Embedded. You init / create your cache from within your app. No external services
need to run.

- Persistent. Data (key-value pairs) are written to files in a folder given at
creation time.

- Supports concurrent access.

Notes

Since element removal is not supported, the usability of this cache implementation
is limited, but in exchange it provides very compact storage: basically the
required storage size equals to the size of keys and the associated values
(plus a very tiny overhead).

*/
package kvcache
