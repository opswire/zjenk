package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"mcp-jenkins/internal/config"
	"mcp-jenkins/internal/dto"
)

type Jenkins struct {
	jsonClient *resty.Client
	baseURL    string
}

func New(cfg *config.Config) *Jenkins {
	r := resty.New().
		SetBaseURL(cfg.Jenkins.URL).
		SetBasicAuth(cfg.Jenkins.Username, cfg.Jenkins.Password)

	return &Jenkins{jsonClient: r, baseURL: strings.TrimRight(cfg.Jenkins.URL, "/")}
}

// setCrumb fetches the Jenkins CSRF crumb and sets it on req.
// If the Jenkins instance has CSRF disabled, the request proceeds unchanged.
func (j *Jenkins) setCrumb(ctx context.Context, req *resty.Request) {
	var data struct {
		CrumbRequestField string `json:"crumbRequestField"`
		Crumb             string `json:"crumb"`
	}

	resp, err := j.jsonClient.R().SetContext(ctx).SetResult(&data).Get("/crumbIssuer/api/json")
	if err != nil || resp.IsError() || data.CrumbRequestField == "" {
		return
	}
	req.SetHeader(data.CrumbRequestField, data.Crumb)
}

// subdomainURL строит URL с projectName как поддоменом: "proj" + "https://jenkins.domain.com" → "https://proj.jenkins.domain.com".
func (j *Jenkins) subdomainURL(projectName string) (string, error) {
	u, err := url.Parse(j.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}

	u.Host = projectName + "." + u.Host

	return u.String(), nil
}

func (j *Jenkins) ListJobs(ctx context.Context, projectName string) ([]dto.Job, error) {
	base, err := j.subdomainURL(projectName)
	if err != nil {
		return nil, err
	}

	var data jobsResponse

	resp, err := j.jsonClient.R().
		SetContext(ctx).
		SetResult(&data).
		Get(base + "/api/json?tree=jobs[name,url,color,inQueue]")
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("list jobs: HTTP %d", resp.StatusCode())
	}

	result := make([]dto.Job, len(data.Jobs))
	for i, jb := range data.Jobs {
		result[i] = dto.Job{Name: jb.Name, URL: jb.URL, Color: jb.Color, InQueue: jb.InQueue}
	}

	return result, nil
}

