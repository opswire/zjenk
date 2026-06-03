package dto

type Job struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Color   string `json:"color,omitempty"`
	InQueue bool   `json:"in_queue"`
}

type Build struct {
	Number    int64    `json:"number"`
	URL       string   `json:"url"`
	Result    string   `json:"result,omitempty"` // empty while building
	Building  bool     `json:"building"`
	Duration  float64  `json:"duration_ms"`
	Timestamp int64    `json:"timestamp_ms"`
	Causes    []string `json:"causes,omitempty"`
}

type Node struct {
	Name               string `json:"name"`
	Offline            bool   `json:"offline"`
	TemporarilyOffline bool   `json:"temporarily_offline"`
	NumExecutors       int64  `json:"num_executors"`
	OfflineCauseReason string `json:"offline_cause_reason,omitempty"`
}

type QueueItem struct {
	ID    int64  `json:"id"`
	Task  string `json:"task"`
	Why   string `json:"why,omitempty"`
	Stuck bool   `json:"stuck"`
}

type TriggerResult struct {
	JobPath string `json:"job_path"`
	QueueID int64  `json:"queue_id"`
}

type SearchMatch struct {
	Line    int    `json:"line"`
	Content string `json:"content"`
}
