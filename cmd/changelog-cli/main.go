package main

import (
	"changelog-cli/internal/cliutil"
	"changelog-cli/internal/config"
	"changelog-cli/internal/generator"
	"changelog-cli/internal/validation"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "changelog-cli",
		Usage: "Generate a changelog from Git commit messages",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate a changelog from Git commit messages",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "from",
						Usage: "The starting commit reference (e.g., a tag or commit hash). If omitted (and 'to' is also omitted), the tool automatically determines it as the tag preceding the latest tag. If only one tag exists, it covers commits from the repository's beginning. If provided, 'to' must also be provided.",
					},
					&cli.StringFlag{
						Name:  "to",
						Usage: "The ending commit reference (e.g., a tag or commit hash). If both 'from' and 'to' are omitted, the tool automatically uses the latest Git tag as the 'to' reference. If provided, 'from' must also be provided.",
					},
					&cli.StringFlag{
						Name:  "repo",
						Value: ".",
						Usage: "Path to the Git repository.",
					},
					&cli.StringFlag{
						Name:  "type-map",
						Usage: "Comma-separated list of commit type to display title mappings (e.g., 'feat=Features,fix=Bug Fixes'). If provided, only these mappings will be used.",
					},
					&cli.StringFlag{
						Name:  "fallback-type-title",
						Value: "Other",
						Usage: "The title for the section containing commits of unknown types.",
					},
					&cli.StringFlag{
						Name:  "ticket-pattern",
						Usage: "Regex pattern to extract ticket numbers from commit messages.",
					},
					&cli.StringFlag{
						Name:  "ticket-template-url",
						Usage: "URL template for generating ticket links. Use {ticket} as a placeholder.",
					},
					&cli.StringFlag{
						Name:  "output",
						Usage: "Path to the output file. If the file already exists, the command will fail.",
					},
				},
				Action: runChangelogAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runChangelogAction contains the core logic for generating the changelog.
func runChangelogAction(c *cli.Context) error {
	cliConfig := config.GetConfig(c)
	cliConfig.CleanRepoPath()
	if err := validation.PrepareAndValidate(cliConfig); err != nil {
		return err
	}

	from, to, err := generator.DetermineRefs(cliConfig.RepoPath, cliConfig.From, cliConfig.To)
	if err != nil {
		return err
	}
	cliConfig.From = from
	cliConfig.To = to

	cliutil.PrintChangelogParameters(cliConfig)

	return generator.Generate(cliConfig)
}
