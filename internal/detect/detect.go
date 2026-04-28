package detect

import (
	"os"
	"path/filepath"
)

type indicator struct {
	file string
	cmd  string
}

var indicators = []indicator{
	{"package.json", "npm run dev"},
	{"manage.py", "python manage.py runserver"},
	{"Cargo.toml", "cargo run"},
	{"go.mod", "go run ."},
	{"pyproject.toml", "python main.py"},
	{"main.py", "python main.py"},
}

// Command returns the best-guess run command for the given directory.
func Command(dir string) (string, bool) {
	for _, ind := range indicators {
		if fileExists(filepath.Join(dir, ind.file)) {
			return ind.cmd, true
		}
	}
	return "", false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
