mcp-jenkins/
├── cmd/main.go                  — точка входа: fx.New, все fx.Provide/fx.Invoke, регистрация тулов
├── config.yaml                  — URL/credentials Jenkins + MCP name/version + конфиг тулов
├── internal/
│   ├── config/config.go         — загрузка конфига, валидация; Config.Tool(key) — lookup с дефолтами
│   ├── client/jenkins.go        — клиентский слой (обёртка над gojenkins + resty)
│   └── tool/
│       ├── job.go               — jenkins_list_jobs, jenkins_get_job; toolJSON/toolError helpers
│       ├── build.go             — list/get/trigger/stop/log builds
│       ├── node.go              — jenkins_list_nodes, jenkins_get_queue
│       └── search.go            — jenkins_search_log (regex/substring поиск по логу сборки)
└── server/server.go             — New() создаёт mcp.Server, Run() вешает fx lifecycle

Инструменты (10 шт.):

jenkins_list_jobs / jenkins_get_job
jenkins_list_builds / jenkins_get_build / jenkins_get_build_log
jenkins_trigger_build / jenkins_stop_build
jenkins_list_nodes / jenkins_get_queue
jenkins_search_log

Паттерны:

Весь uber/fx (fx.New, fx.Module, fx.Provide, fx.Invoke) — только в cmd/main.go
Регистрация тулов (Register вызовы) — только в cmd/main.go через fx.Invoke
Каждая группа тулов — struct с Register(s *mcp.Server); описания берутся из cfg.Tool(key)
config.yaml содержит секцию tools: с is-enabled, name, description.ru/en для каждого тула
Если тул disabled (is-enabled: false) — Register его пропускает
Клиентский слой возвращает типизированные domain-объекты (JobInfo, BuildInfo, NodeInfo)
Ошибки инструментов идут в IsError: true результат

Запуск:

go build -o mcp-jenkins ./cmd/
./mcp-jenkins

Для подключения к Claude Desktop добавь в claude_desktop_config.json:

{"mcpServers": {"jenkins": {"command": "/path/to/mcp-jenkins"}}}
