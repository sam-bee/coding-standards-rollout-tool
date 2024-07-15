package csfixing

import (
	"fmt"
	"io"
	"log"
	"testing"
)

func TestFix(t *testing.T) {

	// Set up application config file
	conf := getConfig()

	// Set up mock version control
	git := &gitTestDouble{}

	git.branchesToReturn = []string{"origin/main", "origin/feature/branch1", "origin/feature/branch2"}
	git.filesToReturn = []string{"file1.php", "file2.php"}

	// Set up mock system caller
	systemCaller := &systemCallerTestDouble{}

	// Null logger
	logger := log.New(io.Discard, "", 0)

	// Run the function
	Fix(conf, git, systemCaller, logger)

	// Check that the correct Git commands were run
	expectedCommands := []string{
		"git fetch origin",
		"git branch -r",
		"git diff --name-only  origin/main origin/main",
		"git diff --name-only  origin/feature/branch1 origin/main",
		"git diff --name-only  origin/feature/branch2 origin/main",
		"git checkout origin/main -- file1.php",
		"git checkout origin/main -- file2.php",
	}

	if len(git.commandsRun) != len(expectedCommands) {
		t.Errorf("Expected %d commands to be run, but got %d", len(expectedCommands), len(git.commandsRun))
		return
	}
	for i, cmd := range expectedCommands {
		if git.commandsRun[i] != cmd {
			t.Errorf("Expected command %s to be run, but got %s", cmd, git.commandsRun[i])
		}
	}

	// Check that the correct coding standards fixer command was run
	if systemCaller.commandRun != "/path/to/fixer" {
		t.Errorf("Expected command to be %s, but got %s", "/path/to/fixer", systemCaller.commandRun)
	}
	if len(systemCaller.argsRun) != 2 {
		t.Errorf("Expected 2 arguments to be passed to the command, but got %d", len(systemCaller.argsRun))
	}
	if systemCaller.argsRun[0] != "fixcommand" {
		t.Errorf("Expected first argument to be %s, but got %s", "fixcommand", systemCaller.argsRun[0])
	}
	if systemCaller.argsRun[1] != "--a-flag" {
		t.Errorf("Expected second argument to be %s, but got %s", "--a-flag", systemCaller.argsRun[1])
	}
}

func getConfig() ApplicationConfig {
	return BuildConfig(
		map[string]interface{}{
			"git": map[string]interface{}{
				"mainline-branch-name": "main",
				"remote-name":          "origin",
			},
			"codingstandards": map[string]interface{}{
				"command-to-run":    "/path/to/fixer",
				"command-arguments": []interface{}{"fixcommand", "--a-flag"},
			},
		},
	)
}

// gitTestDouble implements same interface as git struct in git.go. Used as a test double to test Fix function.

type gitTestDouble struct {
	filesToReturn    []string
	branchesToReturn []string
	commandsRun      []string
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

// systemCallerTestDouble implements same interface as SystemCaller struct in systemcall.go. Used as a test double to test Fix function.

type systemCallerTestDouble struct {
	commandRun string
	argsRun    []string
}

func (sc *systemCallerTestDouble) doSystemCall(command string, args []string) ([]string, error) {
	sc.commandRun = command
	sc.argsRun = args
	return []string{}, nil
}

func TestSystemCallerImplementsInterface(t *testing.T) {
	var _ systemCallerInterface = (*systemCallerTestDouble)(nil)
}
