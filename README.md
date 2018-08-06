# QA API Server

Write API integration tests for the following functionality

## Test Cases
1. Listing tasks initially returns an empty array
2. Creating a task fails if any of the required args are missing (title, description)
3. The newly created task is returned by the `list_tasks` api
4. After deleting a task it is not returned by `list_tasks`
5. Deleting an unknown task id errors

## APIs

The API has a common response structure.
If the operation is a success, the status will be `success`, and the `data` field may be populated.
If an error occurs the status will be `error` and the `error` field will be populated with a message.

Depending on the error, the API may return HTTP status code 400 or 500.

```json
    {
        "status":"success|error",
        "data":<object>,
        "error":"string"
    }
```

The task object has the following attributes

- *id:int* populated after the task is created
- *title:string* Required
- *description:string* Required

### GET /api/list_tasks

```bash
$ curl http://localhost:8080/api/list_tasks
{"status":"success","data":[{"id":10,"title":"foo 2","description":"foo"}]}
```

### POST /api/create_task

Successful create
```bash
$ curl http://localhost:8080/api/create_task -d '{"title":"My Task", "description":"I have some work to do"}'
{"status":"success","data":{"id":10,"title":"My Task","description":"I have some work to do"}}
```

Missing required arg
```bash
$ curl http://localhost:8080/api/create_task -d '{"title":"My Task"}'
{"status":"error","error":"A description must be provided"}
```

### POST /api/delete_task

Success deleting a task
```bash
$ curl http://localhost:8080/api/delete_task -d '{"id":10}'
{"status":"success"}
```

Error deleting an unknown task id
```bash
$ curl http://localhost:8080/api/delete_task -d '{"id":2}'
{"status":"error","error":"The item with id 2 does not exist"}
```
