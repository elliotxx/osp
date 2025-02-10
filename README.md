![OSP](https://socialify.git.ci/elliotxx/osp/image?font=Raleway&language=1&name=1&owner=1&pattern=Plus&theme=Light)

# OSP - Open Source Pilot for Community Governance

[![Go Report Card](https://goreportcard.com/badge/github.com/elliotxx/osp)](https://goreportcard.com/report/github.com/elliotxx/osp)
[![GoDoc](https://godoc.org/github.com/elliotxx/osp?status.svg)](https://godoc.org/github.com/elliotxx/osp)
[![License](https://img.shields.io/github/license/elliotxx/osp.svg)](https://github.com/elliotxx/osp/blob/main/LICENSE)

[English](README.md) | [简体中文](README_zh.md)

OSP (Open Source Pilot) is an automated tool focused on open source community governance. It integrates various best practices in open source community governance, providing maintainers with a comprehensive toolkit for efficient operations, precise tracking, and data-driven decision making.

OSP offers two usage modes: a CLI tool and GitHub Action workflows. The CLI tool is suitable for local interactive management, while GitHub Action enables fully automated operations through github event subscriptions - configure once, serve continuously.

Actual Example:
- [Automated updated community task list](https://github.com/KusionStack/karpor/issues/463)
- [Automated updated project planning](https://github.com/KusionStack/karpor/issues/723)

## ✨ Features

### Implemented
- 🔑 GitHub Authentication - Secure identity authentication, same as GitHub CLI
- 📊 Project Statistics - Multi-dimensional data analysis
- 📝 Community Tasks & Project Planning - Auto-updates through GitHub event subscriptions
- 📈 Star History - Project growth tracking

### Roadmap
- 📋 Roadmap Generation - Auto-updates through GitHub event subscriptions
- 📅 Community Activity Aggregation - Auto-aggregates recent comments, new PRs/Issues/Discussions
- 🤖 Smart PR Review - LLM-based code review with automated comments
- 💡 Smart Issue Creation - One-line issue generation for improved efficiency
- 🔌 GitHub App Integration - Enhanced integration capabilities
- 📝 Release Note Generation - Auto-summarizes core changes, contributors, and community participation metrics
- Support reading `<!-- CUSTOM -->` tag content from issue will not be overwritten
- Support display `diff` content before plan and onboard excute
- Add descriptions for each label in `plan` and `onboard` templates
- Support recent activity (latest finish issue, etc) in `plan` and `onboard` templates, such as showing recently closed issues
- Add explanation for difficulty symbol `!` in `osp plan` template

## 🚀 Installation

For more installation options, see [Advanced Installation Guide](docs/guide/advanced-installation.md).

### 🐙 Via Go

```bash
go install github.com/elliotxx/osp@latest
```

### 🍺 Via Homebrew

```bash
brew tap elliotxx/tap
brew install osp
```

## 🚀 Usage

### 🖥️ CLI Tool

1. Configure GitHub Authentication
```bash
# Login with GitHub CLI
gh auth login

# Verify authentication status
osp auth status
```

2. Project Management
```bash
# Add a project
osp repo add owner/repo

# Switch projects
osp repo

# View current project
osp repo current
```

3. Core Features
```bash
# Generate version planning
osp plan

# Manage community tasks
osp onboard

# View project statistics
osp stats

# Analyze Star history
osp star history
```

For more details, see the [CLI Usage Guide](docs/guide/cli.md).

### 🤖 GitHub Action

> For osp-action implementation, see [osp-action](https://github.com/elliotxx/osp-action)

Here's an example of automating community task generation and updates using osp-action. For more automation scenarios, see the [documentation](docs/guide/github-action.md).

1. Create workflow file `.github/workflows/community-task-updater.yml` in main branch:
```yaml
name: Community Task Updater

on:
  # Manually triggered
  workflow_dispatch:
  # Trigger on issue events
  issues:
    types: [opened, edited, deleted, transferred, milestoned, demilestoned, labeled, unlabeled, assigned, unassigned]

jobs:
  osp-run:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Update Community Tasks
        uses: elliotxx/osp-action@main
        with:
          # Optional: version of OSP to use (default: latest)
          version: 'latest'
          
          # Optional: working directory (default: project root)
          working-directory: '.'
          
          # Optional: GitHub token (default: ${{ github.token }})
          github-token: ${{ secrets.GITHUB_TOKEN }}
          
          # Optional: enable debug mode (default: false)
          debug: false
          
          # Optional: skip caching (default: false)
          skip-cache: false
          
          # Optional: additional OSP arguments
          args: >-
            onboard
            --yes
            --onboard-labels 'help wanted,good first issue'
            --difficulty-labels 'good first issue,help wanted'
            --category-labels bug,documentation,enhancement
            --target-title 'Community Tasks 🎯'
```

2. Configure Required Permissions
- Navigate to Settings -> Actions -> General
- Enable "Read and write permissions" under "Workflow permissions"
- Save the changes

3. Usage
- Automatic: Workflow executes automatically when configured GitHub events occur
- Manual:
  1. Go to Actions page
  2. Select "Community Task Updater"
  3. Click "Run workflow"

## 📚 Documentation

- [User Guide](docs/guide/README.md) - Detailed usage instructions
- [Design Docs](docs/design/README.md) - Architecture and implementation
- [CLI Reference](docs/cli/osp.md) - Command-line tool reference

## 👥 Who's using it

- [karpor](https://github.com/KusionStack/karpor)

## 🤝 Contributing

We welcome all forms of contributions! Whether it's new features, documentation improvements, or bug fixes. See our [Contributing Guide](CONTRIBUTING.md) for details.

## 👀 Similar Projects

- [Oscar](https://github.com/golang/oscar)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
