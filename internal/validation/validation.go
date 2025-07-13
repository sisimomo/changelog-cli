package validation

import (
	"changelog-cli/internal/config"
	"changelog-cli/internal/git"
	"fmt"
	"os"
	"path/filepath"
)

// PrepareAndValidate validates the configuration.
func PrepareAndValidate(cfg *config.Config) error {
	cfg.RepoPath = filepath.FromSlash(cfg.RepoPath)

	// Validate repository path existence
	if err := ValidateRepositoryPath(cfg.RepoPath); err != nil {
		return err
	}

	// Validate if it's a Git repository
	if err := ValidateGitRepository(cfg.RepoPath); err != nil {
		return err
	}

	// Validate 'from' tag existence if provided
	if err := ValidateTagExistence(cfg.RepoPath, cfg.From); err != nil {
		return err
	}

	if err := ValidateTagExistence(cfg.RepoPath, cfg.To); err != nil {
		return err
	}

	return nil
}

// ValidateRepositoryPath checks if the given path exists and is a directory.
func ValidateRepositoryPath(repositoryPath string) error {
	information, err := os.Stat(repositoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("repository path does not exist: %s", repositoryPath)
		}
		return fmt.Errorf("error accessing repository path %q: %w", repositoryPath, err)
	}
	if !information.IsDir() {
		return fmt.Errorf("repository path is not a directory: %s", repositoryPath)
	}
	return nil
}

// ValidateGitRepository checks if the given path is a valid Git repository.
func ValidateGitRepository(repositoryPath string) error {
	return git.IsGitRepository(repositoryPath)
}

// ValidateTagExistence checks if a given tag exists in the repository.
func ValidateTagExistence(repositoryPath, tag string) error {
	if tag == "" {
		return nil // No tag to validate
	}

	exists, err := git.TagExists(repositoryPath, tag)
	if err != nil {
		return fmt.Errorf("error checking tag existence for %q: %w", tag, err)
	}
	if !exists {
		return fmt.Errorf("tag does not exist in the repository: %s", tag)
	}
	return nil
}
