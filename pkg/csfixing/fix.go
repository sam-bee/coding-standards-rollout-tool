package csfixing

import (
	"log"
	"strings"
	"sync"
)

var logs *log.Logger

func Fix(conf ApplicationConfig, git gitInterface, systemCaller systemCallerInterface, logger *log.Logger) {

	remoteName := conf.getRemoteName()
	mainlineTrackingBranch := remoteName + "/" + conf.getMainlineBranchName()
	logs = logger

	// main algorithm
	update(git, remoteName)
	trackingBranches := getTrackingBranches(git, remoteName, mainlineTrackingBranch)
	exemptFiles := getExemptFiles(git, trackingBranches, mainlineTrackingBranch)
	systemCaller.doSystemCall(conf.getCommandToRun(), conf.getCommandArguments())
	revertChangesToFiles(git, mainlineTrackingBranch, exemptFiles)
}

func update(git gitInterface, remoteName string) {
	git.fetch(remoteName)
	logs.Printf("Fetching from remote %s\n", remoteName)
}

func getTrackingBranches(git gitInterface, remoteName string, mainlineTrackingBranch string) []string {
	allBranches, _ := git.getRemoteBranches()
	allBranches = trim(allBranches)
	filteredBranches := filterForRelevantTrackingBranches(allBranches, remoteName, mainlineTrackingBranch)

	logs.Printf("There are %d tracking branches starting with '%s' from the remote\n", len(filteredBranches), remoteName+"/")
	for i := 0; i < min(5, len(filteredBranches)); i++ {
		logs.Printf("  (Example tracking branch %d: %s)\n", i, filteredBranches[i])
	}
	return filteredBranches
}

func getExemptFiles(git gitInterface, trackingBranches []string, mainlineTrackingBranch string) []string {
	exemptFiles := []string{}
	for _, trackingBranch := range trackingBranches {
		files, _ := git.getFilesEditedInBranch(trackingBranch, mainlineTrackingBranch)
		exemptFiles = append(exemptFiles, files...)
	}
	uniqueExemptFiles := unique(exemptFiles)

	logs.Printf("There are %d exempt files which should be reverted after coding standards fixes\n", len(uniqueExemptFiles))
	for i := 0; i < min(5, len(uniqueExemptFiles)); i++ {
		logs.Printf("  (Example exempt file %d: %s)\n", i, uniqueExemptFiles[i])
	}

	return uniqueExemptFiles
}

func revertChangesToFiles(git gitInterface, mainlineTrackingBranch string, files []string) {
	fileCh := make(chan string, len(files))
	for _, file := range files {
		fileCh <- file
	}
	runChangeRevertingWorkers(100, fileCh, git, mainlineTrackingBranch)
}

func runChangeRevertingWorkers(noOfWorkers int, fileCh chan string, git gitInterface, mainlineTrackingBranch string) {
	wg := sync.WaitGroup{}
	wg.Add(noOfWorkers)
	for i := 0; i < noOfWorkers; i++ {
		go func() {
			defer wg.Done()
			fileChangeRevertingWorker(fileCh, git, mainlineTrackingBranch)
		}()
	}
	close(fileCh)
	wg.Wait()
}

func fileChangeRevertingWorker(fileCh <-chan string, git gitInterface, mainlineTrackingBranch string) {
	for file := range fileCh {
		git.revertChangesToFile(mainlineTrackingBranch, file)
	}
}

func filterForRelevantTrackingBranches(allBranches []string, remoteName string, mainlineTrackingBranch string) (ret []string) {
	branches := filter(
		allBranches,
		func(branch string) bool {
			if !strings.HasPrefix(branch, remoteName+"/") {
				return false
			}
			if strings.HasPrefix(branch, remoteName+"/HEAD ->") {
				return false
			}
			if branch == mainlineTrackingBranch {
				return false
			}
			return true
		},
	)
	return branches
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func unique(in []string) (ret []string) {
	m := map[string]struct{}{}
	for _, s := range in {
		m[s] = struct{}{}
	}
	for s := range m {
		ret = append(ret, s)
	}
	return
}

func trim(s []string) (ret []string) {
	for _, str := range s {
		ret = append(ret, strings.TrimSpace(str))
	}
	return
}
