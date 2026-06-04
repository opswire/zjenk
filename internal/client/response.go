package client

// Jenkins REST API wire types — used only inside the client package.

type apiParameterDefinition struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	DefaultParameterValue struct {
		Value interface{} `json:"value"`
	} `json:"defaultParameterValue"`
}

type apiJobProperty struct {
	ParameterDefinitions []apiParameterDefinition `json:"parameterDefinitions"`
}

type apiJob struct {
	Name     string           `json:"name"`
	URL      string           `json:"url"`
	Color    string           `json:"color"`
	InQueue  bool             `json:"inQueue"`
	Property []apiJobProperty `json:"property"`
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
	Duration  int64            `json:"duration"`
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
