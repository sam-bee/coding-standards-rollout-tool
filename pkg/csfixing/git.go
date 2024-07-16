package csfixing

import (
	"os/exec"
	"strings"
)

type gitInterface interface {
	fetch(remoteName string) error
	getRemoteBranches() ([]string, error)
	getFilesEditedInBranch(featureTrackingBranch string, mainlineTrackingBranch string) ([]string, error)
	revertChangesToFile(mainlineTrackingBranch string, file string) error
}

type Git struct{}

/**
 * Runs something like `git fetch origin`
 * This fetches the latest changes from the remote
 */
func (g *Git) fetch(remote string) error {
	_, err := issueCommand("git", []string{"fetch", remote})
	return err
}

/**
 * Runs `git branch -r`
 * This gets a list of remote branches
 */
func (g *Git) getRemoteBranches() ([]string, error) {
	return issueCommand("git", []string{"branch", "-r"})
}

/**
 * Runs something like `git diff --name-only origin/feature-branch-name origin/main`
 * This gets a list of files edited in a branch
 */
func (g *Git) getFilesEditedInBranch(branch string, basisBranch string) ([]string, error) {
	return issueCommand("git", []string{"diff", "--name-only", branch, basisBranch})
}

/**
 * Runs something like `git checkout origin/main -- ./path/to/file.php`
 * This reverts local changes to a file
 */
func (g *Git) revertChangesToFile(mainBranch, file string) error {
	_, err := issueCommand("git", []string{"checkout", mainBranch, "--", file})
	return err
}

func issueCommand(command string, args []string) ([]string, error) {
	cmd := exec.Command(command, args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	return lines, nil
}
