package main

import (
	"os"
	"path/filepath"
)

// PathValue is a flag that is used to point to a file. PathValue picks up appropriate path.
//
//	1. A path passed by command line options if it passed
//	2. A fallback path that it is set with SetFallback if it exists
//
// Otherwise PathValue uses Path field value.
type PathValue struct {
	Path string

	updated bool
	file    string
}

// String implements flag.Value interface.
func (f *PathValue) String() string {
	return f.Path
}

// Set implements flag.Value interface.
func (f *PathValue) Set(s string) error {
	f.Path = s
	f.updated = true
	return nil
}

// SetFallback sets the fallback target.
func (f *PathValue) SetFallback(file string) {
	f.file = file
}

// SetDefault sets the default path.
func (f *PathValue) SetDefault(s string) {
	if f.updated {
		return
	}
	f.Path = s
}

// ResolveFile returns file only if f.Path is not exist and file is exist. Otherwise it returns f.Path.
func (f *PathValue) ResolveFile() string {
	if f.updated {
		return f.Path
	}
	if _, err := os.Stat(f.Path); err == nil || !os.IsNotExist(err) {
		return f.Path
	}
	if f.file == "" {
		return f.Path
	}
	if _, err := os.Stat(f.file); err != nil {
		return f.Path
	}
	return f.file
}

// ResolveDir returns dir only if f.Path/name is not exist and dir/name is exist. Otherwise it returns f.Path.
func (f *PathValue) ResolveDir() string {
	if f.updated {
		return f.Path
	}
	if f.file == "" {
		return f.Path
	}
	dir := filepath.Dir(f.file)
	name := filepath.Base(f.file)
	file := filepath.Join(f.Path, name)
	if _, err := os.Stat(file); err == nil || !os.IsNotExist(err) {
		return f.Path
	}
	if _, err := os.Stat(f.file); err != nil {
		return f.Path
	}
	return dir
}
