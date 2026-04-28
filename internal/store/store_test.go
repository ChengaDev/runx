package store

import (
	"os"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("RUNX_CONFIG_DIR", dir)
	s, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	return s
}

func TestAdd_and_Get(t *testing.T) {
	s := newTestStore(t)

	p := Project{Name: "myapp", Path: "/tmp/myapp", Cmd: "go run ."}
	if err := s.Add(p); err != nil {
		t.Fatalf("Add() error: %v", err)
	}

	got, ok := s.Get("myapp")
	if !ok {
		t.Fatal("Get() returned false, want true")
	}
	if got != p {
		t.Errorf("Get() = %+v, want %+v", got, p)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := newTestStore(t)

	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("Get() returned true for missing project, want false")
	}
}

func TestRemove(t *testing.T) {
	s := newTestStore(t)

	p := Project{Name: "myapp", Path: "/tmp/myapp", Cmd: "go run ."}
	_ = s.Add(p)

	if err := s.Remove("myapp"); err != nil {
		t.Fatalf("Remove() error: %v", err)
	}

	_, ok := s.Get("myapp")
	if ok {
		t.Error("Get() returned true after Remove, want false")
	}
}

func TestRemove_NotFound(t *testing.T) {
	s := newTestStore(t)

	err := s.Remove("ghost")
	if err == nil {
		t.Error("Remove() expected error for missing project, got nil")
	}
}

func TestList_Empty(t *testing.T) {
	s := newTestStore(t)

	if got := s.List(); len(got) != 0 {
		t.Errorf("List() len = %d, want 0", len(got))
	}
}

func TestList_Sorted(t *testing.T) {
	s := newTestStore(t)

	_ = s.Add(Project{Name: "zebra", Path: "/z", Cmd: "z"})
	_ = s.Add(Project{Name: "alpha", Path: "/a", Cmd: "a"})
	_ = s.Add(Project{Name: "mango", Path: "/m", Cmd: "m"})

	list := s.List()
	if len(list) != 3 {
		t.Fatalf("List() len = %d, want 3", len(list))
	}

	names := []string{list[0].Name, list[1].Name, list[2].Name}
	want := []string{"alpha", "mango", "zebra"}
	for i, name := range names {
		if name != want[i] {
			t.Errorf("List()[%d].Name = %q, want %q", i, name, want[i])
		}
	}
}

func TestAdd_Overwrite(t *testing.T) {
	s := newTestStore(t)

	_ = s.Add(Project{Name: "app", Path: "/old", Cmd: "old cmd"})
	_ = s.Add(Project{Name: "app", Path: "/new", Cmd: "new cmd"})

	got, _ := s.Get("app")
	if got.Path != "/new" || got.Cmd != "new cmd" {
		t.Errorf("Add() did not overwrite: got %+v", got)
	}
}

func TestPersistence(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RUNX_CONFIG_DIR", dir)

	s1, _ := New()
	_ = s1.Add(Project{Name: "persisted", Path: "/p", Cmd: "run"})

	// Load a fresh store from the same directory
	s2, err := New()
	if err != nil {
		t.Fatalf("second New() error: %v", err)
	}

	got, ok := s2.Get("persisted")
	if !ok {
		t.Fatal("Get() returned false after reload, want true")
	}
	if got.Path != "/p" {
		t.Errorf("Path = %q, want %q", got.Path, "/p")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RUNX_CONFIG_DIR", dir)

	// Remove the file if it was created
	_ = os.Remove(dir + "/projects.json")

	s, err := New()
	if err != nil {
		t.Fatalf("New() should not error when file is absent: %v", err)
	}
	if len(s.List()) != 0 {
		t.Error("expected empty store when file is absent")
	}
}
