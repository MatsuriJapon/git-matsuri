[![CircleCI](https://circleci.com/gh/MatsuriJapon/git-matsuri.svg?style=svg&circle-token=d8dafe865a3b6e3dee7dd0f56e51e393f4ffc0be)](https://circleci.com/gh/MatsuriJapon/git-matsuri)

# git-matsuri
Git subcommands to use with repositories under the MatsuriJapon Organization.

## Prerequisite
[Git](https://git-scm.com/downloads) must be installed.

## Installation
Download and extract the latest release. The binary is found in the respective directory for the platform of choice (windows_amd64 for Windows, darwin_amd64 for Mac OS, etc.).

### Windows
- Place `git-matsuri.exe` into the `C:\Program Files\Git\mingw64\libexec\git-core` folder

### Linux/Mac
- Move the `git-matsuri` binary file to a location of your choice, or keep it in the downloaded directory
- Add the file location to PATH: 
```sh
export PATH="$(pwd):$PATH"
```
- To make this permanent, add the above to your `~/.bashrc` file

## Show the current kanban
Use sparingly. It is usually meant for admins to prepare their report. Displays the full kanban, if available, in text format as would otherwise be available in the GitHub Projects page.
```sh
git matsuri kanban ${YEAR}
```

## Show open issues
Usually used by specifying a year to get the open issues for the current Project year.
```sh
git matsuri todo
git matsuri todo ${YEAR}
```

## Start working on an issue
```sh
git matsuri start
git matsuri start ${ISSUE}
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

## Create a pull request
When your work is ready for review, send a Pull Request specifying the issue number. If the Pull Request is not meant to close the issue once it is merged, add the `-noclose` flag
```sh
git matsuri pr ${ISSUE}
git matsuri pr -noclose ${ISSUE}
```

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