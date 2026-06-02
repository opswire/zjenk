# mcp-jenkins

MCP server that exposes Jenkins as a set of tools for LLM agents (Claude, etc.).

## Tools

| Tool | Description |
|---|---|
| `jenkins_list_jobs` | List all jobs with name, URL and status |
| `jenkins_get_job` | Get details of a job by name |
| `jenkins_list_builds` | List all builds for a job |
| `jenkins_get_build` | Get build details (result, duration, causes) |
| `jenkins_get_build_log` | Get full console output of a build |
| `jenkins_trigger_build` | Trigger a new build with optional parameters |
| `jenkins_stop_build` | Abort a running build |
| `jenkins_list_nodes` | List all agents/nodes with their status |
| `jenkins_get_queue` | Get the current build queue |
| `jenkins_search_log` | Search build log by regex or substring |

## Requirements

- Go 1.23+
- Jenkins with API access (URL + username + API token)

## Build

```bash
go build -o mcp-jenkins ./cmd/
```

## Configuration

Copy `config.yaml` and fill in your Jenkins credentials:

```yaml
jenkins:
  url: "http://your-jenkins:8080"
  username: "admin"
  password: "your-api-token"   # Jenkins API token, not password

mcp:
  name: "jenkins"
  version: "1.0.0"
```

Config path can be overridden via the `CONFIG_PATH` environment variable.

### Disabling tools

Each tool can be disabled individually in `config.yaml`:

```yaml
tools:
  jenkins_get_build_log:
    is-enabled: false
    name: jenkins_get_build_log
    description:
      ru: "..."
      en: "..."
```

## Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "jenkins": {
      "command": "/path/to/mcp-jenkins",
      "env": {
        "CONFIG_PATH": "/path/to/config.yaml"
      }
    }
  }
}
```

## Integration tests

Tests connect to a real Jenkins instance and skip automatically if `JENKINS_URL` is not set.

```bash
JENKINS_URL=http://localhost:8080 \
JENKINS_USERNAME=admin \
JENKINS_PASSWORD=your-api-token \
go test ./internal/client/ -v
```

| Variable | Required | Description |
|---|---|---|
| `JENKINS_URL` | yes | Jenkins base URL |
| `JENKINS_USERNAME` | yes | Username |
| `JENKINS_PASSWORD` | yes | API token |
| `JENKINS_TEST_JOB` | no | Job name for read tests (auto-discovered if empty) |
| `JENKINS_ALLOW_MUTATIONS` | no | Set to `true` to enable `TriggerBuild` tests |

## Project structure

```
cmd/main.go                  — entry point: fx wiring, tool registrations
config.yaml                  — Jenkins credentials, MCP metadata, tool config
internal/
  config/config.go           — config loading and validation
  client/jenkins.go          — Jenkins API client (gojenkins + resty)
  tool/
    job.go                   — jenkins_list_jobs, jenkins_get_job
    build.go                 — list/get/trigger/stop/log builds
    node.go                  — jenkins_list_nodes, jenkins_get_queue
    search.go                — jenkins_search_log
server/server.go             — MCP server lifecycle
```
