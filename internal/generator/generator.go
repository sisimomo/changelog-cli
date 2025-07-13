package generator

import (
	"changelog-cli/internal/changelog"
	"changelog-cli/internal/config"
	"changelog-cli/internal/fileutil"
	"changelog-cli/internal/git"
	"fmt"
)

// Generate orchestrates the changelog generation process.
func Generate(cfg *config.Config) error {
	commits, err := git.GetParsedCommits(cfg.RepoPath, cfg.From, cfg.To)
	if err != nil {
		return err
	}

	changelogContent := changelog.GenerateMarkdown(commits, cfg)

	if cfg.OutputFile != "" {
		return fileutil.WriteChangelogToFile(cfg.OutputFile, changelogContent)
	} else {
		fmt.Print(changelogContent)
	}
	return nil
}

// DetermineRefs resolves from/to references based on CLI input or Git tags.
func DetermineRefs(repo, from, to string) (string, string, error) {
	if from == "" && to == "" {
		to, prev, err := git.GetLastTwoTags(repo)
		if err != nil {
			return "", "", fmt.Errorf("failed to get last two tags: %w", err)
		}
		if prev != nil {
			from = *prev
		} else {
			// If there is only one tag, use the initial commit as the starting point.
			initialCommit, err := git.GetInitialCommit(repo)
			if err != nil {
				return "", "", err
			}
			from = initialCommit
		}
		return from, to, nil
	}

	if (from == "") != (to == "") {
		return "", "", fmt.Errorf("if either 'from' or 'to' is provided, both must be provided")
	}

	return from, to, nil
}
