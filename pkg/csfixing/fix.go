package csfixing

import (
	"log"
	"strings"
	"sync"
)

var logs *log.Logger

func Fix(conf ApplicationConfig, git gitInterface, exec systemCallerInterface, logger *log.Logger) {

	remoteName := conf.getRemoteName()
	mainlineTrackingBranch := remoteName + "/" + conf.getMainlineBranchName()
	logs = logger

	// main algorithm
	update(git, remoteName)
	trackingBranches := getTrackingBranches(git, remoteName, mainlineTrackingBranch)
	exemptFiles := getExemptFiles(git, trackingBranches, mainlineTrackingBranch)
	fixCodingStandards(exec, conf.getCommandToRun(), conf.getCommandArguments())
	revertChangesToFiles(git, mainlineTrackingBranch, exemptFiles)
}

func update(git gitInterface, remote string) {
	git.fetch(remote)
	logs.Printf("Fetching from remote %s\n", remote)
}

func getTrackingBranches(git gitInterface, remote string, mainBranch string) []string {
	allBranches, _ := git.getRemoteBranches()
	allBranches = trim(allBranches)
	branches := filterForRelevantTrackingBranches(allBranches, remote, mainBranch)

	logs.Printf("There are %d tracking branches starting with '%s' from the remote\n", len(branches), remote+"/")
	for i := 0; i < min(5, len(branches)); i++ {
		logs.Printf("  (example tracking branch %d: %s)\n", i, branches[i])
	}
	return branches
}

func getExemptFiles(git gitInterface, branches []string, mainBranch string) []string {
	exemptFiles := []string{}
	for _, trackingBranch := range branches {
		files, _ := git.getFilesEditedInBranch(trackingBranch, mainBranch)
		exemptFiles = append(exemptFiles, files...)
	}
	uniqueExemptFiles := unique(exemptFiles)

	logs.Printf("There are %d exempt files which should be reverted after coding standards fixes\n", len(uniqueExemptFiles))
	for i := 0; i < min(5, len(uniqueExemptFiles)); i++ {
		logs.Printf("  (example exempt file %d: %s)\n", i, uniqueExemptFiles[i])
	}

	return uniqueExemptFiles
}

func fixCodingStandards(exec systemCallerInterface, command string, args []string) {
	logs.Printf("Running coding standards fixer command: %s %s\n", command, strings.Join(args, " "))
	_, exitCode, _ := exec.doSystemCall(command, args)
	if exitCode != 0 {
		logs.Printf("Command failed with exit code %d. Your configured CS fixing command is probably not working.\n", exitCode)
	}
}

func revertChangesToFiles(git gitInterface, mainBranch string, files []string) {
	logs.Printf("Reverting changes to %d files\n", len(files))
	fileCh := make(chan string, len(files))
	for _, file := range files {
		fileCh <- file
	}
	runChangeRevertingWorkers(100, fileCh, git, mainBranch)
}

func runChangeRevertingWorkers(workers int, fileCh chan string, git gitInterface, mainBranch string) {
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			fileChangeRevertingWorker(fileCh, git, mainBranch)
		}()
	}
	close(fileCh)
	wg.Wait()
}

func fileChangeRevertingWorker(fileCh <-chan string, git gitInterface, mainBranch string) {
	for file := range fileCh {
		git.revertChangesToFile(mainBranch, file)
	}
}

func filterForRelevantTrackingBranches(allBranches []string, remote string, mainBranch string) (ret []string) {
	branches := filter(
		allBranches,
		func(branch string) bool {
			if !strings.HasPrefix(branch, remote+"/") {
				return false
			}
			if strings.HasPrefix(branch, remote+"/HEAD ->") {
				return false
			}
			if branch == mainBranch {
				return false
			}
			return true
		},
	)
	return branches
}

func filter(strs []string, test func(string) bool) (filtered []string) {
	for _, s := range strs {
		if test(s) {
			filtered = append(filtered, s)
		}
	}
	return
}

func unique(in []string) (unique []string) {
	m := map[string]struct{}{}
	for _, s := range in {
		m[s] = struct{}{}
	}
	for s := range m {
		unique = append(unique, s)
	}
	return
}

func trim(s []string) (trimmed []string) {
	for _, str := range s {
		trimmed = append(trimmed, strings.TrimSpace(str))
	}
	return
}
