package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDir_Dot(t *testing.T) {
	cwd, _ := os.Getwd()
	want, _ := filepath.EvalSymlinks(cwd)
	got, err := resolveDir(".")
	if err != nil {
		t.Fatalf("resolveDir(\".\") error: %v", err)
	}
	gotReal, _ := filepath.EvalSymlinks(got)
	if gotReal != want {
		t.Errorf("resolveDir(\".\") = %q, want %q", gotReal, want)
	}
}

func TestResolveDir_Absolute(t *testing.T) {
	dir := t.TempDir()
	want, _ := filepath.EvalSymlinks(dir)
	got, err := resolveDir(dir)
	if err != nil {
		t.Fatalf("resolveDir(%q) error: %v", dir, err)
	}
	gotReal, _ := filepath.EvalSymlinks(got)
	if gotReal != want {
		t.Errorf("resolveDir(%q) = %q, want %q", dir, gotReal, want)
	}
}

func TestResolveDir_Relative(t *testing.T) {
	dir := t.TempDir()
	want, _ := filepath.EvalSymlinks(dir)
	parent := filepath.Dir(dir)
	base := filepath.Base(dir)

	cwd, _ := os.Getwd()
	_ = os.Chdir(parent)
	defer os.Chdir(cwd)

	got, err := resolveDir(base)
	if err != nil {
		t.Fatalf("resolveDir(%q) error: %v", base, err)
	}
	gotReal, _ := filepath.EvalSymlinks(got)
	if gotReal != want {
		t.Errorf("resolveDir(%q) = %q, want %q", base, gotReal, want)
	}
}

func TestResolveDir_NotExist(t *testing.T) {
	_, err := resolveDir("/nonexistent/path/xyz")
	if err == nil {
		t.Error("resolveDir() expected error for non-existent path, got nil")
	}
}

func TestResolveDir_File(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	_ = os.WriteFile(f, []byte{}, 0o644)

	_, err := resolveDir(f)
	if err == nil {
		t.Error("resolveDir() expected error for file path, got nil")
	}
}

func TestValidateDir_Valid(t *testing.T) {
	dir := t.TempDir()
	if err := validateDir(dir); err != nil {
		t.Errorf("validateDir(%q) unexpected error: %v", dir, err)
	}
}

func TestValidateDir_NotExist(t *testing.T) {
	if err := validateDir("/does/not/exist"); err == nil {
		t.Error("validateDir() expected error for missing path, got nil")
	}
}

func TestValidateDir_IsFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "notadir.txt")
	_ = os.WriteFile(f, []byte{}, 0o644)

	if err := validateDir(f); err == nil {
		t.Error("validateDir() expected error for file, got nil")
	}
}
