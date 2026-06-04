package dto

type (
	ParameterDefinition struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Description  string `json:"description,omitempty"`
		DefaultValue string `json:"default_value,omitempty"`
	}

	JobBuild struct {
		Number int64  `json:"number"`
		URL    string `json:"url"`
	}

	Job struct {
		Name            string                `json:"name"`
		URL             string                `json:"url"`
		Color           string                `json:"color,omitempty"`
		InQueue         bool                  `json:"in_queue"`
		Buildable       bool                  `json:"buildable"`
		Description     string                `json:"description,omitempty"`
		NextBuildNumber int64                 `json:"next_build_number"`
		LastBuild       *JobBuild             `json:"last_build,omitempty"`
		Parameters      []ParameterDefinition `json:"parameters,omitempty"`
	}

	Build struct {
		Number            int64    `json:"number"`
		URL               string   `json:"url"`
		DisplayName       string   `json:"display_name,omitempty"`
		Result            string   `json:"result,omitempty"`
		Building          bool     `json:"building"`
		Duration          float64  `json:"duration_ms"`
		EstimatedDuration float64  `json:"estimated_duration_ms,omitempty"`
		Timestamp         int64    `json:"timestamp_ms"`
		QueueID           int64    `json:"queue_id,omitempty"`
		BuiltOn           string   `json:"built_on,omitempty"`
		Causes            []string `json:"causes,omitempty"`
	}

	Node struct {
		Name               string `json:"name"`
		Offline            bool   `json:"offline"`
		TemporarilyOffline bool   `json:"temporarily_offline"`
		Idle               bool   `json:"idle"`
		NumExecutors       int64  `json:"num_executors"`
		OfflineCauseReason string `json:"offline_cause_reason,omitempty"`
	}

	QueueItem struct {
		ID           int64  `json:"id"`
		Task         string `json:"task"`
		Why          string `json:"why,omitempty"`
		Stuck        bool   `json:"stuck"`
		Blocked      bool   `json:"blocked"`
		Buildable    bool   `json:"buildable"`
		InQueueSince int64  `json:"in_queue_since_ms,omitempty"`
	}

	TriggerResult struct {
		JobPath string `json:"job_path"`
		QueueID int64  `json:"queue_id"`
	}

	SearchMatch struct {
		Line    int    `json:"line"`
		Content string `json:"content"`
	}
)
