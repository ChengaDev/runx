package cmd

import (
	"fmt"

	"github.com/ChengaDev/runx/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all saved projects",
	RunE:  runList,
}

func runList(_ *cobra.Command, _ []string) error {
	s, err := store.New()
	if err != nil {
		return err
	}

	projects := s.List()
	if len(projects) == 0 {
		fmt.Println("No projects saved yet. Use `runx add` to add one.")
		return nil
	}

	fmt.Printf("%-20s  %-40s  %s\n", "NAME", "PATH", "CMD")
	fmt.Printf("%-20s  %-40s  %s\n", "----", "----", "---")
	for _, p := range projects {
		fmt.Printf("%-20s  %-40s  %s\n", p.Name, p.Path, p.Cmd)
	}

	return nil
}
