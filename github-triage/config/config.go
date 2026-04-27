package config

import (
	"log"
	"os"
)

type Config struct {
	GithubToken string
	RepoOwner   string
	RepoName    string
}

func Load() *Config {
	cfg := &Config{
		GithubToken: os.Getenv("GITHUB_TOKEN"),
		RepoOwner:   os.Getenv("REPO_OWNER"),
		RepoName:    os.Getenv("REPO_NAME"),
	}

	if cfg.GithubToken == "" || cfg.RepoOwner == "" || cfg.RepoName == "" {
		log.Fatal("Missing required environment variables: GITHUB_TOKEN, REPO_OWNER, REPO_NAME")
	}

	return cfg
}
