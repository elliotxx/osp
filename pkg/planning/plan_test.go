package planning

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateProgressBar(t *testing.T) {
	tests := []struct {
		name      string
		completed int
		total     int
		length    int
		want      string
	}{
		{
			name:      "empty progress",
			completed: 0,
			total:     0,
			length:    10,
			want:      "░░░░░░░░░░ 0%",
		},
		{
			name:      "half progress",
			completed: 5,
			total:     10,
			length:    10,
			want:      "█████░░░░░ 50%",
		},
		{
			name:      "full progress",
			completed: 10,
			total:     10,
			length:    10,
			want:      "██████████ 100%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateProgressBar(tt.completed, tt.total, tt.length)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrepareTemplateData(t *testing.T) {
	// Mock milestone
	milestone := Milestone{
		Title:       "Test Milestone",
		Description: "Test Description",
		DueOn:       &time.Time{},
		Number:      1,
		State:       "open",
		HTMLURL:     "https://github.com/elliotxx/osp/milestone/1",
	}

	// Mock issues
	issues := []Issue{
		{
			Title:  "Issue 1",
			Number: 1,
			State:  "closed",
			Labels: []Label{
				{Name: "bug"},
				{Name: "priority/high"},
			},
			Assignee: &User{Login: "user1"},
		},
		{
			Title:  "Issue 2",
			Number: 2,
			State:  "open",
			Labels: []Label{
				{Name: "enhancement"},
				{Name: "good first issue"},
				{Name: "priority/medium"},
			},
			Assignee: &User{Login: "user2"},
		},
		{
			Title:  "Issue 3",
			Number: 3,
			State:  "open",
			Labels: []Label{
				{Name: "question"},
				{Name: "priority/low"},
			},
		},
	}

	// Test categories
	categories := []string{"bug", "enhancement"}

	// Create manager and options
	m := &Manager{}
	opts := Options{
		Categories: categories,
		Priorities: []string{"priority/high", "priority/medium", "priority/low"},
	}

	// Prepare template data
	data := m.prepareTemplateData(milestone, issues, opts)

	// Assertions
	t.Run("milestone data", func(t *testing.T) {
		assert.Equal(t, milestone.Title, data.Milestone.Title)
		assert.Equal(t, milestone.Description, data.Milestone.Description)
		assert.Equal(t, milestone.Number, data.Milestone.Number)
	})

	t.Run("statistics", func(t *testing.T) {
		assert.Equal(t, 3, data.Stats.TotalIssues)
		assert.Equal(t, 1, data.Stats.CompletedIssues)
		assert.InDelta(t, 33.33, data.Stats.Progress, 0.01)
		assert.ElementsMatch(t, []string{"user1"}, data.Stats.Contributors)
	})

	t.Run("categorized issues", func(t *testing.T) {
		// Check bug category
		bugIssues := data.Issues["bug"]
		assert.Len(t, bugIssues, 1)
		assert.Equal(t, 1, bugIssues[0].Number)

		// Check enhancement category
		enhancementIssues := data.Issues["enhancement"]
		assert.Len(t, enhancementIssues, 1)
		assert.Equal(t, 2, enhancementIssues[0].Number)
	})

	t.Run("uncategorized issues", func(t *testing.T) {
		assert.Len(t, data.UncategorizedIssues, 1)
		assert.Equal(t, 3, data.UncategorizedIssues[0].Number)
	})

	t.Run("progress bar", func(t *testing.T) {
		assert.Contains(t, data.ProgressBar, "33%")
	})
}

func TestFindPlanningIssue(t *testing.T) {
	// Mock existing issues with duplicate planning issues
	existingIssues := []Issue{
		{
			Title:  "Planning: v1.0.0",
			Number: 5,
			State:  "open",
		},
		{
			Title:  "Planning: v1.0.0",
			Number: 3,
			State:  "open",
		},
		{
			Title:  "Planning: v1.0.0",
			Number: 8,
			State:  "open",
		},
		{
			Title:  "Other Issue",
			Number: 1,
			State:  "open",
		},
	}

	planningTitle := "Planning: v1.0.0"
	var planningIssue *Issue
	minIssueNumber := 2147483647
	for _, issue := range existingIssues {
		if issue.Title == planningTitle {
			if planningIssue == nil || issue.Number < minIssueNumber {
				planningIssue = &issue
				minIssueNumber = issue.Number
			}
		}
	}

	// Assertions
	assert.NotNil(t, planningIssue)
	assert.Equal(t, 3, planningIssue.Number, "Should select the planning issue with the smallest number")
}

func TestGeneratePlanningContent(t *testing.T) {
	// Create a fixed time for testing
	fixedTime := time.Date(2025, 1, 30, 15, 0o4, 0o5, 0, time.UTC)
	dueDate := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)

	// Mock template data
	data := TemplateData{
		Milestone: Milestone{
			Title:       "v1.0.0",
			Description: "First stable release",
			DueOn:       &dueDate,
			Number:      1,
			State:       "open",
			HTMLURL:     "https://github.com/elliotxx/osp/milestone/1",
		},
		Stats: MilestoneStats{
			TotalIssues:     3,
			CompletedIssues: 1,
			Progress:        33.33,
			Contributors:    []string{"user1"},
		},
		Categories: []string{"bug", "enhancement", "documentation"},
		Priorities: []string{"priority/high", "priority/medium", "priority/low"},
		Issues: map[string][]Issue{
			"bug": {
				{
					Title:  "Critical Bug",
					Number: 1,
					State:  "closed",
					Labels: []Label{
						{Name: "bug"},
						{Name: "priority/high"},
					},
					Assignee: &User{
						Login: "user1",
					},
				},
			},
			"enhancement": {
				{
					Title:  "New Feature",
					Number: 2,
					State:  "open",
					Labels: []Label{
						{Name: "enhancement"},
						{Name: "good first issue"},
						{Name: "priority/medium"},
					},
					Assignee: &User{
						Login: "user2",
					},
				},
			},
		},
		UncategorizedIssues: []Issue{
			{
				Title:  "Question",
				Number: 3,
				State:  "open",
				Labels: []Label{
					{Name: "question"},
					{Name: "priority/low"},
				},
			},
		},
		ProgressBar: "█████░░░░░░░░░ 33%",
	}

	// Create manager
	m := &Manager{}

	// Generate content
	content, err := m.generatePlanningContentWithTime(data, fixedTime)

	// Assertions
	t.Run("content generation", func(t *testing.T) {
		assert.NoError(t, err)

		// Check each section
		sections := []struct {
			name     string
			expected []string
		}{
			{"header", []string{
				"## Overview",
				"- Progress: █████░░░░░░░░░ 33%",
				"- Total Issues: 3",
				"  - ✅ Completed: 1",
				"  - 🚧 In Progress: 2",
				"- Due Date: February 28, 2025",
				"- Data comes from [Milestone #1](https://github.com/elliotxx/osp/milestone/1)",
				"## Description",
				"First stable release",
				"## Tasks by Category",
			}},
			{"bug section", []string{
				"### bug (1)",
				"- [x] !!! #1 (@user1) `bug` `priority/high`",
			}},
			{"enhancement section", []string{
				"### enhancement (1)",
				"- [ ] !! #2 (@user2) `enhancement` `good first issue` `priority/medium`",
			}},
			{"uncategorized section", []string{
				"### Uncategorized (1)",
				"- [ ] ! #3 `question` `priority/low`",
			}},
			{"contributors section", []string{
				"## Contributors",
				"Thanks to all our contributors for their efforts on completed issues:",
				"- @user1",
			}},
			{"footer", []string{
				"## Links",
				"- 📋 [Issues without priority](https://github.com///issues?q=is%3Aopen+is%3Aissue+milestone%3Av1.0.0+-label%3Apriority/high+-label%3Apriority/medium+-label%3Apriority/low)",
				"- 👥 [Unassigned issues](https://github.com///issues?q=is%3Aopen+is%3Aissue+milestone%3Av1.0.0+no%3Aassignee)",
				"- 📊 [All milestone issues](https://github.com///milestone/1)",
				"---",
				"> 🤖 Auto-generated by [OSP](https://github.com/elliotxx/osp). DO NOT EDIT.",
				"> Last Updated: January 30, 2025 15:04 UTC",
			}},
		}

		// Split content into lines and filter out empty lines
		lines := []string{}
		for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
			if strings.TrimSpace(line) != "" {
				lines = append(lines, line)
			}
		}

		// Check each section
		currentLine := 0
		for _, section := range sections {
			for i, expectedLine := range section.expected {
				if currentLine >= len(lines) {
					t.Errorf("Section %s: Missing line %d: expected '%s'", section.name, i+1, expectedLine)
					continue
				}
				assert.Equal(t, expectedLine, lines[currentLine], "Section %s: Line %d should match", section.name, i+1)
				currentLine++
			}
		}
	})
}
