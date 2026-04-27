package main

import (
	"fmt"
	"log"
	"strings"

	"github-triage/config"
	"github-triage/github"
	"github-triage/llm"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}
}
func main() {
	log.Println("🚀 Starting GitHub Triage Agent")

	cfg := config.Load()

	ghClient := github.NewClient(cfg.GithubToken, cfg.RepoOwner, cfg.RepoName)
	analyzer := llm.NewOpenAIClient()

	issues, err := ghClient.FetchOpenIssues()
	if err != nil {
		log.Fatal("[MAIN][ERROR] Failed to fetch issues:", err)
	}

	if len(issues) == 0 {
		log.Println("[MAIN] No open issues found. Exiting.")
		return
	}

	for _, issue := range issues {
		log.Printf("[MAIN] Processing issue #%d: %s\n", issue.Number, issue.Title)

		analysis, err := analyzer.Analyze(issue.Title + " " + issue.Body)
		if err != nil {
			log.Println("[MAIN][ERROR] Analysis failed:", err)
			continue
		}

		// Add label
		err = ghClient.AddLabel(issue.Number, analysis.Category)
		if err != nil {
			log.Println("[MAIN][ERROR] Failed to add label:", err)
		} else {
			log.Printf("[MAIN] Label '%s' added to issue #%d\n", analysis.Category, issue.Number)
		}

		// Build comment
		comment := fmt.Sprintf(
			"### 🤖 AI Triage\n\n**Summary:** %s\n\n**Action Items:**\n- %s",
			analysis.Summary,
			strings.Join(analysis.ActionItems, "\n- "),
		)

		// Add comment
		err = ghClient.Comment(issue.Number, comment)
		if err != nil {
			log.Println("[MAIN][ERROR] Failed to add comment:", err)
		} else {
			log.Printf("[MAIN] Comment added to issue #%d\n", issue.Number)
		}
	}

	log.Println("✅ Processing completed")
}
