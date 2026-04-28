package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/ChengaDev/runx/internal/detect"
	"github.com/ChengaDev/runx/internal/store"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Save a project",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runAdd,
}

func init() {
	addCmd.Flags().StringP("path", "p", "", "Directory path of the project")
	addCmd.Flags().StringP("cmd", "c", "", "Command to run (auto-detected if omitted)")
}

func runAdd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return runAddWizard()
	}

	name := args[0]
	path, _ := cmd.Flags().GetString("path")
	command, _ := cmd.Flags().GetString("cmd")

	if path == "" {
		return fmt.Errorf("required flag \"path\" not set")
	}

	absPath, err := resolveDir(path)
	if err != nil {
		return err
	}

	if command == "" {
		detected, ok := detect.Command(absPath)
		if !ok {
			return fmt.Errorf(
				"could not auto-detect a run command for %q\n"+
					"Please provide one explicitly with --cmd",
				absPath,
			)
		}
		command = detected
		fmt.Printf("Auto-detected command: %s\n", command)
	}

	return saveProject(name, absPath, command)
}

func runAddWizard() error {
	cwd, _ := os.Getwd()
	var name, path, command string

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Placeholder("my-app").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}).
				Value(&name),

			huh.NewFilePicker().
				Title("Project path").
				Description("Navigate with arrow keys, press Enter to select a directory.").
				CurrentDirectory(cwd).
				DirAllowed(true).
				FileAllowed(false).
				Value(&path),
		),
	).Run(); err != nil {
		return err
	}

	detected, _ := detect.Command(path)
	command = detected

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Run command").
				Description("Auto-detected from project files. Edit or leave as-is.").
				Placeholder("e.g. go run . / npm run dev").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("command cannot be empty")
					}
					return nil
				}).
				Value(&command),
		),
	).Run(); err != nil {
		return err
	}

	return saveProject(name, path, command)
}

func saveProject(name, path, command string) error {
	s, err := store.New()
	if err != nil {
		return err
	}

	if err := s.Add(store.Project{Name: name, Path: path, Cmd: command}); err != nil {
		return err
	}

	fmt.Printf("Project %q saved.\n  path: %s\n  cmd:  %s\n", name, path, command)
	return nil
}

// resolveDir resolves a path (including ".") to an absolute path and validates it's a directory.
func resolveDir(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", path, err)
	}
	return abs, validateDir(abs)
}

func validateDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path %q does not exist", path)
	}
	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", path)
	}
	return nil
}
