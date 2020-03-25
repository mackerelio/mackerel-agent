package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMakeFallbackEnv_386(t *testing.T) {
	if runtime.GOARCH != "386" {
		t.Skip()
	}
	wdir := filepath.Join(os.Getenv("PROGRAMFILES"), "Mackerel", "mackerel-agent")
	dir := FallbackConfigDir(wdir)
	if dir != "" {
		t.Errorf("FallbackConfigDir(%q) = %q; want %q", wdir, dir, "")
	}
}

func TestMakeFallbackEnv_amd64(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip()
	}
	wdir := filepath.Join(os.Getenv("PROGRAMFILES"), "Mackerel", "mackerel-agent")
	dir := FallbackConfigDir(wdir)
	want := filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Mackerel", "mackerel-agent")
	if dir != want {
		t.Errorf("makeFallbackEnv(%q) = %q; want %q", wdir, dir, want)
	}
}
