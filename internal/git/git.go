package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// gitCmd builds an exec.Cmd for git with UTF-8 forced on every platform.
func gitCmd(args ...string) *exec.Cmd {
	fullArgs := append([]string{
		"-c", "core.quotepath=false",
		"-c", "i18n.logOutputEncoding=UTF-8",
		"-c", "i18n.commitEncoding=UTF-8",
	}, args...)

	cmd := exec.Command("git", fullArgs...)
	cmd.Env = append(os.Environ(),
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
		"GIT_TERMINAL_PROMPT=0",
	)
	return cmd
}

// Run executes a git command and returns trimmed stdout.
func Run(args ...string) (string, error) {
	cmd := gitCmd(args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// MustRun is like Run but returns an empty string on error.
func MustRun(args ...string) string {
	out, _ := Run(args...)
	return out
}

// Log returns commit messages for the last n commits.
func Log(n int) string {
	return MustRun("log", fmt.Sprintf("-n%d", n), "--pretty=format:Commit: %s\nDesc: %b\n")
}

// DiffStat returns the diff --stat for the last n commits.
func DiffStat(n int) string {
	ref := fmt.Sprintf("HEAD~%d", n)
	return MustRun("diff", "--stat", ref, "HEAD")
}

// Diff returns the full unified diff for the last n commits.
func Diff(n int) string {
	ref := fmt.Sprintf("HEAD~%d", n)
	return MustRun("diff", ref, "HEAD")
}

// FilteredDiff returns the diff for the last n commits, excluding files that
// match any of the given patterns (uses git pathspec :(exclude)).
// If ignorePatterns is empty, falls back to plain Diff.
func FilteredDiff(n int, ignorePatterns []string) string {
	ref := fmt.Sprintf("HEAD~%d", n)
	args := []string{"diff", ref, "HEAD", "--"}
	for _, p := range ignorePatterns {
		args = append(args, ":(exclude)"+p)
	}
	return MustRun(args...)
}

// StagedDiff returns the diff of currently staged (cached) changes.
// This is what prgen commit uses to analyze what's about to be committed.
func StagedDiff() string {
	return MustRun("diff", "--cached")
}

// StagedStat returns the --stat of currently staged changes.
func StagedStat() string {
	return MustRun("diff", "--cached", "--stat")
}

// DiffBetween returns the unified diff between two refs (branches, tags, SHAs).
// Used by --from / --to and prgen review.
func DiffBetween(from, to string) string {
	return MustRun("diff", from+"..."+to)
}

// StatBetween returns the --stat between two refs.
func StatBetween(from, to string) string {
	return MustRun("diff", "--stat", from+"..."+to)
}

// LogBetween returns the commit log between two refs.
func LogBetween(from, to string) string {
	return MustRun("log", from+".."+to, "--pretty=format:Commit: %s\nDesc: %b\n")
}

// Branch returns the current branch name.
func Branch() string {
	return MustRun("rev-parse", "--abbrev-ref", "HEAD")
}

// HeadHash returns the full HEAD commit hash.
func HeadHash() string {
	return MustRun("rev-parse", "HEAD")
}

// IsClean returns true if there are no uncommitted changes.
func IsClean() bool {
	out := MustRun("status", "--porcelain")
	return out == ""
}

// HasStagedChanges returns true if there are staged changes to commit.
func HasStagedChanges() bool {
	out := MustRun("diff", "--cached", "--name-only")
	return out != ""
}

// FetchAndDiff fetches origin and returns commits ahead of local HEAD.
func FetchAndDiff(remote, branch string) (string, error) {
	if _, err := Run("fetch", remote); err != nil {
		return "", fmt.Errorf("git fetch: %w", err)
	}
	ref := fmt.Sprintf("HEAD..%s/%s", remote, branch)
	log, err := Run("log", ref, "--oneline")
	if err != nil {
		return "", err
	}
	return log, nil
}

// Pull does git pull origin <branch>.
func Pull(remote, branch string) error {
	_, err := Run("pull", remote, branch)
	return err
}
