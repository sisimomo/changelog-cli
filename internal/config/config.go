package config

import (
	"strings"

	"github.com/urfave/cli/v2"
)

// TypeMappingEntry represents a single mapping from a commit type to a display title.
type TypeMappingEntry struct {
	Key   string
	Title string
}

// TypeMappings is an array of TypeMappingEntry.
type TypeMappings []TypeMappingEntry

// Config holds all the configuration for the changelog generation.
type Config struct {
	RepoPath          string
	From              string
	To                string
	TypeMap           TypeMappings
	FallbackTypeTitle string
	TicketPattern     string
	TicketTemplateURL string
	OutputFile        string
}

var defaultTypeMapping = TypeMappings{
	{"feat", "New Features"},
	{"fix", "Fixes"},
	{"perf", "Performance Improvements"},
	{"refactor", "Code Refactoring"},
	{"style", "Code Style"},
	{"test", "Add or Update Tests"},
	{"docs", "Documentation"},
	{"build", "Build System"},
	{"ci", "Continuous Integration"},
}

// GetConfig creates a new Config object from the cli context.
func GetConfig(c *cli.Context) *Config {
	return &Config{
		RepoPath:          c.String("repo"),
		From:              c.String("from"),
		To:                c.String("to"),
		TypeMap:           prepareTypeMappings(c.String("type-map")),
		FallbackTypeTitle: c.String("fallback-type-title"),
		TicketPattern:     c.String("ticket-pattern"),
		TicketTemplateURL: c.String("ticket-template-url"),
		OutputFile:        c.String("output"),
	}
}

func prepareTypeMappings(cliTypeMap string) TypeMappings {
	if cliTypeMap == "" {
		return defaultTypeMapping
	}
	var typeMapping TypeMappings
	pairs := strings.Split(cliTypeMap, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			typeMapping = append(typeMapping, TypeMappingEntry{
				Key:   strings.TrimSpace(parts[0]),
				Title: strings.TrimSpace(parts[1]),
			})
		}
	}
	return typeMapping
}

// CleanRepoPath removes leading/trailing quotes from the repository path.
func (c *Config) CleanRepoPath() {
	c.RepoPath = strings.Trim(c.RepoPath, "\"")
}
