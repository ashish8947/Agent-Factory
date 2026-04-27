package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github-triage/config"
	"github-triage/github"
	"github-triage/llm"
	"github-triage/similarity"

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

	// 🔥 Load embedding cache (persistent)
	similarity.LoadCache()
	defer similarity.SaveCache()

	cfg := config.Load()

	ghClient := github.NewClient(cfg.GithubToken, cfg.RepoOwner, cfg.RepoName)
	analyzer := llm.NewOpenAIClient()

	issueNumber := os.Getenv("ISSUE_NUMBER")

	issues, err := ghClient.FetchOpenIssues()
	if err != nil {
		log.Fatal("[MAIN][ERROR] Failed to fetch issues:", err)
	}

	if len(issues) == 0 {
		log.Println("[MAIN] No open issues found. Exiting.")
		return
	}

	for _, issue := range issues {

		// ✅ Process only triggered issue (GitHub Actions)
		if issueNumber != "" && fmt.Sprintf("%d", issue.Number) != issueNumber {
			continue
		}

		log.Printf("[MAIN] Processing issue #%d: %s\n", issue.Number, issue.Title)

		issueText := issue.Title + " " + issue.Body

		// 🔹 Step 1: LLM Analysis
		analysis, err := analyzer.Analyze(issueText)
		if err != nil {
			log.Println("[MAIN][ERROR] Analysis failed:", err)
			continue
		}

		// 🔹 Step 2: Fetch recent issues
		recentIssues, err := ghClient.FetchRecentIssues(50)
		if err != nil {
			log.Println("[MAIN][ERROR] Fetch recent issues failed:", err)
			continue
		}

		// 🔹 Step 3: Prepare batch texts
		texts := []string{issueText}

		filteredIssues := []int{} // store indexes of valid issues

		for _, other := range recentIssues {
			if other.Number == issue.Number {
				continue
			}
			texts = append(texts, other.Title+" "+other.Body)
			filteredIssues = append(filteredIssues, other.Number)
		}

		// 🔹 Step 4: Get embeddings in ONE call
		embeddings, err := similarity.GetEmbeddings(texts)
		if err != nil {
			log.Println("[MAIN][ERROR] Batch embedding failed:", err)
			continue
		}

		currentEmbedding := embeddings[0]

		// 🔹 Step 5: Duplicate detection
		duplicates := []string{}
		idx := 1

		for _, other := range recentIssues {

			if other.Number == issue.Number {
				continue
			}

			score := similarity.CosineSimilarity(currentEmbedding, embeddings[idx])
			idx++

			if score > 0.85 {
				duplicates = append(duplicates,
					fmt.Sprintf("#%d %s", other.Number, other.Title))
			}
		}

		log.Printf("[MAIN] Found %d possible duplicates\n", len(duplicates))

		// 🔹 Step 6: Quality check
		needsMoreInfo, err := analyzer.CheckQuality(issueText)
		if err != nil {
			log.Println("[MAIN][ERROR] Quality check failed:", err)
		}

		// 🔹 Step 7: Add label
		err = ghClient.AddLabel(issue.Number, analysis.Category)
		if err != nil {
			log.Println("[MAIN][ERROR] Failed to add label:", err)
		} else {
			log.Printf("[MAIN] Label '%s' added to issue #%d\n",
				analysis.Category, issue.Number)
		}

		// 🔹 Step 8: Build single smart comment
		var commentBuilder strings.Builder

		commentBuilder.WriteString("### 🤖 AI Triage\n\n")

		commentBuilder.WriteString("**Summary:** ")
		commentBuilder.WriteString(analysis.Summary)
		commentBuilder.WriteString("\n\n")

		if len(duplicates) > 0 {
			commentBuilder.WriteString("**Possible Duplicates:**\n")
			for _, d := range duplicates {
				commentBuilder.WriteString("- ")
				commentBuilder.WriteString(d)
				commentBuilder.WriteString("\n")
			}
			commentBuilder.WriteString("\n")
		}

		if needsMoreInfo {
			commentBuilder.WriteString("**Suggestion:** Please add:\n")
			commentBuilder.WriteString("- Steps to reproduce\n")
			commentBuilder.WriteString("- Expected vs actual behavior\n")
			commentBuilder.WriteString("- Logs/screenshots\n")
		}

		comment := commentBuilder.String()

		// 🔹 Step 9: Post comment
		err = ghClient.Comment(issue.Number, comment)
		if err != nil {
			log.Println("[MAIN][ERROR] Failed to add comment:", err)
		} else {
			log.Printf("[MAIN] Comment added to issue #%d\n", issue.Number)
		}
	}

	log.Println("✅ Processing completed")
}
