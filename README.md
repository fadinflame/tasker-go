# Tasker

> Simple implementation of Go-kit microservice

## Layout

- Transport
- Endpoint
- Service

## Build

```
go build .
```

## Run

```
go run .
```

## Usage

The service supports gRPC and HTTP protocols

| Name | gRPC | Method | HTTP |
| --- | --- | --- | --- |
| CreateTask | :8082 | POST | :8080/create-task |
| GetTask | :8082 | GET | :8080/task/{id} |
| UpdateTask | :8082 | POST | :8080/task/{id} |
| DeleteTask | :8082 | DELETE | :8080/task/{id} |
