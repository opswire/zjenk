package client

// Jenkins REST API wire types — used only inside the client package.

type apiJob struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Color   string `json:"color"`
	InQueue bool   `json:"inQueue"`
}

type apiJobsResponse struct {
	Jobs []apiJob `json:"jobs"`
}

type apiBuildAction struct {
	Causes []struct {
		ShortDescription string `json:"shortDescription"`
	} `json:"causes"`
}

type apiBuild struct {
	Number    int64            `json:"number"`
	URL       string           `json:"url"`
	Result    string           `json:"result"`
	Building  bool             `json:"building"`
	Duration  float64          `json:"duration"`
	Timestamp int64            `json:"timestamp"`
	Actions   []apiBuildAction `json:"actions"`
}

type apiBuildsResponse struct {
	Builds []apiBuild `json:"builds"`
}

type apiNodesResponse struct {
	Computer []struct {
		DisplayName        string `json:"displayName"`
		Offline            bool   `json:"offline"`
		TemporarilyOffline bool   `json:"temporarilyOffline"`
		NumExecutors       int64  `json:"numExecutors"`
		OfflineCauseReason string `json:"offlineCauseReason"`
	} `json:"computer"`
}

type apiQueueResponse struct {
	Items []struct {
		ID   int64 `json:"id"`
		Task struct {
			Name string `json:"name"`
		} `json:"task"`
		Why   string `json:"why"`
		Stuck bool   `json:"stuck"`
	} `json:"items"`
}
