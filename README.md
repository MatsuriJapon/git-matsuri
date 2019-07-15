[![CircleCI](https://circleci.com/gh/MatsuriJapon/git-matsuri.svg?style=svg&circle-token=d8dafe865a3b6e3dee7dd0f56e51e393f4ffc0be)](https://circleci.com/gh/MatsuriJapon/git-matsuri) [![GoDoc](https://godoc.org/github.com/MatsuriJapon/git-matsuri?status.svg)](https://godoc.org/github.com/MatsuriJapon/git-matsuri) [![Go Report Card](https://goreportcard.com/badge/github.com/MatsuriJapon/git-matsuri)](https://goreportcard.com/report/github.com/MatsuriJapon/git-matsuri)

# git-matsuri
Git subcommands to use with repositories under the MatsuriJapon Organization.

## Prerequisite
[Git](https://git-scm.com/downloads) must be installed.

## Installation
Download and extract the latest [release](https://github.com/MatsuriJapon/git-matsuri/releases) or compile from source (requires Go 1.12+). The binary is found in the respective directory for the platform of choice (`windows_amd64` for Windows, `darwin_amd64` for Mac OS, etc.).

### Windows
- Place `git-matsuri.exe` into the `C:\Program Files\Git\mingw64\libexec\git-core` folder

### Linux/Mac
- Move the `git-matsuri` binary file to a location of your choice, or keep it in the downloaded directory
- Add the file location to PATH: 
```sh
export PATH="$(pwd):$PATH"
```
- To make this permanent, add the above to your `~/.bashrc` or `~/.bash_profile` file and reload it using `source ~/.bashrc` or `source ~/.bash_profile`

### GitHub token
Visit https://github.com/settings/tokens/new and create a new token with `repo` and `user:email` permissions and save it to your system environment variables under the name `MATSURI_TOKEN`

#### Windows 10
`Win + S` and search for `environment variables`. Add one named `MATSURI_TOKEN` with the token you created as a value.

#### Linux/Mac
```sh
export MATSURI_TOKEN=<enter token here>
```

To make this permanent, add the above to your `~/.bashrc` or `~/.bash_profile` file and reload it using `source ~/.bashrc` or `source ~/.bash_profile`

## Clone a MatsuriJapon repository
To start working on a repository, you must first clone it. By default, we use git over SSH, although it is also possible to use it via HTTP (not recommended).
```sh
git matsuri setup ${REPO_NAME}
# Alternatively, you can clone using HTTP (not recommended)
git matsuri setup -http ${REPO_NAME}
# Move into the newly cloned directory
cd ${REPO_NAME}
```

## Show the current kanban
Use sparingly. It is usually meant for admins to prepare their report. Displays the full kanban, if available, in text format as would otherwise be available in the GitHub Projects page.
```sh
git matsuri kanban ${YEAR}
```

### Plain git equivalent
The current kanban can only be viewed on [GitHub](https://github.com/MatsuriJapon/matsuri-japon/projects)

## Show open issues
Usually used by specifying a year to get the open issues for the current Project year.
```sh
git matsuri todo
git matsuri todo ${YEAR}
```

### Plain git equivalent
The current kanban can only be viewed on [GitHub](https://github.com/MatsuriJapon/matsuri-japon/projects)

## Start working on an issue
```sh
git matsuri start
git matsuri start ${ISSUE}
```

### Plain git equivalent
Suppose the curent project year is 2020, then the default branch will be `v2020` for the `matsuri-japon` repository. For other repositories, check what the default branch is on GitHub (it is usually `master`).
```sh
# Assuming you have already cloned the repository
git checkout v2020
git pull
git checkout -b ISSUE-${ISSUE}
```

## Save current work to GitHub in a topic branch
```sh
# first commit your work
git add <modified files>
git commit -m "<a meaningful commit message>"
git matsuri save ${ISSUE}
# for subsequent saves, a simple `git push` would do
git push
```

### Plain git equivalent
```sh
# first commit your work
git add <modified files>
git commit -m "<a meaningful commit message>"
git push -u origin ISSUE-${ISSUE}
# for subsequent saves, a simple `git push` would do
git push
```

## Create a pull request
When your work is ready for review, send a Pull Request specifying the issue number. If the Pull Request is not meant to close the issue once it is merged, add the `-noclose` flag
```sh
git matsuri pr ${ISSUE}
git matsuri pr -noclose ${ISSUE}
```

### Plain git equivalent
Go to the repository page on GitHub and manually create a new Pull Request via the GUI. The PR title must start with `ISSUE-XYZ` where XYZ is the issue number you were working on, and the PR message must contain a message like `Closes #XYZ` for the Issue to be automatically closed when the PR is merged (we usually want that).

## Fix your pull request
When a reviewer requests changes to your PR, you can simply make the requested changes and push to the topic branch.
```sh
git add <modified files>
git commit -m "<a meaningful commit message>"
git push
```

## Create a fix pull request
If there was a problem with an already merged pull request, instead create a fix PR.
```sh
git matsuri fix ${ISSUE}
```

### Plain git equivalent
Go to the repository page on GitHub and manually create a new Pull Request via the GUI. The PR title must start with `ISSUE-XYZ-fix` where XYZ is the issue number you were working on, and the PR message must contain a message like `Closes #XYZ` for the Issue to be automatically closed when the PR is merged (we usually want that).

## Rebasing
When too many commits have been added to the PR, the reviewer may request you squash them into a single commit to avoid polluting the log. For example, if you made 16 commits in a PR:
```sh
git rebase -i HEAD~16
```
- In the editor that appears (usually `vim`), leave the first line intact, then replace "pick" in all the following lines with "fix" or "f"
- Exit the editor (usually `vim`) by using the command `:wq`
- In case the editor that opens up is `nano`, save and exit by using `ctrl+x`

After rebasing, you may have to push using the `--force` option
```sh
git push --force
```
