package cliutil

import (
	"changelog-cli/internal/config"
	"fmt"
	"strings"
)

// PrintChangelogParameters prints the initial parameters before changelog generation.
func PrintChangelogParameters(cfg *config.Config) {
	fmt.Println("Changelog Generation Parameters:")
	fmt.Printf("  Repository: %s\n", cfg.RepoPath)
	fmt.Printf("  From: %s\n", cfg.From)
	fmt.Printf("  To: %s\n", cfg.To)
	fmt.Printf("  Ticket Pattern: %s\n", cfg.TicketPattern)
	fmt.Printf("  Ticket Template URL: %s\n", cfg.TicketTemplateURL)

	// Print type mappings
	var sb strings.Builder
	sb.WriteString("  Type Mappings:\n")
	for _, entry := range cfg.TypeMap {
		sb.WriteString(fmt.Sprintf("    - %s: %s\n", entry.Key, entry.Title))
	}
	fmt.Print(sb.String())
	fmt.Printf("  Fallback Type Title: %s\n", cfg.FallbackTypeTitle)
	fmt.Printf("--------------------------------------------------------------\n")
}
