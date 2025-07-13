package changelog

import (
	"changelog-cli/internal/config"
	"changelog-cli/internal/git"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func formatCommit(commit *git.Commit) string {
	if commit.Body == "" {
		return "* " + commit.Description
	}
	return fmt.Sprintf("* <details>\n     <summary>%s</summary>\n%s   </details>",
		commit.Description,
		renderCommitBody(commit.Body),
	)
}

func renderCommitBody(body string) string {
	var formatted strings.Builder
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		formatted.WriteString("       ")
		if strings.TrimSpace(line) != "" {
			formatted.WriteString(line)
		} else {
			formatted.WriteString("<br><br>")
		}
		formatted.WriteString("\n")
	}
	return formatted.String()
}

func sortCommitsByDate(commits []*git.Commit) {
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.Before(commits[j].Date)
	})
}

func groupCommitsByType(commits []*git.Commit, validTypes map[string]string) (map[string][]*git.Commit, []*git.Commit) {
	grouped := make(map[string][]*git.Commit)
	var fallback []*git.Commit
	for _, commit := range commits {
		if _, ok := validTypes[commit.Type]; ok {
			grouped[commit.Type] = append(grouped[commit.Type], commit)
		} else {
			fallback = append(fallback, commit)
		}
	}
	return grouped, fallback
}

func appendCommitGroupSection(builder *strings.Builder, title string, commits []*git.Commit) {
	if len(commits) == 0 {
		return
	}
	builder.WriteString(fmt.Sprintf("\n## %s\n", title))
	for _, commit := range commits {
		builder.WriteString(fmt.Sprintf("%s\n", formatCommit(commit)))
	}
}

func replaceTickets(markdown, pattern, urlTemplate string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(markdown, fmt.Sprintf("<a href=\"%s\">%s</a>", strings.Replace(urlTemplate, "{ticket}", "$0", -1), "$0"))
}

// GenerateMarkdown generates the changelog markdown from a list of commits and a config.
func GenerateMarkdown(commits []*git.Commit, cfg *config.Config) string {
	sortCommitsByDate(commits)

	validTypes := make(map[string]string)
	for _, entry := range cfg.TypeMap {
		validTypes[entry.Key] = entry.Title
	}

	groupedCommits, fallbackCommits := groupCommitsByType(commits, validTypes)

	var builder strings.Builder
	for _, entry := range cfg.TypeMap {
		appendCommitGroupSection(&builder, entry.Title, groupedCommits[entry.Key])
	}
	appendCommitGroupSection(&builder, cfg.FallbackTypeTitle, fallbackCommits)

	if cfg.TicketPattern != "" {
		result := builder.String()
		return replaceTickets(result, cfg.TicketPattern, cfg.TicketTemplateURL)
	}

	return builder.String()
}
