package onboard

import (
	"strings"
	"testing"

	"github.com/elliotxx/osp/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateProgressBar(t *testing.T) {
	tests := []struct {
		name      string
		completed int
		total     int
		want      string
	}{
		{
			name:      "empty progress",
			completed: 0,
			total:     10,
			want:      "░░░░░░░░░░░░░░░░░░░░ 0.0%",
		},
		{
			name:      "half progress",
			completed: 5,
			total:     10,
			want:      "██████████░░░░░░░░░░ 50.0%",
		},
		{
			name:      "full progress",
			completed: 10,
			total:     10,
			want:      "████████████████████ 100.0%",
		},
		{
			name:      "no tasks",
			completed: 0,
			total:     0,
			want:      "░░░░░░░░░░░░░░░░░░░░",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateProgressBar(tt.completed, tt.total)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateContent(t *testing.T) {
	// Mock issues
	issues := []OnboardIssue{
		{
			Difficulty: "difficulty/easy",
			Status:     "closed",
			Assignee:   "user1",
			Number:     1,
			Category:   "bug",
		},
		{
			Difficulty: "difficulty/medium",
			Status:     "open",
			Number:     2,
			Category:   "enhancement",
		},
		{
			Difficulty: "difficulty/hard",
			Status:     "open",
			Number:     3,
			Category:   "documentation",
		},
	}

	// Mock options
	opts := Options{
		// Issue labels configuration
		OnboardLabels:    []string{"help wanted", "good first issue"},
		DifficultyLabels: []string{"difficulty/easy", "difficulty/medium", "difficulty/hard"},
		CategoryLabels:   []string{"bug", "enhancement", "documentation"},

		// Target issue configuration
		TargetLabel: "onboarding",
		TargetTitle: "Community Tasks",

		// Command behavior
		DryRun:      false,
		AutoConfirm: true,
	}

	// Create manager
	cfg := &config.Config{
		Current: "elliotxx/osp",
		Auth: struct {
			Token string `yaml:"token"`
		}{
			Token: "test-token",
		},
	}
	m := NewManager(cfg, nil)

	// Generate content
	content, err := m.GenerateContent(issues, "elliotxx/osp", opts)
	assert.NoError(t, err)

	// Verify content sections
	sections := []struct {
		name    string
		expects []string
	}{
		{
			name: "header",
			expects: []string{
				"## Overview ",
				"- Progress: ██████░░░░░░░░░░░░░░ 33.3%",
				"- [Total Issues: 3](https://github.com/elliotxx/osp/issues?q=is%3Aissue+(label%3A%22help+wanted%22+OR+label%3A%22good+first+issue%22))",
				"  - ✅ [Completed: 1](https://github.com/elliotxx/osp/issues?q=is%3Aissue+is%3Aclosed+(label%3A%22help+wanted%22+OR+label%3A%22good+first+issue%22))",
				"  - 🚧 [In Progress: 0](https://github.com/elliotxx/osp/issues?q=is%3Aissue+is%3Aopen+assignee%3A*+(label%3A%22help+wanted%22+OR+label%3A%22good+first+issue%22))",
				"  - 📋 [Unassigned: 2](https://github.com/elliotxx/osp/issues?q=is%3Aissue+(label%3A%22help+wanted%22+OR+label%3A%22good+first+issue%22)+no%3Aassignee)",
				"## Description",
				"As a programming enthusiast, have you ever felt that you want to participate in the development of an open source project, but don't know where to start?",
				"In order to help everyone better participate in open source projects, we regularly publish issues suitable for new contributors to help everyone learn by doing!",
			},
		},
		{
			name: "contributors",
			expects: []string{
				"## Contributors (1)",
				"Thanks to all our contributors who have completed onboarding issues! Your contributions help make our project better:",
				"@user1 ",
			},
		},
		{
			name: "issues",
			expects: []string{
				"## Issue List (3)",
				"> The following onboarding issues are organized first by difficulty level (from easy to hard), and then by category within each difficulty level.",
				"### Difficulty: difficulty/easy (1)",
				"#### Category: bug (1)",
				"- [x] #1 **[@user1 did it! Cheers! 🍻]**",
				"### Difficulty: difficulty/medium (1)",
				"#### Category: enhancement (1)",
				"- [ ] #2",
				"### Difficulty: difficulty/hard (1)",
				"#### Category: documentation (1)",
				"- [ ] #3",
			},
		},
	}

	// Helper function to filter empty lines
	filterEmptyLines := func(lines []string) []string {
		var result []string
		for _, line := range lines {
			if line != "" {
				result = append(result, line)
			}
		}
		return result
	}

	// Split content into lines and filter empty lines
	contentLines := filterEmptyLines(strings.Split(strings.TrimSpace(content), "\n"))

	// Verify each section
	currentLine := 0
	for _, section := range sections {
		t.Run(section.name, func(t *testing.T) {
			expectedLines := filterEmptyLines(section.expects)
			for i, expectedLine := range expectedLines {
				if currentLine+i >= len(contentLines) {
					t.Errorf("content ended too early, expected line %q at position %d", expectedLine, currentLine+i)
					return
				}
				assert.Equal(t, expectedLine, contentLines[currentLine+i], "content mismatch at line %d", currentLine+i)
			}
			currentLine += len(expectedLines)
		})
	}
}
