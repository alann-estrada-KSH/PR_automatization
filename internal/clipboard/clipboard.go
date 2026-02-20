package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Copy puts text into the system clipboard.
// Works on Windows (via PowerShell), macOS (pbcopy), and Linux (xclip/xsel).
func Copy(text string) error {
	switch runtime.GOOS {
	case "windows":
		return copyWindows(text)
	case "darwin":
		return copyMac(text)
	default:
		return copyLinux(text)
	}
}

// copyWindows uses PowerShell Set-Clipboard which correctly handles UTF-8 /
// Unicode. Using clip.exe directly mangles non-ASCII characters because it
// reads stdin in the system code page (cp1252 on most Spanish Windows installs).
func copyWindows(text string) error {
	// Write to a temp file so we avoid shell-escaping nightmares with the
	// content (special chars, quotes, backticks, etc.)
	tmp, err := os.CreateTemp("", "prgen-clipboard-*.txt")
	if err != nil {
		return fmt.Errorf("clipboard: creating temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(text); err != nil {
		tmp.Close()
		return fmt.Errorf("clipboard: writing temp file: %w", err)
	}
	tmp.Close()

	// PowerShell reads the file as UTF-8 and places it in the clipboard as Unicode.
	script := fmt.Sprintf(
		"[System.IO.File]::ReadAllText('%s', [System.Text.Encoding]::UTF8) | Set-Clipboard",
		tmp.Name(),
	)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("clipboard: PowerShell Set-Clipboard failed: %w\n%s", err, out)
	}
	return nil
}

func copyMac(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func copyLinux(text string) error {
	// Try xclip first, then xsel
	for _, tool := range []string{"xclip -selection clipboard", "xsel --clipboard --input"} {
		parts := strings.Fields(tool)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	return nil // Silently fail if neither tool is available
}
