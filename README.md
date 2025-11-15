# **[Тестовое задание](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2025/Backend-trainee-assignment-autumn-2025.md#%D1%82%D0%B5%D1%81%D1%82%D0%BE%D0%B2%D0%BE%D0%B5-%D0%B7%D0%B0%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5-%D0%B4%D0%BB%D1%8F-%D1%81%D1%82%D0%B0%D0%B6%D1%91%D1%80%D0%B0-backend-%D0%BE%D1%81%D0%B5%D0%BD%D0%BD%D1%8F%D1%8F-%D0%B2%D0%BE%D0%BB%D0%BD%D0%B0-2025) для стажировки в Авито**
### **Сервис назначения ревьюеров для Pull Request’ов**
Внутри команды требуется единый микросервис, который автоматически назначает ревьюеров на Pull Request’ы (PR), а также позволяет управлять командами и участниками. Взаимодействие происходит исключительно через HTTP API.
### **Описание задачи**
Необходимо реализовать сервис, который назначает ревьюеров на PR из команды автора, позволяет выполнять переназначение ревьюверов и получать список PR’ов, назначенных конкретному пользователю, а также управлять командами и активностью пользователей. После merge PR изменение состава ревьюверов запрещено.
### Было сделано основное задание: 

1.Используйте этот API (OpenAPI-спецификация будет предоставлена отдельным файлом — openapi.yaml). 


2.Объём данных умеренный (до 20 команд и до 200 пользователей), RPS — 5, SLI времени ответа — 300 мс, SLI успешности — 99.9%. 


3.Пользователь с isActive = false не должен назначаться на ревью. 


4.Операция merge должна быть идемпотентной — повторный вызов не приводит к ошибке и возвращает актуальное состояние PR. 


5.Сервис и его зависимости должны подниматься командой docker-compose up. Если решение предусматривает миграции, они также должны применяться при выполнении этой команды. Сервис должен быть доступен на порту 8080. 


6.Учтите, что соблюдение условий по поднятию сервиса ускорит и упростит проверку вашей работы наставниками. 

**А также были выполнены доп. задания:**

* [Добавить простой эндпоинт статистики (например, количество назначений по пользователям и/или по PR).](./internal/service/stats_service.go)

* [Добавить метод массовой деактивации пользователей команды и безопасную переназначаемость открытых PR (стремиться уложиться в 100 мс для средних объёмов данных).](./internal/service/team_service.go)


* [Реализовать интеграционное или E2E-тестирование.](./tests/integration)


* [Описать конфигурацию линтера.](./.golangci.yaml)


# Инструкция по запуску


