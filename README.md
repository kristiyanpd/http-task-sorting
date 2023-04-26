# HTTP Task Sorting

## Overview
The application takes a collection of tasks, where each task has a name and a shell command. Tasks may depend on other tasks and require that those are executed beforehand. The service takes care of sorting the tasks to create a proper execution order.

# REST API

## Get sorted tasks

### Request

`POST /tasks/`

```shell
curl -d @examples/tasks.json http://localhost:8000/tasks 
```

### Response

```json
[
    {
        "name": "task-1",
        "command": "touch ./file1"
    },
    {
        "name": "task-3",
        "command": "echo 'Hello World!' > ./file1"
    },
    {
        "name": "task-2",
        "command": "cat ./file1"
    },
    {
        "name": "task-4",
        "command": "rm ./file1"
    }
]
```

## Get sorted tasks directly as bash command

### Request

`POST /tasks/`

```shell
curl -H "Accept: application/bash" -d @examples/tasks.json http://localhost:8000/tasks | bash
```

### Response

```shell
#!/usr/bin/env bash
touch ./file1
echo 'Hello World!' > ./file1
cat ./file1
rm ./file1
```
