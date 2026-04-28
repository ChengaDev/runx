package cmd

import (
	"fmt"
	"os"

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
	// Launch wizard when no name argument is provided
	if len(args) == 0 {
		return runAddWizard()
	}

	name := args[0]
	path, _ := cmd.Flags().GetString("path")
	command, _ := cmd.Flags().GetString("cmd")

	if path == "" {
		return fmt.Errorf("required flag \"path\" not set")
	}

	if err := validateDir(path); err != nil {
		return err
	}

	if command == "" {
		detected, ok := detect.Command(path)
		if !ok {
			return fmt.Errorf(
				"could not auto-detect a run command for %q\n"+
					"Please provide one explicitly with --cmd",
				path,
			)
		}
		command = detected
		fmt.Printf("Auto-detected command: %s\n", command)
	}

	return saveProject(name, path, command)
}

func runAddWizard() error {
	var name, path, command string

	// Step 1 & 2: name and path
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

			huh.NewInput().
				Title("Project path").
				Placeholder("/path/to/project").
				Validate(func(s string) error {
					return validateDir(s)
				}).
				Value(&path),
		),
	).Run(); err != nil {
		return err
	}

	// Step 3: command — pre-fill with auto-detected value if available
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