**1. Скачать [ZIP-архив](https://github.com/northerf/AvitoTest/archive/refs/heads/main.zip) или клонировать [репозиторий](https://github.com/northerf/AvitoTest.git)**
  
 
**3. Убедиться в наличие у вас установленного [Docker](https://www.docker.com/), Golang**


**4. Ввести в консоль "docker-compose up --build"** --- ** ВАЖНОЕ замечание: проверьте порты перед запуском и освободите 8080 и 5432 **


**5. Можете начинать пользоваться!**


## Как проверить, все ли успешно?

### Если написано: Connected to DB, Starting on port 8080... - тогда все отлично!


### Также можно проверить, зайдя в БД "docker exec -it avitotest-postgres-1 psql -U user -d mydb и введя запрос "SELECT * FROM users;" "

# **Как дергать ручки в моём сервисе?**


## **P.S. Все показанные примеры являются показательными, БД будет загружена !ВОЗМОЖНО! не этими данными**

### Чтобы дергать ручки вам понадобиться Postman или его аналог(я использовал Postman)

### 1. Health Check
**Эндпоинт Health**, необходимый для проверки работоспособности сервиса.


**GET http://localhost:8080/health**


Ответ: 200 ОК
```json 
{
"status": "ok"
}
```
### 2. Team
#### 2.1 **Эндпоинт team/add**, необходимый для создания команды с участниками


**POST http://localhost:8080/team/add**

Headers: Content-Type: application/json 

Body:
```json
{
  "team_name": "frontend",
  "members": [
    {
      "user_id": "u6",
      "username": "Igor",
      "is_active": true
    },
    {
      "user_id": "u7",
      "username": "Egor",
      "is_active": true
    }
  ]
}
```

Ответ: 201 Created

```json
{
    "team": {
        "team_name": "frontend",
        "members": [
            {
                "user_id": "u6",
                "username": "Igor",
                "is_active": true
            },
            {
                "user_id": "u7",
                "username": "Egor",
                "is_active": true
            }
        ]
    }
}
```

**В случае повторном создании команды ответ следующий:**

Ответ: 400 Bad Request

```json
{
    "error": {
        "code": "TEAM_EXISTS",
        "message": "Team already exists"
    }
}
```

#### 2.2 **Эндпоинт team/get**, необходимый для получения команды с участниками

**GET http://localhost:8080/team/get?team_name=backend**

Ответ: 200 ОК

```json
{
    "team_name": "frontend",
    "members": [
        {
            "user_id": "u7",
            "username": "Egor",
            "is_active": true
        },
        {
            "user_id": "u6",
            "username": "Igor",
            "is_active": true
        }
    ]
}
```

**В случае, если команда не найдена ответ следующий:**

Ответ: 404 Not Found

```json
{
    "error": {
        "code": "NOT_FOUND",
        "message": "Team not found"
    }
}
```

#### 2.2 **Эндпоинт team/deactivate**, необходимый для массовой деактивации участников

**POST http://localhost:8080/team/deactivate**

Headers: Content-Type: application/json 

Body:
```json
{
  "team_name": "frontend",
  "user_ids": ["u6", "u7"]
}
```
Ответ: 200 OK

```json
{
    "deactivated_count": 2,
    "reassigned_prs": 0 //в случае, если некого переназначать. В нашем случае, в команде в frontend было 2 человека, мы их деактивировали, соответственно, больше нету => переназначений тоже не будет
}
```

**В случае ошибки ответ будет:**


Ответ: 400 Bad Request


```json
{
    "error": {
        "code": "NOT_FOUND",
        "message": "Invalid request"
    }
}
```



### 3. Users
#### 3.1 **Эндпоинт users/setIsActive**, необходимый, чтобы установить флаг активности пользователя.

**POST http://localhost:8080/users/setIsActive**

Headers: Content-Type: application/json 

Body:
```json
{
  "user_id": "user1",
  "is_active": false
}
```

Ответ: 200 OK
```json
{
    "user": {
        "UserID": "user1",
        "Username": "Alice",
        "TeamName": "backend",
        "IsActive": false
    }
}
```

**В случае неверных данных ответ будет:**

Ответ: 404 Not Found

```json
{
    "error": {
        "code": "NOT_FOUND",
        "message": "User not found"
    }
}
```

#### 3.2 **Эндпоинт users/getReview**, необходимый, чтобы получить PR'ы, где пользователь назначен ревьюером

**GET http://localhost:8080/users/getReview**

Ответ: 200 OK

```json
{
    "pull_requests": [
        {
            "pull_request_id": "pr1",
            "pull_request_name": "Add login feature",
            "author_id": "user1",
            "status": "OPEN"
        }
    ],
    "user_id": "user2"
}
```

### 4. PullRequests
#### 4.1 **Эндпоинт pullRequest/create**, необходимый, чтобы создать PR и автоматически назначить до 2 ревьюеров из команды автора

**POST http://localhost:8080/pullRequest/create**

Headers: Content-Type: application/json 

Body:
```json
{
  "pull_request_id": "pr1111",
  "pull_request_name": "Fix DB",
  "author_id": "user1"
}
```

Ответ: 201 Created

```json
{
    "pr": {
        "pull_request_id": "pr1111",
        "pull_request_name": "Fix DB",
        "author_id": "user1",
        "status": "OPEN",
        "assigned_reviewers": [
            "user3",
            "user4"
        ],
        "createdAt": "2025-11-15T15:37:46.33683Z"
    }
}
```

**В случае повтора данных ответ будет:**

Ответ: 400 Conflict

```json
{
    "error": {
        "code": "PR_EXISTS",
        "message": "PR already exists"
    }
}
```

#### 4.2 **Эндпоинт pullRequest/merge**, необходимый, чтобы пометить PR как MERGED

**POST http://localhost:8080/pullRequest/merge**

Headers: Content-Type: application/json 

Body:


```json
{
  "pull_request_id": "pr1"
}
```

Ответ: 200 OK

```json
{
    "pr": {
        "pull_request_id": "pr1",
        "pull_request_name": "Add login feature",
        "author_id": "user1",
        "status": "MERGED",
        "assigned_reviewers": [
            "user2",
            "user3"
        ],
        "createdAt": "2025-11-15T14:21:58.17311Z",
        "mergedAt": "2025-11-15T15:41:34.200302Z"
    }
}
```

**При повторе запроса:**

Ответ: 200 OK

```json
{ //идемпотента
    "pr": {
        "pull_request_id": "pr1",
        "pull_request_name": "Add login feature",
        "author_id": "user1",
        "status": "MERGED",
        "assigned_reviewers": [
            "user2",
            "user3"
        ],
        "createdAt": "2025-11-15T14:21:58.17311Z",
        "mergedAt": "2025-11-15T15:41:34.200302Z"
    }
}
```

**При неизвестном PR:**

Ответ: 404 Not Found

```json
{
    "error": {
        "code": "NOT_FOUND",
        "message": "PR not found"
    }
}
```

#### 4.3 **Эндпоинт pullRequest/reassign**, необходимый, чтобы переназначить конкретного ревьювера на другого из его команды

**POST http://localhost:8080/pullRequest/reassign**

Headers: Content-Type: application/json 

Body:

```json
{
  "pull_request_id": "pr1",
  "old_reviewer_id": "user3"
}
```

Ответ: 409 Conflict

```json
{ //при переназначении мердженего PR
    "error": {
        "code": "PR_MERGED",
        "message": "Cannot modify merged PR"
    }
}
```

**При нормальном запросе:**


Body:


```json
{
  "pull_request_id": "pr2",
  "old_reviewer_id": "user3"
}
```

Ответ: 200 OK

```json
{
    "pr": {
        "pull_request_id": "pr2",
        "pull_request_name": "Fix bug #123",
        "author_id": "user2",
        "status": "OPEN",
        "assigned_reviewers": [
            "user4",
            "user5"
        ],
        "createdAt": "2025-11-15T14:21:58.17311Z"
    },
    "replaced_by": "user5"
}
```

**При неизвестном PR:**

Ответ: 404 Not Found

```json
{
    "error": {
        "code": "NOT_FOUND",
        "message": "PR not found"
    }
}
```

### 5. stats
#### 5.1 **Эндпоинт stats/allstats**, необходимый, чтобы получить всю статистику


**GET http://localhost:8080/stats/allstats**

Ответ: 200 OK

```json
{
    "total_users": 7,
    "total_active_users": 7,
    "total_prs": 4,
    "total_open_prs": 2,
    "total_merged_prs": 2,
    "top_reviewers": [
        {
            "user_id": "user4",
            "username": "David",
            "reviews_assigned": 3,
            "reviews_completed": 1
        },
        {
            "user_id": "user3",
            "username": "Charlie",
            "reviews_assigned": 2,
            "reviews_completed": 1
        },
        {
            "user_id": "user2",
            "username": "Bob",
            "reviews_assigned": 1,
            "reviews_completed": 1
        },
        {
            "user_id": "user5",
            "username": "Eve",
            "reviews_assigned": 1,
            "reviews_completed": 0
        }
    ],
    "prs_without_reviewers": 0
}
```

#### 5.2 **Эндпоинт stats/user?user_id=<ID пользователя>**, необходимый, чтобы получить статистику по пользователю


**GET http://localhost:8080/stats/user?user_id=<ID пользователя>**

Ответ: 200 OK

```json
{
    "user_id": "user2",
    "username": "Bob",
    "reviews_assigned": 1,
    "reviews_completed": 1
}
```

## Линтер

**Для начала вам необходимо** 

    * 1. установить с помощью команды в терминале: "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

    
    * 2. Ввести в терминал команду: "export PATH=$PATH:$(go env GOPATH)/bin"

    
    * 3. Запустить линтер


### Чтобы проверить его работоспособность необходимо ввести в терминал команду "golangci-lint run", "golangci-lint run --fix"

## Тесты

### Чтобы запустить тесты необходимо перейти в директорию /tests/integration и в терминале ввести команду "go test -v"
