# Sparrow API

## Functionality
- `GET /health` — returns JSON with ok status, current time and version of the project
  ```json
  {
    "status": 200,
    "time": "2025-10-01T12:56:19Z",
    "status": "0.1.0" 
  }
  ```

## Project structure
```text
sparrow-api
├── cmd
│   └── server
│       └── main.go
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── lab1
│   │   ├── assets
│   │   └── README.md
│   ├── lab2
│   │   ├── assets
│   │   └── README.md
│   ├── lab3
│   │   ├── assets
│   │   └── README.md
│   └── lab4
│       ├── assets
│       └── README.md
├── go.mod
├── go.sum
├── internal
│   └── ...
└── README.md
```

## How to start localy
```bash
go run ./cmd/server
```

## Run with Docker
```bash
# build image
docker build -t sparrow-api:latest .

# run container
docker run --rm -p 8080:8080 -e PORT=8080 sparrow-api:latest
```

## Run with Docker Compose
```bash
docker-compose up --build
```

## Deploy
The application is deployed on Render.

Accessible at: https://sparrow-api-l8pp.onrender.com

## Git

Using Conventional Commits style.

Examples:
```text
feat: add chi router with /health endpoint
chore(docker): add Dockerfile and docker-compose setup
```

## Author
- GitHub: [@Pliffdax](https://github.com/Pliffdax)  
- Telegram: [@Pliffdax](https://t.me/Pliffdax)
