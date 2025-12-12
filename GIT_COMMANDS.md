# Basic Git Commands Guide

This guide explains the basic Git commands that are essential for version control and collaboration. Git is a distributed version control system that helps you track changes in your code and collaborate with others.

## Table of Contents
- [git init](#git-init)
- [git clone](#git-clone)
- [git status](#git-status)
- [git add](#git-add)
- [git commit](#git-commit)
- [git push](#git-push)
- [git pull](#git-pull)

---

## git init

**Purpose:** Initialize a new Git repository in the current directory.

**Description:** The `git init` command creates a new Git repository. It creates a `.git` directory with subdirectories for objects, refs/heads, refs/tags, and template files. This command is used when you want to start tracking an existing project with Git.

**Syntax:**
```bash
git init [directory]
```

**Examples:**

1. Initialize a repository in the current directory:
```bash
git init
```
Output:
```
Initialized empty Git repository in /path/to/your/project/.git/
```

2. Initialize a repository in a new directory:
```bash
git init my-new-project
```
Output:
```
Initialized empty Git repository in /path/to/my-new-project/.git/
```

3. Initialize a bare repository (for use as a remote):
```bash
git init --bare my-repo.git
```

**When to use:**
- Starting a new project that you want to track with Git
- Converting an existing project to use Git version control

---

## git clone

**Purpose:** Create a local copy of a remote repository.

**Description:** The `git clone` command copies an entire repository from a remote location (like GitHub, GitLab, or another server) to your local machine. It automatically sets up a remote connection called "origin" pointing to the source repository.

**Syntax:**
```bash
git clone <repository-url> [directory]
```

**Examples:**

1. Clone a repository using HTTPS:
```bash
git clone https://github.com/deep-diver/hf-daily-paper-newsletter.git
```
Output:
```
Cloning into 'hf-daily-paper-newsletter'...
remote: Enumerating objects: 1234, done.
remote: Counting objects: 100% (1234/1234), done.
remote: Compressing objects: 100% (567/567), done.
remote: Total 1234 (delta 890), reused 1100 (delta 756)
Receiving objects: 100% (1234/1234), 2.50 MiB | 3.00 MiB/s, done.
Resolving deltas: 100% (890/890), done.
```

2. Clone a repository using SSH:
```bash
git clone git@github.com:username/repository.git
```

3. Clone into a specific directory:
```bash
git clone https://github.com/username/repo.git my-local-folder
```

4. Clone only the latest commit (shallow clone):
```bash
git clone --depth 1 https://github.com/username/repo.git
```

**When to use:**
- Starting work on an existing project
- Creating a local copy of a repository for development
- Downloading open-source projects

---

## git status

**Purpose:** Display the state of the working directory and staging area.

**Description:** The `git status` command shows which files have been modified, which files are staged for commit, and which files are not being tracked by Git. It's one of the most frequently used commands as it helps you understand the current state of your repository.

**Syntax:**
```bash
git status [options]
```

**Examples:**

1. Check the status of your repository:
```bash
git status
```
Output example:
```
On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   README.md

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	new-feature.js

no changes added to commit (use "git add" and/or "git commit -a")
```

2. Get a shorter, more compact status:
```bash
git status -s
```
Output example:
```
 M README.md
?? new-feature.js
```

3. Show branch and tracking information:
```bash
git status -sb
```
Output example:
```
## main...origin/main
 M README.md
?? new-feature.js
```

**When to use:**
- Before staging files to see what has changed
- Before committing to verify what will be included
- To check if your branch is up to date with the remote
- To see which files are untracked

---

## git add

**Purpose:** Add file contents to the staging area (index).

**Description:** The `git add` command stages changes for the next commit. When you modify files, Git doesn't automatically include them in your next commit. You must explicitly stage the changes you want to commit using `git add`.

**Syntax:**
```bash
git add <pathspec>...
```

**Examples:**

1. Stage a specific file:
```bash
git add README.md
```

2. Stage multiple specific files:
```bash
git add file1.js file2.js file3.css
```

3. Stage all changes in the current directory and subdirectories:
```bash
git add .
```

4. Stage all changes in the entire repository:
```bash
git add -A
```
or
```bash
git add --all
```

5. Stage all modified and deleted files (but not new files):
```bash
git add -u
```

6. Stage files interactively (choose which changes to stage):
```bash
git add -p
```
This will prompt you for each change:
```
Stage this hunk [y,n,q,a,d,e,?]?
```

7. Stage all files with a specific extension:
```bash
git add *.js
```

**When to use:**
- After creating new files that you want to track
- After modifying files that you want to commit
- To selectively choose which changes to include in a commit

---

## git commit

**Purpose:** Record changes to the repository.

**Description:** The `git commit` command saves your staged changes to the local repository. Each commit creates a snapshot of your project at that point in time, along with a message describing the changes.

**Syntax:**
```bash
git commit [options]
```

**Examples:**

1. Commit with a message:
```bash
git commit -m "Add user authentication feature"
```
Output:
```
[main 3a8f7b9] Add user authentication feature
 3 files changed, 45 insertions(+), 2 deletions(-)
 create mode 100644 auth.js
```

2. Commit with a multi-line message:
```bash
git commit -m "Add user authentication feature" -m "This commit adds login and logout functionality with JWT tokens."
```

3. Stage all tracked files and commit (skip `git add` for modified files):
```bash
git commit -a -m "Update documentation"
```
or
```bash
git commit -am "Update documentation"
```

4. Commit and open your default editor for a detailed message:
```bash
git commit
```

5. Amend the previous commit (add more changes or edit the message):
```bash
git add forgotten-file.js
git commit --amend -m "Updated commit message"
```

6. Commit with a detailed message template:
```bash
git commit -m "feat: add user profile page

- Created user profile component
- Added profile editing functionality
- Implemented avatar upload feature

Resolves: #123"
```

**Best practices for commit messages:**
- Use present tense: "Add feature" not "Added feature"
- Keep the first line under 50 characters
- Add a blank line before a detailed description
- Explain what and why, not how

**When to use:**
- After staging changes that represent a logical unit of work
- When you've completed a feature, bug fix, or other meaningful change
- Regularly to create a history of your project's evolution

---

## git push

**Purpose:** Upload local repository content to a remote repository.

**Description:** The `git push` command transfers commits from your local repository to a remote repository. This is how you share your work with others and back up your code to a remote server.

**Syntax:**
```bash
git push [remote] [branch]
```

**Examples:**

1. Push commits to the default remote (origin) and branch:
```bash
git push
```
Output:
```
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Delta compression using up to 8 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 356 bytes | 356.00 KiB/s, done.
Total 3 (delta 1), reused 0 (delta 0)
To https://github.com/username/repository.git
   a1b2c3d..e4f5g6h  main -> main
```

2. Push to a specific remote and branch:
```bash
git push origin main
```

3. Push and set upstream tracking (first time pushing a new branch):
```bash
git push -u origin feature-branch
```
or
```bash
git push --set-upstream origin feature-branch
```

4. Push all branches:
```bash
git push --all origin
```

5. Push tags to remote:
```bash
git push --tags
```

6. Force push (use with caution - overwrites remote history):
```bash
git push --force
```
or safer alternative:
```bash
git push --force-with-lease
```

7. Push to a different branch name on remote:
```bash
git push origin local-branch:remote-branch
```

**When to use:**
- After committing changes that you want to share with others
- To back up your work to a remote server
- To deploy code to a production server (in some workflows)
- Before creating a pull request

**Important notes:**
- Always pull before pushing to avoid conflicts
- Be careful with force push as it can overwrite others' work
- Make sure you have permission to push to the remote repository

---

## git pull

**Purpose:** Fetch from and integrate with another repository or local branch.

**Description:** The `git pull` command downloads content from a remote repository and immediately updates the local repository to match that content. It's essentially a combination of `git fetch` (download changes) and `git merge` (integrate changes).

**Syntax:**
```bash
git pull [remote] [branch]
```

**Examples:**

1. Pull changes from the default remote branch:
```bash
git pull
```
Output:
```
remote: Enumerating objects: 5, done.
remote: Counting objects: 100% (5/5), done.
remote: Compressing objects: 100% (3/3), done.
remote: Total 3 (delta 1), reused 0 (delta 0)
Unpacking objects: 100% (3/3), done.
From https://github.com/username/repository
   a1b2c3d..e4f5g6h  main       -> origin/main
Updating a1b2c3d..e4f5g6h
Fast-forward
 README.md | 5 ++++-
 1 file changed, 4 insertions(+), 1 deletion(-)
```

2. Pull from a specific remote and branch:
```bash
git pull origin main
```

3. Pull with rebase instead of merge:
```bash
git pull --rebase
```

4. Pull from a specific remote branch to your current branch:
```bash
git pull origin feature-branch
```

5. Pull and automatically stash and reapply local changes:
```bash
git pull --autostash
```

6. Pull without committing the merge (to review changes first):
```bash
git pull --no-commit
```

**When to use:**
- Before starting new work to ensure you have the latest changes
- Before pushing your changes to avoid conflicts
- When collaborating with others to sync their changes
- Regularly to keep your local repository up to date

**Common scenarios:**

1. **Fast-forward pull** (no conflicts):
```bash
$ git pull
Updating a1b2c3d..e4f5g6h
Fast-forward
 file.txt | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
```

2. **Pull with automatic merge**:
```bash
$ git pull
Auto-merging file.txt
Merge made by the 'recursive' strategy.
 file.txt | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
```

3. **Pull with conflicts** (requires manual resolution):
```bash
$ git pull
Auto-merging file.txt
CONFLICT (content): Merge conflict in file.txt
Automatic merge failed; fix conflicts and then commit the result.
```

**Best practices:**
- Pull frequently to stay synchronized with the team
- Commit your changes before pulling to avoid losing work
- Use `git status` after pulling to see what changed
- Consider using `git pull --rebase` for a cleaner history

---

## Common Workflow Example

Here's a typical workflow combining these commands:

```bash
# 1. Clone a repository (first time only)
git clone https://github.com/username/project.git
cd project

# 2. Check current status
git status

# 3. Pull latest changes (before starting work)
git pull

# 4. Make changes to files...
# (edit README.md, add new files, etc.)

# 5. Check what changed
git status

# 6. Stage your changes
git add README.md
git add new-feature.js

# Or stage all changes
git add .

# 7. Commit your changes
git commit -m "Add new feature and update README"

# 8. Push your changes to remote
git push

# 9. Verify everything is clean
git status
```

## Additional Tips

1. **Check your Git version:**
```bash
git --version
```

2. **Get help for any command:**
```bash
git help <command>
# or
git <command> --help
```

3. **Configure Git (first time setup):**
```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

4. **View commit history:**
```bash
git log
# or for a condensed view
git log --oneline
```

5. **See differences:**
```bash
# See unstaged changes
git diff

# See staged changes
git diff --staged
```

## Conclusion

These basic Git commands form the foundation of version control workflows. Mastering them will help you:
- Track changes in your projects effectively
- Collaborate with team members
- Maintain a history of your work
- Recover from mistakes
- Deploy code safely

Remember to commit often, write clear commit messages, and pull before you push to ensure smooth collaboration!
