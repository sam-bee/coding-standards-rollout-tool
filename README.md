# Coding Standards Rollout Tool

## Introduction

This tool is for rolling out coding standards changes on software projects. It is for large projects with multiple
developers working on them. The tool will fix coding standards in all your source code files, _except_ those where
changes exist in any branch. If someone is making changes to a file in a feature branch, that file will not
be changed. This prevents merge conflicts.

It is written with languages like PHP in mind, where things like whitespace aren't enforced by the language,
but accepted standards exist. Using this tool is a step towards rolling out or changing coding standards in a
large legacy project. You may need to run the tool again a couple of weeks later to achieve near-total adoption
of your coding standards. A possible future step would be for you to add a CS fixer to your build server.

## Suggested Process

To roll out new or different coding standards in your project, consider using a process like the following:
1. Clean up any old branches you don't need on your project.
1. Use this tool to change the majority of your source code. Every file will be fixed unless a feature branch exists where that file has been edited.
1. Get the CS changes merged quickly, before too much work happens on other branches.
1. Ensure every feature branch in your project is updated immediately. Rebasing/updating now prevents massive merge conflicts later
1. Wait a couple of weeks, then do the same thing again. This will fix a few more files, where feature development work was in progress last time.
1. Add a CS fixer or some sort of static analysis check to your build process, so that in future the project will
adhere to your coding standards.

## Configuration

The tool is language agnostic, so you will need to provide working shell commands that run your choice of CS fixer. You will need a config file similar to `config/config.toml.example`. Provide the location of the config to the tool with the `--config` flag. Here is an example of its contents:

```toml
[git]
mainline-branch-name = 'main' # typically 'master' or 'main'
remote-name = 'origin'

[codingstandards]
# fixer command appropriate to your chosen language/tool
command-to-run = './vendor/bin/php-cs-fixer fix ./src/ --rules=@PSR12'
```

## What is it Doing?

The tool never pushes anything, but it does mess around running Git commands in your folder. You could run it in
a clean clone folder if you wanted to, but it's not strictly necessary. This is what the code is doing in detail:

- `git fetch` to update things
- `git branch --remote` to get a list of all the feature branches on the remote
- For each branch on the remote: `git diff --name-only origin/master origin/featurebranchname` to get a list of
'exempt' files someone is working on in that branch. (Main branch name and remote name configurable.)
- Command to fix coding standards in your project (configurable)
- For each file that is 'exempt' `git checkout origin/master -- [filename]` to undo local changes

Note that it didn't commit anything - you will need to do that yourself.

## Running

[]