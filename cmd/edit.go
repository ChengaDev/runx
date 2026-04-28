package cmd

import (
	"fmt"

	"github.com/ChengaDev/runx/internal/detect"
	"github.com/ChengaDev/runx/internal/store"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit a saved project's path or command",
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	editCmd.Flags().StringP("path", "p", "", "New directory path")
	editCmd.Flags().StringP("cmd", "c", "", "New command")
}

func runEdit(cmd *cobra.Command, args []string) error {
	name := args[0]
	newPath, _ := cmd.Flags().GetString("path")
	newCmd, _ := cmd.Flags().GetString("cmd")

	if newPath == "" && newCmd == "" {
		return fmt.Errorf("provide at least one of --path or --cmd")
	}

	s, err := store.New()
	if err != nil {
		return err
	}

	p, ok := s.Get(name)
	if !ok {
		return fmt.Errorf("project %q not found — run `runx list` to see available projects", name)
	}

	if newPath != "" {
		if err := validateDir(newPath); err != nil {
			return err
		}
		p.Path = newPath

		// Re-detect command when path changes and no explicit cmd given
		if newCmd == "" {
			if detected, ok := detect.Command(p.Path); ok {
				p.Cmd = detected
				fmt.Printf("Auto-detected command: %s\n", detected)
			}
		}
	}

	if newCmd != "" {
		p.Cmd = newCmd
	}

	if err := s.Add(p); err != nil {
		return err
	}

	fmt.Printf("Project %q updated.\n  path: %s\n  cmd:  %s\n", name, p.Path, p.Cmd)
	return nil
}
