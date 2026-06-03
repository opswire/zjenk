package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
)

// Run with:
//
//	JENKINS_URL=http://localhost:8080 \
//	JENKINS_USERNAME=admin \
//	JENKINS_PASSWORD=token \
//	go test ./internal/client/ -v
//
// Optional:
//
//	JENKINS_TEST_JOB=my-job        — job used for read tests (auto-discovered if empty)
//	JENKINS_ALLOW_MUTATIONS=true   — enables TriggerBuild tests

type JenkinsSuite struct {
	suite.Suite
	client    *client.Jenkins
	ctx       context.Context
	testJob   string
	testBuild int64
}

func TestJenkins(t *testing.T) {
	suite.Run(t, new(JenkinsSuite))
}

func (s *JenkinsSuite) SetupSuite() {
	url := os.Getenv("JENKINS_URL")
	if url == "" {
		s.T().Skip("set JENKINS_URL to run integration tests")
	}

	cfg := &config.Config{
		Jenkins: config.JenkinsConfig{
			URL:      url,
			Username: os.Getenv("JENKINS_USERNAME"),
			Password: os.Getenv("JENKINS_PASSWORD"),
		},
	}

	s.client = client.New(cfg)
	s.ctx = context.Background()

	s.testJob = os.Getenv("JENKINS_TEST_JOB")
	if s.testJob == "" {
		if jobs, err := s.client.ListJobs(s.ctx); err == nil && len(jobs) > 0 {
			s.testJob = jobs[0].Name
		}
	}

	if s.testJob != "" {
		if builds, err := s.client.ListBuilds(s.ctx, s.testJob); err == nil && len(builds) > 0 {
			s.testBuild = builds[0].Number
		}
	}
}

// --- Jobs ---

func (s *JenkinsSuite) TestListJobs() {
	jobs, err := s.client.ListJobs(s.ctx)
	s.NoError(err)
	s.NotNil(jobs)
	for _, j := range jobs {
		s.NotEmpty(j.Name)
		s.NotEmpty(j.URL)
	}
}

func (s *JenkinsSuite) TestGetJob() {
	if s.testJob == "" {
		s.T().Skip("no jobs available")
	}

	cases := []struct {
		name    string
		job     string
		wantErr bool
	}{
		{"existing", s.testJob, false},
		{"nonexistent", "no-such-job-xyz-99999", true},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			job, err := s.client.GetJob(s.ctx, tc.job)
			if tc.wantErr {
				s.Error(err)
				s.Nil(job)
			} else {
				s.NoError(err)
				s.Require().NotNil(job)
				s.NotEmpty(job.Name)
				s.NotEmpty(job.URL)
			}
		})
	}
}

// --- Builds ---

func (s *JenkinsSuite) TestListBuilds() {
	if s.testJob == "" {
		s.T().Skip("no jobs available")
	}

	cases := []struct {
		name    string
		job     string
		wantErr bool
	}{
		{"existing job", s.testJob, false},
		{"nonexistent job", "no-such-job-xyz-99999", true},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			builds, err := s.client.ListBuilds(s.ctx, tc.job)
			if tc.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(builds)
				for _, b := range builds {
					s.Positive(b.Number)
					s.NotEmpty(b.URL)
				}
			}
		})
	}
}

func (s *JenkinsSuite) TestGetBuild() {
	if s.testBuild == 0 {
		s.T().Skip("no builds available")
	}

	cases := []struct {
		name    string
		job     string
		number  int64
		wantErr bool
	}{
		{"valid", s.testJob, s.testBuild, false},
		{"nonexistent build number", s.testJob, 999999999, true},
		{"nonexistent job", "no-such-job-xyz-99999", 1, true},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			build, err := s.client.GetBuild(s.ctx, tc.job, tc.number)
			if tc.wantErr {
				s.Error(err)
				s.Nil(build)
			} else {
				s.NoError(err)
				s.Require().NotNil(build)
				s.Equal(tc.number, build.Number)
				s.NotEmpty(build.URL)
			}
		})
	}
}

func (s *JenkinsSuite) TestGetBuildLog() {
	if s.testBuild == 0 {
		s.T().Skip("no builds available")
	}

	cases := []struct {
		name    string
		job     string
		number  int64
		wantErr bool
	}{
		{"valid", s.testJob, s.testBuild, false},
		{"nonexistent job", "no-such-job-xyz-99999", 1, true},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			log, err := s.client.GetBuildLog(s.ctx, tc.job, tc.number)
			if tc.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotEmpty(log)
			}
		})
	}
}

func (s *JenkinsSuite) TestTriggerBuild() {
	if os.Getenv("JENKINS_ALLOW_MUTATIONS") == "" {
		s.T().Skip("set JENKINS_ALLOW_MUTATIONS=true to run mutation tests")
	}
	if s.testJob == "" {
		s.T().Skip("no jobs available")
	}

	cases := []struct {
		name    string
		job     string
		params  map[string]string
		wantErr bool
	}{
		{"no params", s.testJob, nil, false},
		{"with params", s.testJob, map[string]string{"BRANCH": "main"}, false},
		{"nonexistent job", "no-such-job-xyz-99999", nil, true},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			queueID, err := s.client.TriggerBuild(s.ctx, tc.job, tc.params)
			if tc.wantErr {
				s.Error(err)
				s.Zero(queueID)
			} else {
				s.NoError(err)
				s.Positive(queueID)
			}
		})
	}
}

// --- Nodes & Queue ---

func (s *JenkinsSuite) TestListNodes() {
	nodes, err := s.client.ListNodes(s.ctx)
	s.NoError(err)
	s.NotNil(nodes)
	for _, n := range nodes {
		s.NotEmpty(n.Name)
	}
}

func (s *JenkinsSuite) TestGetQueue() {
	items, err := s.client.GetQueue(s.ctx)
	s.NoError(err)
	s.NotNil(items)
}
