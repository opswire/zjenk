package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"mcp-jenkins/internal/config"
)

type Jenkins struct {
	http *resty.Client
}

type JobInfo struct {
	Name    string
	URL     string
	Color   string
	InQueue bool
}

type BuildInfo struct {
	Number    int64
	URL       string
	Result    string
	Building  bool
	Duration  float64
	Timestamp int64
	Causes    []string
}

type NodeInfo struct {
	Name               string
	Offline            bool
	TemporarilyOffline bool
	NumExecutors       int64
	OfflineCauseReason string
}

type QueueItem struct {
	ID    int64
	Task  string
	Why   string
	Stuck bool
}

// Jenkins REST API response shapes (unexported).

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

func New(cfg *config.Config) (*Jenkins, error) {
	r := resty.New().
		SetBaseURL(cfg.Jenkins.URL).
		SetBasicAuth(cfg.Jenkins.Username, cfg.Jenkins.Password)

	resp, err := r.R().SetContext(context.Background()).Get("/api/json")
	if err != nil {
		return nil, fmt.Errorf("connect to jenkins at %s: %w", cfg.Jenkins.URL, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("connect to jenkins at %s: HTTP %d", cfg.Jenkins.URL, resp.StatusCode())
	}

	return &Jenkins{http: r}, nil
}

func (j *Jenkins) ListJobs(ctx context.Context) ([]JobInfo, error) {
	var data apiJobsResponse
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		SetQueryParam("tree", "jobs[name,url,color,inQueue]").
		Get("/api/json")
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("list jobs: HTTP %d", resp.StatusCode())
	}

	result := make([]JobInfo, 0, len(data.Jobs))
	for _, jb := range data.Jobs {
		result = append(result, JobInfo{Name: jb.Name, URL: jb.URL, Color: jb.Color, InQueue: jb.InQueue})
	}
	return result, nil
}

