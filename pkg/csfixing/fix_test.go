package csfixing

import (
	"fmt"
	"testing"
)



func TestFix(t *testing.T) {
	conf := getConfig()
	git := &gitTestDouble{}
	Fix(conf, git)
}

func getConfig() ApplicationConfig {
	return BuildConfig(
		map[string]interface{}{
			"git": map[string]interface{}{
				"mainline-branch-name": "main",
				"remote-name":          "origin",
			},
			"codingstandards": map[string]interface{}{
				"command-to-run": "echo 'hello i am a command to fix coding standards in php or something'",
			},
		},
	)
}

// gitTestDouble implements same interface as git struct in git.go. Used as a test double to test Fix function.

type gitTestDouble struct {
	filesToReturn []string
	branchesToReturn []string
	commandsRun []string
}

func (g *gitTestDouble) fetch(remoteName string) error {
	g.commandsRun = append(g.commandsRun, fmt.Sprintf("git fetch %s", remoteName))
	return nil
}

func (g *gitTestDouble) getRemoteBranches() ([]string, error) {
	g.commandsRun = append(g.commandsRun, "git branch -r")
	return g.branchesToReturn, nil
}

func (g *gitTestDouble) getFilesEditedInBranch(featureTrackingBranch string, mainlineTrackingBranch string) ([]string, error) {
	g.commandsRun = append(g.commandsRun, fmt.Sprintf("git diff --name-only  %s %s", featureTrackingBranch, mainlineTrackingBranch))
	return g.filesToReturn, nil
}

func (g *gitTestDouble) revertChangesToFile(mainlineTrackingBranch, file string) error {
	g.commandsRun = append(g.commandsRun, fmt.Sprintf("git checkout %s -- %s", mainlineTrackingBranch, file))
	return nil
}

func TestGitMockImplementsInterface(t *testing.T) {
	var _ gitInterface = (*gitTestDouble)(nil)
}
