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

var hereCmd = &cobra.Command{
	Use:   "here",
	Short: "Save the current directory as a project",
	Args:  cobra.NoArgs,
	RunE:  runHere,
}

func runHere(_ *cobra.Command, _ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	defaultName := filepath.Base(cwd)

	detected, ok := detect.Command(cwd)
	if !ok {
		return fmt.Errorf(
			"could not auto-detect a run command for %q\n"+
				"Use `runx add` to set one explicitly with --cmd",
			cwd,
		)
	}

	var name, command string
	name = defaultName
	command = detected

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("Defaults to current directory name.").
				Value(&name).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Run command").
				Description("Auto-detected from project files. Edit or leave as-is.").
				Value(&command).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("command cannot be empty")
					}
					return nil
				}),
		),
	).Run(); err != nil {
		return err
	}

	s, err := store.New()
	if err != nil {
		return err
	}

	if err := s.Add(store.Project{Name: name, Path: cwd, Cmd: command}); err != nil {
		return err
	}

	fmt.Printf("Project %q saved.\n  path: %s\n  cmd:  %s\n", name, cwd, command)
	return nil
}