func (j *Jenkins) GetJob(ctx context.Context, name string) (*JobInfo, error) {
	var data apiJob
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		Get(jobPath(name) + "/api/json")
	if err != nil {
		return nil, fmt.Errorf("get job %q: %w", name, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get job %q: HTTP %d", name, resp.StatusCode())
	}

	return &JobInfo{Name: data.Name, URL: data.URL, Color: data.Color, InQueue: data.InQueue}, nil
}

func (j *Jenkins) ListBuilds(ctx context.Context, jobName string) ([]BuildInfo, error) {
	var data apiBuildsResponse
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		SetQueryParam("tree", "builds[number,url,result,building,duration,timestamp,actions[causes[shortDescription]]]").
		Get(jobPath(jobName) + "/api/json")
	if err != nil {
		return nil, fmt.Errorf("list builds for job %q: %w", jobName, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("list builds for job %q: HTTP %d", jobName, resp.StatusCode())
	}

	result := make([]BuildInfo, 0, len(data.Builds))
	for _, b := range data.Builds {
		result = append(result, buildInfoFrom(b))
	}
	return result, nil
}

func (j *Jenkins) GetBuild(ctx context.Context, jobName string, number int64) (*BuildInfo, error) {
	var data apiBuild
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		Get(fmt.Sprintf("%s/%d/api/json", jobPath(jobName), number))
	if err != nil {
		return nil, fmt.Errorf("get build #%d for job %q: %w", number, jobName, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get build #%d for job %q: HTTP %d", number, jobName, resp.StatusCode())
	}

	info := buildInfoFrom(data)
	return &info, nil
}

func (j *Jenkins) GetBuildLog(ctx context.Context, jobName string, number int64) (string, error) {
	resp, err := j.http.R().
		SetContext(ctx).
		Get(fmt.Sprintf("%s/%d/consoleText", jobPath(jobName), number))
	if err != nil {
		return "", fmt.Errorf("get log for build #%d of job %q: %w", number, jobName, err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("get log for build #%d of job %q: HTTP %d", number, jobName, resp.StatusCode())
	}

	return resp.String(), nil
}

func (j *Jenkins) TriggerBuild(ctx context.Context, jobName string, params map[string]string) (int64, error) {
	path := jobPath(jobName)
	endpoint := path + "/build"
	if len(params) > 0 {
		endpoint = path + "/buildWithParameters"
	}

	req := j.http.R().SetContext(ctx)
	if len(params) > 0 {
		req = req.SetFormData(params)
	}

	resp, err := req.Post(endpoint)
	if err != nil {
		return 0, fmt.Errorf("trigger build for job %q: %w", jobName, err)
	}
	if resp.IsError() {
		return 0, fmt.Errorf("trigger build for job %q: HTTP %d", jobName, resp.StatusCode())
	}

	queueID, err := parseQueueID(resp.Header().Get("Location"))
	if err != nil {
		return 0, fmt.Errorf("parse queue location: %w", err)
	}
	return queueID, nil
}

func (j *Jenkins) StopBuild(ctx context.Context, jobName string, number int64) error {
	resp, err := j.http.R().
		SetContext(ctx).
		Post(fmt.Sprintf("%s/%d/stop", jobPath(jobName), number))
	if err != nil {
		return fmt.Errorf("stop build #%d for job %q: %w", number, jobName, err)
	}
	if resp.IsError() {
		return fmt.Errorf("stop build #%d for job %q: HTTP %d", number, jobName, resp.StatusCode())
	}
	return nil
}

func (j *Jenkins) ListNodes(ctx context.Context) ([]NodeInfo, error) {
	var data apiNodesResponse
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		SetQueryParam("tree", "computer[displayName,offline,temporarilyOffline,numExecutors,offlineCauseReason]").
		Get("/computer/api/json")
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("list nodes: HTTP %d", resp.StatusCode())
	}

	result := make([]NodeInfo, 0, len(data.Computer))
	for _, n := range data.Computer {
		result = append(result, NodeInfo{
			Name:               n.DisplayName,
			Offline:            n.Offline,
			TemporarilyOffline: n.TemporarilyOffline,
			NumExecutors:       n.NumExecutors,
			OfflineCauseReason: n.OfflineCauseReason,
		})
	}
	return result, nil
}

func (j *Jenkins) GetQueue(ctx context.Context) ([]QueueItem, error) {
	var data apiQueueResponse
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		Get("/queue/api/json")
	if err != nil {
		return nil, fmt.Errorf("get queue: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get queue: HTTP %d", resp.StatusCode())
	}

	result := make([]QueueItem, 0, len(data.Items))
	for _, item := range data.Items {
		result = append(result, QueueItem{
			ID:    item.ID,
			Task:  item.Task.Name,
			Why:   item.Why,
			Stuck: item.Stuck,
		})
	}
	return result, nil
}

func buildInfoFrom(b apiBuild) BuildInfo {
	return BuildInfo{
		Number:    b.Number,
		URL:       b.URL,
		Result:    b.Result,
		Building:  b.Building,
		Duration:  b.Duration,
		Timestamp: b.Timestamp,
		Causes:    buildCauses(b),
	}
}

func buildCauses(b apiBuild) []string {
	var causes []string
	for _, action := range b.Actions {
		for _, cause := range action.Causes {
			if cause.ShortDescription != "" {
				causes = append(causes, cause.ShortDescription)
			}
		}
	}
	return causes
}

// jobPath converts "folder/job" → "/job/folder/job/job" for Jenkins folder support.
func jobPath(name string) string {
	parts := strings.Split(name, "/")
	return "/job/" + strings.Join(parts, "/job/")
}

// parseQueueID extracts the numeric ID from a Jenkins Location header like ".../queue/item/42/".
func parseQueueID(location string) (int64, error) {
	trimmed := strings.TrimRight(location, "/")
	if trimmed == "" {
		return 0, fmt.Errorf("empty location header")
	}
	parts := strings.Split(trimmed, "/")
	return strconv.ParseInt(parts[len(parts)-1], 10, 64)
}
