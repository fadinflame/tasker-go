package endpoint

import (
	"context"
	"tasker/pb"
	models "tasker/pkg/models"
	ts "tasker/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
)

type Endpoints struct {
	CreateTask endpoint.Endpoint
	GetTask    endpoint.Endpoint
	UpdateTask endpoint.Endpoint
	DeleteTask endpoint.Endpoint
}

func MakeEndpoints(ts *ts.TaskService) Endpoints {
	return Endpoints{
		CreateTask: CreateTaskEndpoint(ts),
		GetTask:    GetTaskEndpoint(ts),
		UpdateTask: UpdateTaskEndpoint(ts),
		DeleteTask: DeleteTaskEndpoint(ts),
	}
}

func CreateTaskEndpoint(ts *ts.TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		level.Debug(ts.Logger).Log("endpoint", "CreateTask")
		req := request.(*pb.TaskCreateRequest)
		task, err := ts.CreateTask(models.Task{Title: req.Title, Text: req.Text})
		return &pb.TaskCreateResponse{TaskId: task.Id, Success: err == nil}, err
	}
}

func GetTaskEndpoint(ts *ts.TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		level.Debug(ts.Logger).Log("endpoint", "GetTask")
		req := request.(*pb.TaskGetRequest)
		task, err := ts.GetTask(req.TaskId)
		return &pb.TaskGetResponse{
			TaskId:      task.Id,
			Title:       task.Title,
			Text:        task.Text,
			IsCompleted: task.IsCompleted,
		}, err
	}
}

func UpdateTaskEndpoint(ts *ts.TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		level.Debug(ts.Logger).Log("endpoint", "UpdateTask")
		req := request.(*pb.TaskUpdateRequest)
		updated, err := ts.UpdateTask(req.TaskId, models.Task{Title: req.Title, Text: req.Text, IsCompleted: req.IsCompleted})
		return &pb.TaskUpdateResponse{TaskId: req.TaskId, Success: updated}, err
	}
}

func DeleteTaskEndpoint(ts *ts.TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		level.Debug(ts.Logger).Log("endpoint", "DeleteTask")
		req := request.(*pb.TaskDeleteRequest)
		deleted, err := ts.DeleteTask(req.TaskId)
		return &pb.TaskDeleteResponse{Success: deleted}, err
	}
}