func (j *Jenkins) GetJob(ctx context.Context, jobPath string) (*dto.Job, error) {
	var data job

	resp, err := j.jsonClient.R().
		SetContext(ctx).
		SetResult(&data).
		Get(jobAPIPath(jobPath) + "/api/json?tree=name,url,color,inQueue,buildable,description,nextBuildNumber,lastBuild[number,url],property[parameterDefinitions[name,type,description,defaultParameterValue[name,value]]]")
	if err != nil {
		return nil, fmt.Errorf("get job %q: %w", jobPath, err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("get job %q: HTTP %d", jobPath, resp.StatusCode())
	}

	job := &dto.Job{
		Name:            data.Name,
		URL:             data.URL,
		Color:           data.Color,
		InQueue:         data.InQueue,
		Buildable:       data.Buildable,
		Description:     data.Description,
		NextBuildNumber: data.NextBuildNumber,
	}
	if data.LastBuild != nil {
		job.LastBuild = &dto.JobBuild{Number: data.LastBuild.Number, URL: data.LastBuild.URL}
	}
	for _, prop := range data.Property {
		for _, p := range prop.ParameterDefinitions {
			def := dto.ParameterDefinition{
				Name:        p.Name,
				Type:        p.Type,
				Description: p.Description,
			}
			if p.DefaultParameterValue.Value != nil {
				def.DefaultValue = fmt.Sprintf("%v", p.DefaultParameterValue.Value)
			}
			job.Parameters = append(job.Parameters, def)
		}
	}

	return job, nil
}

func (j *Jenkins) GetLastBuild(ctx context.Context, jobPath string) (*dto.Build, error) {
	var data build

	resp, err := j.jsonClient.R().
		SetContext(ctx).
		SetResult(&data).
		Get(jobAPIPath(jobPath) + "/lastBuild/api/json?tree=number,url,displayName,result,building,duration,estimatedDuration,timestamp,queueId,builtOn,actions[causes[shortDescription]]")
	if err != nil {
		return nil, fmt.Errorf("get last build for %q: %w", jobPath, err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("get last build for %q: HTTP %d", jobPath, resp.StatusCode())
	}

	b := buildFromAPI(data)

	return &b, nil
}

func (j *Jenkins) ListBuilds(ctx context.Context, jobPath string) ([]dto.Build, error) {
	var data buildsResponse

	resp, err := j.jsonClient.R().
		SetContext(ctx).
		SetResult(&data).
		Get(jobAPIPath(jobPath) + "/api/json?tree=builds[number,url,displayName,result,building,duration,estimatedDuration,timestamp,queueId,builtOn,actions[causes[shortDescription]]]")
	if err != nil {
		return nil, fmt.Errorf("list builds for %q: %w", jobPath, err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("list builds for %q: HTTP %d", jobPath, resp.StatusCode())
	}

	result := make([]dto.Build, len(data.Builds))
	for i, b := range data.Builds {
		result[i] = buildFromAPI(b)
	}

	return result, nil
}

func (j *Jenkins) GetBuild(ctx context.Context, jobPath string, number int64) (*dto.Build, error) {
	var data build

	resp, err := j.jsonClient.R().
		SetContext(ctx).
		SetResult(&data).
		Get(fmt.Sprintf("%s/%d/api/json?tree=number,url,displayName,result,building,duration,estimatedDuration,timestamp,queueId,builtOn,actions[causes[shortDescription]]", jobAPIPath(jobPath), number))
	if err != nil {
		return nil, fmt.Errorf("get build #%d for %q: %w", number, jobPath, err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("get build #%d for %q: HTTP %d", number, jobPath, resp.StatusCode())
	}

	b := buildFromAPI(data)

	return &b, nil
}

func (j *Jenkins) GetBuildLog(ctx context.Context, jobPath string, number int64) (string, error) {
	resp, err := j.jsonClient.R().
		SetContext(ctx).
		Get(fmt.Sprintf("%s/%d/consoleText", jobAPIPath(jobPath), number))
	if err != nil {
		return "", fmt.Errorf("get log for build #%d of %q: %w", number, jobPath, err)
	} else if resp.IsError() {
		return "", fmt.Errorf("get log for build #%d of %q: HTTP %d", number, jobPath, resp.StatusCode())
	}

	return resp.String(), nil
}

func (j *Jenkins) TriggerBuild(ctx context.Context, jobPath string, params map[string]string) (int64, error) {
	base := jobAPIPath(jobPath)

	// Parameterised jobs must always use /buildWithParameters (even with empty params),
	// otherwise Jenkins returns 405.
	job, err := j.GetJob(ctx, jobPath)
	if err != nil {
		return 0, fmt.Errorf("trigger build for %q: %w", jobPath, err)
	}
	endpoint := base + "/build"
	if len(job.Parameters) > 0 || len(params) > 0 {
		endpoint = base + "/buildWithParameters"
	}

	req := j.jsonClient.R().SetContext(ctx)
	j.setCrumb(ctx, req)
	if len(params) > 0 {
		req.SetFormData(params)
	}

	resp, err := req.Post(endpoint)
	if err != nil {
		return 0, fmt.Errorf("trigger build for %q: %w", jobPath, err)
	} else if resp.IsError() {
		return 0, fmt.Errorf("trigger build for %q: HTTP %d", jobPath, resp.StatusCode())
	}

	queueID, err := parseQueueID(resp.Header().Get("Location"))
	if err != nil {
		return 0, fmt.Errorf("parse queue location: %w", err)
	}

	return queueID, nil
}

func (j *Jenkins) StopBuild(ctx context.Context, jobPath string, number int64) error {
	req := j.jsonClient.R().SetContext(ctx)
	j.setCrumb(ctx, req)

	resp, err := req.Post(fmt.Sprintf("%s/%d/stop", jobAPIPath(jobPath), number))
	if err != nil {
		return fmt.Errorf("stop build #%d for %q: %w", number, jobPath, err)
	} else if resp.IsError() {
		return fmt.Errorf("stop build #%d for %q: HTTP %d", number, jobPath, resp.StatusCode())
	}

	return nil
}

func (j *Jenkins) ListNodes(ctx context.Context, projectName string) ([]dto.Node, error) {
	base, err := j.subdomainURL(projectName)
	if err != nil {
		return nil, err
	}

	var data nodesResponse

	resp, err := j.jsonClient.R().
		SetContext(ctx).
		SetResult(&data).
		Get(base + "/computer/api/json?tree=computer[displayName,offline,temporarilyOffline,idle,numExecutors,offlineCauseReason]")
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("list nodes: HTTP %d", resp.StatusCode())
	}

	result := make([]dto.Node, len(data.Computer))
	for i, n := range data.Computer {
		result[i] = dto.Node{
			Name:               n.DisplayName,
			Offline:            n.Offline,
			TemporarilyOffline: n.TemporarilyOffline,
			Idle:               n.Idle,
			NumExecutors:       n.NumExecutors,
			OfflineCauseReason: n.OfflineCauseReason,
		}
	}

	return result, nil
}

func (j *Jenkins) GetQueue(ctx context.Context, projectName string) ([]dto.QueueItem, error) {
	base, err := j.subdomainURL(projectName)
	if err != nil {
		return nil, err
	}

	var data queueResponse

	resp, err := j.jsonClient.R().
		SetContext(ctx).
		SetResult(&data).
		Get(base + "/queue/api/json?tree=items[id,task[name],why,stuck,blocked,buildable,inQueueSince]")
	if err != nil {
		return nil, fmt.Errorf("get queue: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("get queue: HTTP %d", resp.StatusCode())
	}

	result := make([]dto.QueueItem, len(data.Items))
	for i, item := range data.Items {
		result[i] = dto.QueueItem{
			ID:           item.ID,
			Task:         item.Task.Name,
			Why:          item.Why,
			Stuck:        item.Stuck,
			Blocked:      item.Blocked,
			Buildable:    item.Buildable,
			InQueueSince: item.InQueueSince,
		}
	}

	return result, nil
}

func buildFromAPI(b build) dto.Build {
	return dto.Build{
		Number:            b.Number,
		URL:               b.URL,
		DisplayName:       b.DisplayName,
		Result:            b.Result,
		Building:          b.Building,
		Duration:          b.Duration,
		EstimatedDuration: b.EstimatedDuration,
		Timestamp:         b.Timestamp,
		QueueID:           b.QueueID,
		BuiltOn:           b.BuiltOn,
		Causes:            buildCauses(b),
	}
}

func buildCauses(b build) []string {
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

// jobAPIPath converts a user-supplied job path (e.g. "folder/job-name") to
// the Jenkins REST API path format (e.g. "/job/folder/job/job-name").
func jobAPIPath(jobPath string) string {
	parts := strings.Split(strings.Trim(jobPath, "/"), "/")

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
