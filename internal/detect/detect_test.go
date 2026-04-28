package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDirWith(t *testing.T, files ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte{}, 0o644); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	return dir
}

func TestCommand_PackageJSON(t *testing.T) {
	dir := tempDirWith(t, "package.json")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "npm run dev")
}

func TestCommand_ManagePy(t *testing.T) {
	dir := tempDirWith(t, "manage.py")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "python manage.py runserver")
}

func TestCommand_CargoToml(t *testing.T) {
	dir := tempDirWith(t, "Cargo.toml")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "cargo run")
}

func TestCommand_GoMod(t *testing.T) {
	dir := tempDirWith(t, "go.mod")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "go run .")
}

func TestCommand_PyprojectToml(t *testing.T) {
	dir := tempDirWith(t, "pyproject.toml")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "python main.py")
}

func TestCommand_MainPy(t *testing.T) {
	dir := tempDirWith(t, "main.py")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "python main.py")
}

func TestCommand_NoMatch(t *testing.T) {
	dir := tempDirWith(t, "README.md", "somefile.txt")
	_, ok := Command(dir)
	if ok {
		t.Error("Command() returned true for unrecognized project, want false")
	}
}

func TestCommand_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	_, ok := Command(dir)
	if ok {
		t.Error("Command() returned true for empty directory, want false")
	}
}

// package.json takes priority over go.mod when both are present.
func TestCommand_Priority_PackageOverGo(t *testing.T) {
	dir := tempDirWith(t, "package.json", "go.mod")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "npm run dev")
}

// pyproject.toml takes priority over main.py when both are present.
func TestCommand_Priority_PyprojectOverMainPy(t *testing.T) {
	dir := tempDirWith(t, "pyproject.toml", "main.py")
	cmd, ok := Command(dir)
	assertDetected(t, cmd, ok, "python main.py")
}

func assertDetected(t *testing.T, cmd string, ok bool, want string) {
	t.Helper()
	if !ok {
		t.Fatalf("Command() returned false, want true")
	}
	if cmd != want {
		t.Errorf("Command() = %q, want %q", cmd, want)
	}
}
