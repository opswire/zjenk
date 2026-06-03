# mcp-jenkins

MCP-сервер, который предоставляет Jenkins как набор инструментов для LLM-агентов (Claude и др.).

## Инструменты

| Инструмент              | Описание                                                   |
|-------------------------|------------------------------------------------------------|
| `jenkins_list_jobs`     | Список всех джоб с именем, URL и статусом                  |
| `jenkins_get_job`       | Получить детали конкретной джобы                           |
| `jenkins_list_builds`   | Список всех сборок джобы                                   |
| `jenkins_get_build`     | Детали сборки (результат, длительность, причины)           |
| `jenkins_get_build_log` | Полный вывод консоли сборки                                |
| `jenkins_trigger_build` | Запустить новую сборку с опциональными параметрами         |
| `jenkins_stop_build`    | Прервать выполняющуюся сборку                              |
| `jenkins_list_nodes`    | Список всех агентов/нод с их статусом                      |
| `jenkins_get_queue`     | Текущая очередь сборок                                     |
| `jenkins_search_log`    | Поиск в логе сборки по регулярному выражению или подстроке |

### Параметр job_path

Инструменты `get_job`, `list_builds`, `get_build`, `get_build_log`, `trigger_build`, `stop_build`, `search_log` принимают `job_path` — путь к джобе **без** префикса `/job/`.

| Тип джобы       | Пример job_path           |
|-----------------|---------------------------|
| Простая джоба   | `my-job`                  |
| Джоба в папке   | `folder/my-job`           |
| Вложенные папки | `folder/subfolder/my-job` |

Сервер автоматически преобразует путь в формат Jenkins REST API: `folder/my-job` → `/job/folder/job/my-job`.

### Параметр project_name

Инструменты `list_jobs`, `list_nodes`, `get_queue` принимают `project_name` — имя проекта, которое подставляется как поддомен Jenkins URL.

Например, при базовом URL `https://jenkins.domain.com` и `project_name: my-project` запрос уйдёт на `https://my-project.jenkins.domain.com`.

Путь к конфигу можно переопределить через переменную окружения `CONFIG_PATH`.

## Интеграционные тесты

Тесты подключаются к реальному Jenkins и пропускаются автоматически, если не задан `JENKINS_URL`. Тестовые джоба и проект задаются константами `testJob` и `testProject` в файле теста.

```bash
JENKINS_URL=http://localhost:8080 \
JENKINS_USERNAME=admin \
JENKINS_PASSWORD=your-api-token \
go test ./internal/client/ -v
```

| Переменная                | Обязательная | Описание                                              |
|---------------------------|--------------|-------------------------------------------------------|
| `JENKINS_URL`             | да           | Базовый URL Jenkins                                   |
| `JENKINS_USERNAME`        | да           | Имя пользователя                                      |
| `JENKINS_PASSWORD`        | да           | API-токен                                             |
| `JENKINS_ALLOW_MUTATIONS` | нет          | Установите `true` чтобы включить тест `TriggerBuild`  |
