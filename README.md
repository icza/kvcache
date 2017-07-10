# kvcache

[![Build Status](https://travis-ci.org/icza/kvcache.svg?branch=master)](https://travis-ci.org/icza/kvcache)
[![GoDoc](https://godoc.org/github.com/icza/kvcache?status.svg)](https://godoc.org/github.com/icza/kvcache)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/kvcache)](https://goreportcard.com/report/github.com/icza/kvcache)

Simple, optimized, embedded, persistent (file-based) key-value cache.

## Features

- Very simple interface. Basically just a Get and a Put operation. Does not support
removing elements, but it supports removing all elements (reset) with the Clear method.

- Optimized. Keys are kept in memory for fast lookups.

- Embedded. You init / create your cache from within your app. No external services
need to run.

- Persistent. Data (key-value pairs) are written to files in a folder given at
creation time.

- Supports concurrent access.
