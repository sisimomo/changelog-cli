package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Commit represents a parsed conventional commit.
type Commit struct {
	Type        string
	Scope       string
	Description string
	Body        string
	Date        time.Time
}

var headerRegex = regexp.MustCompile(`^([a-zA-Z]+)(?:\(([^)]+)\))?:\s*(.*)$`)

// IsGitRepository returns an error if path is not a valid git repo.
func IsGitRepository(path string) error {
	out, err := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return fmt.Errorf("failed to execute git command for path %q: %w", path, err)
	}
	if strings.TrimSpace(string(out)) != "true" {
		return fmt.Errorf("path %q is not a git repository", path)
	}
	return nil
}

// TagExists checks if a tag exists in the repo.
func TagExists(path, tag string) (bool, error) {
	_, err := exec.Command("git", "-C", path, "show-ref", "--tags", tag).Output()
	if err == nil {
		return true, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("failed to check tag existence for %q: %w", tag, err)
}

// GetLastTwoTags returns the last and previous tag, sorted by commit date.
func GetLastTwoTags(path string) (last string, previous *string, err error) {
	out, err := exec.Command("git", "-C", path, "tag", "--sort=-committerdate").Output()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get tags: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(tags) == 0 || tags[0] == "" {
		return "", nil, fmt.Errorf("no tags found")
	}

	last = tags[0]
	if len(tags) > 1 {
		prev := tags[1]
		previous = &prev
	}

	return last, previous, nil
}

// GetInitialCommit returns the hash of the first commit in the repository.
func GetInitialCommit(path string) (string, error) {
	out, err := exec.Command("git", "-C", path, "rev-list", "--max-parents=0", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get initial commit: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetParsedCommits returns all parsed commits (with date) between two refs.
func GetParsedCommits(path, from, to string) ([]*Commit, error) {
	const sep = "-_-_-_-_-_END-_-_-_-_-_"
	format := fmt.Sprintf("%s%%n%%cI%%n%%B", sep)

	out, err := exec.Command("git", "-C", path, "log", "--no-merges", "--pretty=format:"+format, from+".."+to).Output()
	if err != nil {
		return nil, err
	}

	return parseCommitsWithDate(string(out), sep), nil
}

func parseCommitsWithDate(raw, sep string) []*Commit {
	parts := strings.Split(raw, sep+"\n")
	var commits []*Commit

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		lines := strings.SplitN(part, "\n", 2)
		if len(lines) < 2 {
			continue
		}
		commits = append(commits, parseCommitWithDate(lines[1], lines[0]))
	}

	return commits
}

// parseCommitWithDate returns a Commit with parsed message and a parsed time.Time.
func parseCommitWithDate(message, dateStr string) *Commit {
	commit := parseCommit(message)
	t, err := time.Parse(time.RFC3339, dateStr)
	if err == nil {
		commit.Date = t
	}
	return commit
}

// parseCommit parses a raw commit message and returns a Commit struct.
func parseCommit(message string) *Commit {
	lines := strings.Split(message, "\n")

	commit := &Commit{}
	if matches := headerRegex.FindStringSubmatch(lines[0]); len(matches) > 0 {
		commit.Type = matches[1]
		commit.Scope = matches[2]
		commit.Description = matches[3]
	} else {
		commit.Description = lines[0]
	}

	if len(lines) > 1 {
		commit.Body = strings.TrimSpace(strings.Join(lines[1:], "\n"))
	}

	return commit
}
