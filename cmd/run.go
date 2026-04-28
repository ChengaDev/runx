package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ChengaDev/runx/internal/store"
)

var runCmd = &cobra.Command{
	Use:   "run <name> [-- <extra args>]",
	Short: "Run a saved project",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRun,
}

func runRun(cmd *cobra.Command, args []string) error {
	name := args[0]

	var extra []string
	dashIdx := cmd.ArgsLenAtDash()
	if dashIdx >= 0 {
		extra = args[dashIdx:]
	}

	s, err := store.New()
	if err != nil {
		return err
	}

	p, ok := s.Get(name)
	if !ok {
		return fmt.Errorf("project %q not found — run `runx list` to see available projects", name)
	}

	if _, err := os.Stat(p.Path); err != nil {
		return fmt.Errorf("project directory no longer exists: %s\nUpdate it with: runx edit %s --path <new-path>", p.Path, name)
	}

	cmdStr := p.Cmd
	if len(extra) > 0 {
		quoted := make([]string, len(extra))
		for i, a := range extra {
			quoted[i] = shellQuote(a)
		}
		cmdStr += " " + strings.Join(quoted, " ")
	}

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.CommandContext(cmd.Context(), "cmd", "/c", cmdStr)
	} else {
		c = exec.CommandContext(cmd.Context(), "sh", "-c", cmdStr)
	}

	c.Dir = p.Path
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return fmt.Errorf("command exited with error: %w", err)
	}

	return nil
}

// shellQuote wraps s in single quotes, escaping any single quotes within it,
// so the value is treated as a literal argument by the shell.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
