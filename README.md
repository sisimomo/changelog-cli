# changelog-cli

A command-line tool to automatically generate, structured changelogs from your Git commit history. The repo must follows Conventional Commit specifications so that the tool can group changes meaningfully.

The tool is designed to be simple to use, yet flexible enough to integrate into various workflows. It can automatically detect commit ranges based on your Git tags, or you can specify them manually.

## Features

-   Parses conventional commit messages.
-   Groups commits by type (e.g., Features, Bug Fixes, etc.).
-   Automatically determines commit ranges based on Git tags.
-   Supports custom commit type mappings.
-   Links ticket numbers (e.g., JIRA, GitHub issues) automatically.
-   Outputs to the console or a file.
-   Supports expandable commit bodies in the markdown output for a clean and detailed view.

## Intent

The goal of this project is to streamline the changelog creation process. By leveraging the structure of conventional commits, it automates the tedious task of compiling release notes, ensuring consistency and accuracy. This allows developers to focus on building software while maintaining a clear and professional record of changes.

## Usage

The primary command is `generate`, which orchestrates the entire changelog creation process.

```bash
changelog-cli generate [flags]
```

### Default Behavior

When run without `--from` or `--to` flags, the tool automatically determines the commit range:
-   It uses the latest Git tag as the `to` reference.
-   It uses the tag immediately preceding the latest tag as the `from` reference.
-   If only one tag exists in the repository, it uses the very first commit as the `from` reference.

## Commands and Flags

### `generate`

Generates a changelog from Git commit messages.

| Flag                    | Description                                                                                                                                                        | Default Value |
| ----------------------- |--------------------------------------------------------------------------------------------------------------------------------------------------------------------| ------------- |
| `--from`                | The starting commit reference (e.g., a tag or commit hash). If omitted, it's determined automatically from the previous tag. If provided, `--to` must also be set. | (auto)        |
| `--to`                  | The ending commit reference (e.g., a tag or commit hash). If omitted, it defaults to the latest Git tag. If provided, `--from` must also be set.                   | (auto)        |
| `--repo`                | Path to the Git repository.                                                                                                                                        | `.`           |
| `--output`              | Path to the output file. If the file already exists, the command will fail to prevent accidental overwrites. If omitted, the changelog is printed to the console.  | (console)     |
| `--type-map`            | A comma-separated list of `type=Title` pairs to define which commit types are included and how their sections are titled (e.g., `'feat=Features,fix=Bug Fixes'`).  | (see below)   |
| `--fallback-type-title` | The title for the section containing commits of types not defined in the `type-map`.                                                                               | `Other`       |
| `--ticket-pattern`      | A regex pattern to extract ticket numbers from commit messages (e.g., `DKQ-\d+`).                                                                                  | (none)        |
| `--ticket-template-url` | A URL template for generating ticket links. Use `{ticket}` as a placeholder for the extracted ticket number (e.g., `https://my-jira.com/browse/{ticket}`).         | (none)        |

### Default Type Mappings

If the `--type-map` flag is not provided, the following default mappings are used:

-   `feat`: **New Features**
-   `fix`: **Fixes**
-   `perf`: **Performance Improvements**
-   `refactor`: **Code Refactoring**
-   `style`: **Code Style**
-   `test`: **Add or Update Tests**
-   `docs`: **Documentation**
-   `build`: **Build System**
-   `ci`: **Continuous Integration**

## Examples

### 1. Automatic Changelog

Generate a changelog for the commits between the last two tags and print it to the console.

```bash
changelog-cli generate
```

### 2. Specify a Version Range

Generate a changelog for a specific version range.

```bash
changelog-cli generate --from v1.0.0 --to v1.1.0
```

### 3. Write to a File

Generate a changelog and save it to a `CHANGELOG.md` file.

```bash
changelog-cli generate --from v1.0.0 --to v1.1.0 --output CHANGELOG.md
```

### 4. Use Custom Type Mappings

Generate a changelog with a custom set of types and titles, overriding the default ones.

```bash
changelog-cli generate --type-map "feat=✨ New Features,fix=🐛 Bug Fixes,perf=🚀 Performance"
```

### 5. Link to JIRA Tickets

Generate a changelog and automatically link ticket numbers found in commit messages.

```bash
changelog-cli generate \
  --ticket-pattern "PROJ-\d+" \
  --ticket-template-url "https://my-jira.com/browse/{ticket}"
```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
