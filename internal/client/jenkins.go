package client

import (
	"context"
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/go-resty/resty/v2"

	"mcp-jenkins/internal/config"
)

// Jenkins wraps gojenkins and resty for Jenkins API access.
type Jenkins struct {
	core  *gojenkins.Jenkins
	resty *resty.Client
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

func New(cfg *config.Config) (*Jenkins, error) {
	ctx := context.Background()

	core, err := gojenkins.CreateJenkins(
		nil,
		cfg.Jenkins.URL,
		cfg.Jenkins.Username,
		cfg.Jenkins.Password,
	).Init(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to jenkins at %s: %w", cfg.Jenkins.URL, err)
	}

	restyClient := resty.New().
		SetBaseURL(cfg.Jenkins.URL).
		SetBasicAuth(cfg.Jenkins.Username, cfg.Jenkins.Password).
		SetHeader("Content-Type", "application/json")

	return &Jenkins{core: core, resty: restyClient}, nil
}

func (j *Jenkins) ListJobs(ctx context.Context) ([]JobInfo, error) {
	jobs, err := j.core.GetAllJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}

	result := make([]JobInfo, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, JobInfo{
			Name:    job.GetName(),
			URL:     job.Raw.URL,
			Color:   job.Raw.Color,
			InQueue: job.Raw.InQueue,
		})
	}
	return result, nil
}

func (j *Jenkins) GetJob(ctx context.Context, name string) (*JobInfo, error) {
	job, err := j.core.GetJob(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get job %q: %w", name, err)
	}

	return &JobInfo{
		Name:    job.GetName(),
		URL:     job.Raw.URL,
		Color:   job.Raw.Color,
		InQueue: job.Raw.InQueue,
	}, nil
}

func (j *Jenkins) ListBuilds(ctx context.Context, jobName string) ([]BuildInfo, error) {
	job, err := j.core.GetJob(ctx, jobName)
	if err != nil {
		return nil, fmt.Errorf("get job %q: %w", jobName, err)
	}

	builds, err := job.GetAllBuildIds(ctx)
	if err != nil {
		return nil, fmt.Errorf("list builds for job %q: %w", jobName, err)
	}

	result := make([]BuildInfo, 0, len(builds))
	for _, b := range builds {
		build, err := job.GetBuild(ctx, b.Number)
		if err != nil {
			continue
		}
		result = append(result, buildInfoFrom(build))
	}
	return result, nil
}

func (j *Jenkins) GetBuild(ctx context.Context, jobName string, number int64) (*BuildInfo, error) {
	job, err := j.core.GetJob(ctx, jobName)
	if err != nil {
		return nil, fmt.Errorf("get job %q: %w", jobName, err)
	}

	build, err := job.GetBuild(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("get build #%d for job %q: %w", number, jobName, err)
	}

	info := buildInfoFrom(build)
	return &info, nil
}

func (j *Jenkins) GetBuildLog(ctx context.Context, jobName string, number int64) (string, error) {
	job, err := j.core.GetJob(ctx, jobName)
	if err != nil {
		return "", fmt.Errorf("get job %q: %w", jobName, err)
	}

	build, err := job.GetBuild(ctx, number)
	if err != nil {
		return "", fmt.Errorf("get build #%d for job %q: %w", number, jobName, err)
	}

	log := build.GetConsoleOutput(ctx)
	return log, nil
}

func (j *Jenkins) TriggerBuild(ctx context.Context, jobName string, params map[string]string) (int64, error) {
	queueID, err := j.core.BuildJob(ctx, jobName, params)
	if err != nil {
		return 0, fmt.Errorf("trigger build for job %q: %w", jobName, err)
	}
	return queueID, nil
}

func (j *Jenkins) StopBuild(ctx context.Context, jobName string, number int64) error {
	job, err := j.core.GetJob(ctx, jobName)
	if err != nil {
		return fmt.Errorf("get job %q: %w", jobName, err)
	}

	build, err := job.GetBuild(ctx, number)
	if err != nil {
		return fmt.Errorf("get build #%d for job %q: %w", number, jobName, err)
	}

	_, err = build.Stop(ctx)
	if err != nil {
		return fmt.Errorf("stop build #%d for job %q: %w", number, jobName, err)
	}
	return nil
}

func (j *Jenkins) ListNodes(ctx context.Context) ([]NodeInfo, error) {
	nodes, err := j.core.GetAllNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	result := make([]NodeInfo, 0, len(nodes))
	for _, n := range nodes {
		result = append(result, NodeInfo{
			Name:               n.GetName(),
			Offline:            n.Raw.Offline,
			TemporarilyOffline: n.Raw.TemporarilyOffline,
			NumExecutors:       n.Raw.NumExecutors,
			OfflineCauseReason: n.Raw.OfflineCauseReason,
		})
	}
	return result, nil
}

func (j *Jenkins) GetQueue(ctx context.Context) ([]QueueItem, error) {
	queue, err := j.core.GetQueue(ctx)
	if err != nil {
		return nil, fmt.Errorf("get queue: %w", err)
	}

	result := make([]QueueItem, 0, len(queue.Raw.Items))
	for _, item := range queue.Raw.Items {
		result = append(result, QueueItem{
			ID:    item.ID,
			Task:  item.Task.Name,
			Why:   item.Why,
			Stuck: item.Stuck,
		})
	}
	return result, nil
}

func buildInfoFrom(b *gojenkins.Build) BuildInfo {
	return BuildInfo{
		Number:    b.Raw.Number,
		URL:       b.Raw.URL,
		Result:    b.Raw.Result,
		Building:  b.Raw.Building,
		Duration:  b.Raw.Duration,
		Timestamp: b.Raw.Timestamp,
		Causes:    buildCauses(b),
	}
}

func buildCauses(b *gojenkins.Build) []string {
	var causes []string
	for _, action := range b.Raw.Actions {
		for _, cause := range action.Causes {
			if desc, ok := cause["shortDescription"].(string); ok && desc != "" {
				causes = append(causes, desc)
			}
		}
	}
	return causes
}
