package csfixing

import (
	"strings"
	"sync"
)

func Fix(conf ApplicationConfig, git gitInterface, systemCaller systemCallerInterface) {

	remoteName := conf.getRemoteName()
	mainlineBranchName := conf.getMainlineBranchName()
	mainlineTrackingBranch := remoteName + "/" + mainlineBranchName

	git.fetch(remoteName)
	trackingBranches := getTrackingBranches(git, remoteName)
	exemptFiles := getExemptFiles(git, trackingBranches, mainlineTrackingBranch)
	systemCaller.doSystemCall(conf.getCommandToRun(), conf.getCommandArguments())
	revertChangesToFiles(git, mainlineTrackingBranch, exemptFiles)
}

func getExemptFiles(git gitInterface, trackingBranches []string, mainlineTrackingBranch string) []string {
	exemptFiles := []string{}
	for _, trackingBranch := range trackingBranches {
		files, _ := git.getFilesEditedInBranch(trackingBranch, mainlineTrackingBranch)
		exemptFiles = append(exemptFiles, files...)
	}
	uniqueExemptFiles := unique(exemptFiles)
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

func getTrackingBranches(git gitInterface, remoteName string) []string {
	allBranches, _ := git.getRemoteBranches()
	return filterForRelevantTrackingBranches(allBranches, remoteName)
}

func filterForRelevantTrackingBranches(allBranches []string, remoteName string) (ret []string) {
	for _, branch := range allBranches {
		if strings.HasPrefix(branch, remoteName+"/") {
			ret = append(ret, branch)
		}
	}
	return
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
