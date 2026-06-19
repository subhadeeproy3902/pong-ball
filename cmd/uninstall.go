package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/subhadeeproy3902/pong-ball/store"
)

// uninstallCmd builds the `pong-ball uninstall` command. It removes the running
// binary plus all of pong-ball's local data (scores, config, cached sound
// files) after a Y/N confirmation.
func uninstallCmd() *cobra.Command {
	var assumeYes, keepData bool

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the pong-ball binary and its local data",
		Long: "Remove the pong-ball binary together with its saved scores, config,\n" +
			"and cached sound files. Prompts for confirmation unless --yes is given.\n\n" +
			"If pong-ball was installed through a package manager (Homebrew, Scoop,\n" +
			"apt, …), prefer that manager's own uninstall so its receipt is cleared too.",
		Example: "  pong-ball uninstall\n  pong-ball uninstall --yes\n  pong-ball uninstall --keep-data",
		RunE: func(c *cobra.Command, args []string) error {
			return runUninstall(assumeYes, keepData)
		},
	}

	cmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "Skip the confirmation prompt")
	cmd.Flags().BoolVar(&keepData, "keep-data", false, "Remove the binary but keep saved scores and config")
	return cmd
}

func runUninstall(assumeYes, keepData bool) error {
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#cc785c")).Bold(true)
	warn := lipgloss.NewStyle().Foreground(lipgloss.Color("#d8a657"))

	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not locate the pong-ball binary: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(self); err == nil {
		self = resolved
	}

	dataDir := store.DataDir()
	sfxDir := filepath.Join(os.TempDir(), "pong-ball-sfx")

	fmt.Println(accent.Render("pong-ball uninstall"))
	fmt.Println()
	fmt.Println(warn.Render("This will permanently remove:"))
	fmt.Printf("  - the binary           %s\n", self)
	if !keepData {
		fmt.Printf("  - saved scores/config  %s\n", dataDir)
		fmt.Printf("  - cached sound files    %s\n", sfxDir)
	}
	fmt.Println()

	if !assumeYes {
		fmt.Print("Type YES to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "YES" {
			fmt.Println("Aborted - nothing was removed.")
			return nil
		}
	}

	if !keepData {
		if err := os.RemoveAll(dataDir); err != nil {
			fmt.Printf("  ! could not remove %s: %v\n", dataDir, err)
		} else {
			fmt.Printf("  removed %s\n", dataDir)
		}
		if err := os.RemoveAll(sfxDir); err != nil {
			fmt.Printf("  ! could not remove %s: %v\n", sfxDir, err)
		} else {
			fmt.Printf("  removed %s\n", sfxDir)
		}
	}

	if err := removeSelf(self); err != nil {
		fmt.Printf("  ! could not remove %s: %v\n", self, err)
		fmt.Println("    Delete it manually, or use your package manager's uninstall.")
		return nil
	}

	fmt.Println()
	fmt.Println(accent.Render("pong-ball uninstalled. Thanks for playing!"))
	fmt.Println("  (Installed via a package manager? Run its uninstall too, e.g.")
	fmt.Println("   `brew uninstall pong-ball` or `scoop uninstall pong-ball`.)")
	return nil
}

// removeSelf deletes the running binary. On Unix the file can be unlinked while
// the process still runs; on Windows a running .exe is locked, so we hand the
// delete off to a short-lived detached cmd that waits for us to exit first.
func removeSelf(path string) error {
	if runtime.GOOS != "windows" {
		return os.Remove(path)
	}

	// ping is a portable "sleep ~1s" so our process can exit and release the
	// file lock before del runs.
	script := fmt.Sprintf(`ping -n 2 127.0.0.1 >nul & del /f /q "%s"`, path)
	c := exec.Command("cmd", "/C", script)
	c.SysProcAttr = detachedProcAttr()
	if err := c.Start(); err != nil {
		return err
	}
	fmt.Println("  binary will be removed as pong-ball exits.")
	return nil
}
