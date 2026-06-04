package client

// Jenkins REST API wire types — used only inside the client package.

type (
	parameterDefinition struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		DefaultParameterValue struct {
			Name  string      `json:"name"`
			Value interface{} `json:"value"`
		} `json:"defaultParameterValue"`
	}

	jobProperty struct {
		ParameterDefinitions []parameterDefinition `json:"parameterDefinitions"`
	}

	jobBuild struct {
		Number int64  `json:"number"`
		URL    string `json:"url"`
	}

	job struct {
		Name            string        `json:"name"`
		URL             string        `json:"url"`
		Color           string        `json:"color"`
		InQueue         bool          `json:"inQueue"`
		Buildable       bool          `json:"buildable"`
		Description     string        `json:"description"`
		NextBuildNumber int64         `json:"nextBuildNumber"`
		LastBuild       *jobBuild     `json:"lastBuild"`
		Property        []jobProperty `json:"property"`
	}

	jobsResponse struct {
		Jobs []job `json:"jobs"`
	}

	buildAction struct {
		Causes []struct {
			ShortDescription string `json:"shortDescription"`
		} `json:"causes"`
	}

	build struct {
		Number            int64         `json:"number"`
		URL               string        `json:"url"`
		DisplayName       string        `json:"displayName"`
		Result            string        `json:"result"`
		Building          bool          `json:"building"`
		Duration          float64       `json:"duration"`
		EstimatedDuration float64       `json:"estimatedDuration"`
		Timestamp         int64         `json:"timestamp"`
		QueueID           int64         `json:"queueId"`
		BuiltOn           string        `json:"builtOn"`
		Actions           []buildAction `json:"actions"`
	}

	buildsResponse struct {
		Builds []build `json:"builds"`
	}

	nodesResponse struct {
		Computer []struct {
			DisplayName        string `json:"displayName"`
			Offline            bool   `json:"offline"`
			TemporarilyOffline bool   `json:"temporarilyOffline"`
			Idle               bool   `json:"idle"`
			NumExecutors       int64  `json:"numExecutors"`
			OfflineCauseReason string `json:"offlineCauseReason"`
		} `json:"computer"`
	}

	queueResponse struct {
		Items []struct {
			ID   int64 `json:"id"`
			Task struct {
				Name string `json:"name"`
			} `json:"task"`
			Why          string `json:"why"`
			Stuck        bool   `json:"stuck"`
			Blocked      bool   `json:"blocked"`
			Buildable    bool   `json:"buildable"`
			InQueueSince int64  `json:"inQueueSince"`
		} `json:"items"`
	}
)
