package csfixing

import (
	"strings"
)

func Fix(conf ApplicationConfig, git gitInterface, systemCaller systemCallerInterface) {

	remoteName := conf.getRemoteName()
	mainlineBranchName := conf.getMainlineBranchName()
	mainlineTrackingBranch := remoteName + "/" + mainlineBranchName

	git.fetch(remoteName)

	allBranches, _ := git.getRemoteBranches()
	trackingBranches := filterForRelevantTrackingBranches(allBranches, remoteName)

	exemptFiles := []string{}

	for _, trackingBranch := range trackingBranches {
		files, _ := git.getFilesEditedInBranch(trackingBranch, mainlineTrackingBranch)
		exemptFiles = append(exemptFiles, files...)
	}

	uniqueExemptFiles := unique(exemptFiles)

	systemCaller.doSystemCall(conf.getCommandToRun(), conf.getCommandArguments())

	for _, file := range uniqueExemptFiles {
		git.revertChangesToFile(mainlineTrackingBranch, file)
	}
}

func filterForRelevantTrackingBranches(allBranches []string, remoteName string) (ret []string) {
	for _, branch := range allBranches {
		if strings.HasPrefix(branch, remoteName + "/") {
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
