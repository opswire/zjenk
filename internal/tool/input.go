package tool

type (
	// JobPathInput встраивается во все инпуты, где job_path обязателен.
	JobPathInput struct {
		JobPath string `json:"job_path" jsonschema:"[Required] Путь к джобе без префикса /job/, например: my-job или folder/subfolder/my-job"`
	}

	// ProjectNameInput встраивается в инпуты listJobs, listNodes, getQueue.
	// project_name подставляется как поддомен Jenkins URL: specific-name.jenkins.domain.com
	ProjectNameInput struct {
		ProjectName string `json:"project_name" jsonschema:"[Required] Имя проекта, подставляется как поддомен Jenkins URL, например: my-project"`
	}

	ListJobsInput struct {
		ProjectNameInput
	}

	GetJobInput struct {
		JobPathInput
	}

	ListBuildsInput struct {
		JobPathInput
	}

	GetLastBuildInput struct {
		JobPathInput
	}

	BuildRefInput struct {
		JobPathInput
		BuildNumber int64 `json:"build_number" jsonschema:"[Required] Номер сборки (целое положительное число)"`
	}

	TriggerBuildInput struct {
		JobPathInput
		Params map[string]string `json:"params,omitempty" jsonschema:"[Optional] Параметры запуска сборки в формате ключ-значение"`
	}

	SearchLogInput struct {
		JobPathInput
		BuildNumber int64  `json:"build_number" jsonschema:"[Required] Номер сборки (целое положительное число)"`
		Pattern     string `json:"pattern"      jsonschema:"[Required] Регулярное выражение или подстрока для поиска в логе сборки"`
	}

	ListNodesInput struct {
		ProjectNameInput
	}

	GetQueueInput struct {
		ProjectNameInput
	}
)
