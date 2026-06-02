package tool

type GetJobInput struct {
	JobURL string `json:"job_url" jsonschema:"description=Full Jenkins job URL (e.g. http://jenkins:8080/job/my-job)"`
}

type ListBuildsInput struct {
	JobURL string `json:"job_url" jsonschema:"description=Full Jenkins job URL"`
}

type BuildRefInput struct {
	JobURL      string `json:"job_url"      jsonschema:"description=Full Jenkins job URL"`
	BuildNumber int64  `json:"build_number" jsonschema:"description=Build number (positive integer)"`
}

type TriggerBuildInput struct {
	JobURL string            `json:"job_url"          jsonschema:"description=Full Jenkins job URL"`
	Params map[string]string `json:"params,omitempty" jsonschema:"description=Optional build parameters as key-value pairs"`
}

type SearchLogInput struct {
	JobURL      string `json:"job_url"      jsonschema:"description=Full Jenkins job URL"`
	BuildNumber int64  `json:"build_number" jsonschema:"description=Build number (positive integer)"`
	Pattern     string `json:"pattern"      jsonschema:"description=Regex pattern or substring to search for in the build log"`
}
