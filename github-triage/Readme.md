# 🤖 GitHub Triage Agent (MVP)

An AI-powered GitHub agent that automatically analyzes issues, labels them, and adds actionable summaries.

Built with **Golang + LLM (OpenAI)** as a lightweight, production-style MVP.

---

## 🚀 Features

- 🔍 Fetches open issues from a GitHub repository  
- 🧠 Uses LLM to analyze issue content  
- 🏷️ Automatically adds labels (bug, feature, docs, other)  
- 💬 Posts structured summary + action items as comments  
- 📦 Dockerized for easy execution  
- 📊 Standard logging for debugging and observability  

---

## 🏗️ Architecture (MVP)
GitHub API → Go Agent → LLM → GitHub API


Simple, fast, and focused on delivering value.

---

---

## ⚙️ Setup

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd github-triage
```
2. Create .env

```bash
cp .env.example .env
```

Update Values:

```bash
GITHUB_TOKEN=your_github_token
REPO_OWNER=your_username
REPO_NAME=your_repo
OPENAI_API_KEY=your_openai_key
OPENAI_MODEL=gpt-4o-mini
```

3. Run locally

```bash
go mod tidy
go run main.go
```
🐳 Run with Docker

```bash
docker build -t github-triage .
```
Run container
```bash
docker run --env-file .env github-triage
```

🤖 Example Output
Label added:
```
bug
```

📜 License

MIT