package models

type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	State  string `json:"state"`
}

type Analysis struct {
	Category    string   `json:"category"`
	Summary     string   `json:"summary"`
	ActionItems []string `json:"action_items"`
}