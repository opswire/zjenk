package tool

// JobPathInput is embedded by all inputs that reference a Jenkins job.
type JobPathInput struct {
	JobPath string `json:"job_path" jsonschema:"[Required] Путь к джобе без префикса /job/, например: my-job или folder/subfolder/my-job"`
}

type GetJobInput struct {
	JobPathInput
}

type ListBuildsInput struct {
	JobPathInput
}

type BuildRefInput struct {
	JobPathInput
	BuildNumber int64 `json:"build_number" jsonschema:"[Required] Номер сборки (целое положительное число)"`
}

type TriggerBuildInput struct {
	JobPathInput
	Params map[string]string `json:"params,omitempty" jsonschema:"[Optional] Параметры запуска сборки в формате ключ-значение"`
}

type SearchLogInput struct {
	JobPathInput
	BuildNumber int64  `json:"build_number" jsonschema:"[Required] Номер сборки (целое положительное число)"`
	Pattern     string `json:"pattern"      jsonschema:"[Required] Регулярное выражение или подстрока для поиска в логе сборки"`
}
