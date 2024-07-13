package csfixing

type ApplicationConfig struct {
	git struct {
		mainlineBranchName string
		remoteName         string
	}
	codingStandards struct {
		commandToRun string
	}
}

func BuildConfig(conf map[string]interface{}) ApplicationConfig {
	c := ApplicationConfig{}
	c.codingStandards.commandToRun = conf["codingstandards"].(map[string]interface{})["command-to-run"].(string)
	c.git.mainlineBranchName = conf["git"].(map[string]interface{})["mainline-branch-name"].(string)
	c.git.remoteName = conf["git"].(map[string]interface{})["remote-name"].(string)
	c.validateConfig()
	return c
}

func (c *ApplicationConfig) getMainlineBranchName() string {
	return c.git.mainlineBranchName
}
func (c *ApplicationConfig) getRemoteName() string {
	return c.git.remoteName
}
func (c *ApplicationConfig) getCommandToRun() string {
	return c.codingStandards.commandToRun
}

func (c *ApplicationConfig) validateConfig() {
	if c.git.mainlineBranchName == "" {
		panic("Config error: git.mainline-branch-name not set")
	}
	if c.git.remoteName == "" {
		panic("Config error: git.remote-name not set")
	}
	if c.codingStandards.commandToRun == "" {
		panic("Config error: codingstandards.command-to-run not set")
	}
}
