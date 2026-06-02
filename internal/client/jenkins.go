package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"mcp-jenkins/internal/config"
	"mcp-jenkins/internal/dto"
)

type Jenkins struct {
	http *resty.Client
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

func (j *Jenkins) ListJobs(ctx context.Context) ([]dto.Job, error) {
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

	result := make([]dto.Job, len(data.Jobs))
	for i, jb := range data.Jobs {
		result[i] = dto.Job{Name: jb.Name, URL: jb.URL, Color: jb.Color, InQueue: jb.InQueue}
	}
	return result, nil
}

func (j *Jenkins) GetJob(ctx context.Context, jobURL string) (*dto.Job, error) {
	var data apiJob
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		Get(normalizeURL(jobURL) + "/api/json")
	if err != nil {
		return nil, fmt.Errorf("get job %q: %w", jobURL, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get job %q: HTTP %d", jobURL, resp.StatusCode())
	}

	return &dto.Job{Name: data.Name, URL: data.URL, Color: data.Color, InQueue: data.InQueue}, nil
}

func (j *Jenkins) ListBuilds(ctx context.Context, jobURL string) ([]dto.Build, error) {
	var data apiBuildsResponse
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		SetQueryParam("tree", "builds[number,url,result,building,duration,timestamp,actions[causes[shortDescription]]]").
		Get(normalizeURL(jobURL) + "/api/json")
	if err != nil {
		return nil, fmt.Errorf("list builds for %q: %w", jobURL, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("list builds for %q: HTTP %d", jobURL, resp.StatusCode())
	}

	result := make([]dto.Build, len(data.Builds))
	for i, b := range data.Builds {
		result[i] = buildFromAPI(b)
	}
	return result, nil
}

func (j *Jenkins) GetBuild(ctx context.Context, jobURL string, number int64) (*dto.Build, error) {
	var data apiBuild
	resp, err := j.http.R().
		SetContext(ctx).
		SetResult(&data).
		Get(fmt.Sprintf("%s/%d/api/json", normalizeURL(jobURL), number))
	if err != nil {
		return nil, fmt.Errorf("get build #%d for %q: %w", number, jobURL, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get build #%d for %q: HTTP %d", number, jobURL, resp.StatusCode())
	}

	b := buildFromAPI(data)
	return &b, nil
}

func (j *Jenkins) GetBuildLog(ctx context.Context, jobURL string, number int64) (string, error) {
	resp, err := j.http.R().
		SetContext(ctx).
		Get(fmt.Sprintf("%s/%d/consoleText", normalizeURL(jobURL), number))
	if err != nil {
		return "", fmt.Errorf("get log for build #%d of %q: %w", number, jobURL, err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("get log for build #%d of %q: HTTP %d", number, jobURL, resp.StatusCode())
	}

	return resp.String(), nil
}

func (j *Jenkins) TriggerBuild(ctx context.Context, jobURL string, params map[string]string) (int64, error) {
	base := normalizeURL(jobURL)
	endpoint := base + "/build"
	if len(params) > 0 {
		endpoint = base + "/buildWithParameters"
	}

	req := j.http.R().SetContext(ctx)
	if len(params) > 0 {
		req = req.SetFormData(params)
	}

	resp, err := req.Post(endpoint)
	if err != nil {
		return 0, fmt.Errorf("trigger build for %q: %w", jobURL, err)
	}
	if resp.IsError() {
		return 0, fmt.Errorf("trigger build for %q: HTTP %d", jobURL, resp.StatusCode())
	}

	queueID, err := parseQueueID(resp.Header().Get("Location"))
	if err != nil {
		return 0, fmt.Errorf("parse queue location: %w", err)
	}
	return queueID, nil
}

func (j *Jenkins) StopBuild(ctx context.Context, jobURL string, number int64) error {
	resp, err := j.http.R().
		SetContext(ctx).
		Post(fmt.Sprintf("%s/%d/stop", normalizeURL(jobURL), number))
	if err != nil {
		return fmt.Errorf("stop build #%d for %q: %w", number, jobURL, err)
	}
	if resp.IsError() {
		return fmt.Errorf("stop build #%d for %q: HTTP %d", number, jobURL, resp.StatusCode())
	}
	return nil
}

func (j *Jenkins) ListNodes(ctx context.Context) ([]dto.Node, error) {
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

	result := make([]dto.Node, len(data.Computer))
	for i, n := range data.Computer {
		result[i] = dto.Node{
			Name:               n.DisplayName,
			Offline:            n.Offline,
			TemporarilyOffline: n.TemporarilyOffline,
			NumExecutors:       n.NumExecutors,
			OfflineCauseReason: n.OfflineCauseReason,
		}
	}
	return result, nil
}

func (j *Jenkins) GetQueue(ctx context.Context) ([]dto.QueueItem, error) {
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

	result := make([]dto.QueueItem, len(data.Items))
	for i, item := range data.Items {
		result[i] = dto.QueueItem{ID: item.ID, Task: item.Task.Name, Why: item.Why, Stuck: item.Stuck}
	}
	return result, nil
}

func buildFromAPI(b apiBuild) dto.Build {
	return dto.Build{
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

// normalizeURL strips a trailing slash so URL + "/api/json" is always valid.
func normalizeURL(rawURL string) string {
	return strings.TrimRight(rawURL, "/")
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
