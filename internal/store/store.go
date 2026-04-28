package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/viper"
)

type Project struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Cmd  string `json:"cmd"`
}

type Store struct {
	filePath string
	data     map[string]Project
}

func New() (*Store, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cannot create config directory %q: %w", dir, err)
	}

	s := &Store{
		filePath: filepath.Join(dir, "projects.json"),
		data:     make(map[string]Project),
	}

	return s, s.load()
}

func configDir() (string, error) {
	viper.SetEnvPrefix("RUNX")
	viper.AutomaticEnv()

	if override := viper.GetString("CONFIG_DIR"); override != "" {
		return override, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	return filepath.Join(home, ".runx"), nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot read projects file: %w", err)
	}

	if err := json.Unmarshal(data, &s.data); err != nil {
		return fmt.Errorf("projects file is corrupted: %w", err)
	}

	return nil
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot serialize projects: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0o644); err != nil {
		return fmt.Errorf("cannot write projects file: %w", err)
	}

	return nil
}

func (s *Store) Add(p Project) error {
	s.data[p.Name] = p
	return s.save()
}

func (s *Store) Get(name string) (Project, bool) {
	p, ok := s.data[name]
	return p, ok
}

func (s *Store) Remove(name string) error {
	if _, ok := s.data[name]; !ok {
		return fmt.Errorf("project %q not found — run `runx list` to see available projects", name)
	}
	delete(s.data, name)
	return s.save()
}

func (s *Store) List() []Project {
	projects := make([]Project, 0, len(s.data))
	for _, p := range s.data {
		projects = append(projects, p)
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
	return projects
}
