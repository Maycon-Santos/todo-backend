package server

type CheckItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Done  bool   `json:"done"`
}

type Task struct {
	ID          string      `json:"id"`
	Done        bool        `json:"done"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Checklist   []CheckItem `json:"checklist"`
}
