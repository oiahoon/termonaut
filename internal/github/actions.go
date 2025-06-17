package github

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// WorkflowTemplate represents a GitHub Actions workflow template
type WorkflowTemplate struct {
	Name        string
	Description string
	Template    string
}

// ActionsManager handles GitHub Actions integration
type ActionsManager struct {
	repoOwner string
	repoName  string
}

// NewActionsManager creates a new GitHub Actions manager
func NewActionsManager(repoOwner, repoName string) *ActionsManager {
	return &ActionsManager{
		repoOwner: repoOwner,
		repoName:  repoName,
	}
}

// GetWorkflowTemplates returns available workflow templates
func (am *ActionsManager) GetWorkflowTemplates() []WorkflowTemplate {
	return []WorkflowTemplate{
		{
			Name:        "termonaut-stats-update",
			Description: "Automatically update Termonaut badges and stats",
			Template:    termonautStatsWorkflow,
		},
		{
			Name:        "termonaut-profile-sync",
			Description: "Sync Termonaut profile data to repository",
			Template:    termonautProfileSyncWorkflow,
		},
		{
			Name:        "termonaut-weekly-report",
			Description: "Generate weekly productivity reports",
			Template:    termonautWeeklyReportWorkflow,
		},
	}
}

// GenerateWorkflowFile generates a workflow file for a given template
func (am *ActionsManager) GenerateWorkflowFile(templateName string, params map[string]interface{}) (string, error) {
	templates := am.GetWorkflowTemplates()

	var selectedTemplate *WorkflowTemplate
	for _, tmpl := range templates {
		if tmpl.Name == templateName {
			selectedTemplate = &tmpl
			break
		}
	}

	if selectedTemplate == nil {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	// Set default parameters
	if params == nil {
		params = make(map[string]interface{})
	}

	params["RepoOwner"] = am.repoOwner
	params["RepoName"] = am.repoName
	params["Timestamp"] = time.Now().Format(time.RFC3339)

	tmpl, err := template.New(templateName).Parse(selectedTemplate.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, params); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// GetWorkflowFilePath returns the file path for a workflow
func (am *ActionsManager) GetWorkflowFilePath(templateName string) string {
	return filepath.Join(".github", "workflows", fmt.Sprintf("%s.yml", templateName))
}

// Badge endpoint configuration
type BadgeEndpoint struct {
	Path        string
	Label       string
	Description string
}

// GetBadgeEndpoints returns available badge endpoints
func (am *ActionsManager) GetBadgeEndpoints() []BadgeEndpoint {
	return []BadgeEndpoint{
		{
			Path:        "badges/commands.json",
			Label:       "Commands",
			Description: "Total commands executed",
		},
		{
			Path:        "badges/level.json",
			Label:       "Level",
			Description: "Current experience level",
		},
		{
			Path:        "badges/streak.json",
			Label:       "Streak",
			Description: "Current activity streak",
		},
		{
			Path:        "badges/productivity.json",
			Label:       "Productivity",
			Description: "Productivity score",
		},
		{
			Path:        "badges/achievements.json",
			Label:       "Achievements",
			Description: "Achievement progress",
		},
		{
			Path:        "badges/last-active.json",
			Label:       "Last Active",
			Description: "Last activity timestamp",
		},
	}
}

// GenerateBadgeEndpointJSON generates JSON for Shields.io endpoint
func (am *ActionsManager) GenerateBadgeEndpointJSON(label, message, color string) string {
	return fmt.Sprintf(`{
  "schemaVersion": 1,
  "label": "%s",
  "message": "%s",
  "color": "%s",
  "namedLogo": "terminal",
  "logoColor": "white"
}`, label, message, color)
}

// Workflow templates
const termonautStatsWorkflow = `name: Update Termonaut Stats

on:
  schedule:
    # Run every 6 hours
    - cron: '0 */6 * * *'
  workflow_dispatch:
  push:
    paths:
      - '.termonaut/**'

jobs:
  update-stats:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Termonaut
      run: |
        go install github.com/{{.RepoOwner}}/{{.RepoName}}/cmd/termonaut@latest

    - name: Generate Badge Data
      run: |
        mkdir -p badges

        # Generate badges from local data if available
        if [ -f ".termonaut/termonaut.db" ]; then
          termonaut badges generate --output badges/
        else
          # Generate placeholder badges
          echo '{"schemaVersion":1,"label":"Commands","message":"0","color":"lightgrey"}' > badges/commands.json
          echo '{"schemaVersion":1,"label":"Level","message":"1","color":"lightgrey"}' > badges/level.json
          echo '{"schemaVersion":1,"label":"Streak","message":"0 days","color":"red"}' > badges/streak.json
          echo '{"schemaVersion":1,"label":"Productivity","message":"0.0%","color":"red"}' > badges/productivity.json
          echo '{"schemaVersion":1,"label":"Achievements","message":"0/0","color":"lightgrey"}' > badges/achievements.json
        fi

    - name: Update Profile Stats
      run: |
        # Generate profile stats markdown
        if [ -f ".termonaut/termonaut.db" ]; then
          termonaut profile generate --output TERMONAUT_PROFILE.md
        fi

    - name: Commit and push changes
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add badges/ TERMONAUT_PROFILE.md || true
        git diff --staged --quiet || git commit -m "🚀 Update Termonaut stats - $(date)"
        git push
`

const termonautProfileSyncWorkflow = `name: Sync Termonaut Profile

on:
  workflow_dispatch:
    inputs:
      sync_data:
        description: 'Sync profile data'
        required: false
        default: 'true'

jobs:
  sync-profile:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Termonaut
      run: |
        go install github.com/{{.RepoOwner}}/{{.RepoName}}/cmd/termonaut@latest

    - name: Sync Profile Data
      run: |
        # This would sync profile data from a secure source
        # Implementation depends on how users want to sync their data
        echo "Profile sync functionality - coming soon!"

    - name: Generate Profile
      run: |
        mkdir -p profile

        # Generate comprehensive profile
        if [ -f ".termonaut/termonaut.db" ]; then
          termonaut profile generate --format=markdown --output=profile/README.md
          termonaut analytics export --output=profile/analytics.json
          termonaut achievements list --format=json --output=profile/achievements.json
        fi

    - name: Commit changes
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add profile/ || true
        git diff --staged --quiet || git commit -m "📊 Sync Termonaut profile - $(date)"
        git push
`

const termonautWeeklyReportWorkflow = `name: Weekly Termonaut Report

on:
  schedule:
    # Run every Monday at 9 AM UTC
    - cron: '0 9 * * 1'
  workflow_dispatch:

jobs:
  weekly-report:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Termonaut
      run: |
        go install github.com/{{.RepoOwner}}/{{.RepoName}}/cmd/termonaut@latest

    - name: Generate Weekly Report
      run: |
        mkdir -p reports

        if [ -f ".termonaut/termonaut.db" ]; then
          # Generate weekly report
          WEEK=$(date +'%Y-W%U')
          termonaut analytics weekly --output=reports/week-$WEEK.md

          # Update latest report
          cp reports/week-$WEEK.md reports/latest-weekly.md
        fi

    - name: Create Issue with Report
      uses: actions/github-script@v7
      with:
        script: |
          const fs = require('fs');

          try {
            const report = fs.readFileSync('reports/latest-weekly.md', 'utf8');

            await github.rest.issues.create({
              owner: context.repo.owner,
              repo: context.repo.repo,
              title: '📊 Weekly Termonaut Report - Week ' + new Date().toISOString().slice(0, 10),
              body: report,
              labels: ['termonaut', 'weekly-report']
            });
          } catch (error) {
            console.log('No report generated or error creating issue:', error);
          }

    - name: Commit reports
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add reports/ || true
        git diff --staged --quiet || git commit -m "📈 Weekly Termonaut report - $(date +'%Y-W%U')"
        git push
`
