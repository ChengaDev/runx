package cmd

import (
	"fmt"

	"github.com/ChengaDev/runx/internal/store"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Delete a saved project",
	Args:  cobra.ExactArgs(1),
	RunE:  runRemove,
}

func runRemove(_ *cobra.Command, args []string) error {
	name := args[0]

	s, err := store.New()
	if err != nil {
		return err
	}

	if err := s.Remove(name); err != nil {
		return err
	}

	fmt.Printf("Project %q removed.\n", name)
	return nil
}
