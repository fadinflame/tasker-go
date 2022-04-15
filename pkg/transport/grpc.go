package transport

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"tasker/pb"
	task_ep "tasker/pkg/endpoint"
)

type grpcServer struct {
	createTask grpctransport.Handler
	getTask    grpctransport.Handler
	updateTask grpctransport.Handler
	deleteTask grpctransport.Handler
	pb.UnimplementedTaskServiceServer
	logger log.Logger
}

func NewGRPCServer(tm task_ep.Endpoints, logger log.Logger) pb.TaskServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		createTask: grpctransport.NewServer(
			tm.CreateTask,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
		getTask: grpctransport.NewServer(
			tm.GetTask,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
		updateTask: grpctransport.NewServer(
			tm.UpdateTask,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
		deleteTask: grpctransport.NewServer(
			tm.DeleteTask,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
		logger: logger,
	}
}

func (gs *grpcServer) CreateTask(ctx context.Context, request *pb.TaskCreateRequest) (*pb.TaskCreateResponse, error) {
	level.Debug(gs.logger).Log("gRPC", "CreateTask")
	_, resp, err := gs.createTask.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TaskCreateResponse), nil
}

func (gs *grpcServer) GetTask(ctx context.Context, request *pb.TaskGetRequest) (*pb.TaskGetResponse, error) {
	level.Debug(gs.logger).Log("gRPC", "GetTask")
	_, resp, err := gs.getTask.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TaskGetResponse), nil
}

func (gs *grpcServer) UpdateTask(ctx context.Context, request *pb.TaskUpdateRequest) (*pb.TaskUpdateResponse, error) {
	level.Debug(gs.logger).Log("gRPC", "UpdateTask")
	_, resp, err := gs.updateTask.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TaskUpdateResponse), nil
}

func (gs *grpcServer) DeleteTask(ctx context.Context, request *pb.TaskDeleteRequest) (*pb.TaskDeleteResponse, error) {
	level.Debug(gs.logger).Log("gRPC", "DeleteTask")
	_, resp, err := gs.deleteTask.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TaskDeleteResponse), nil
}

func decodeGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func encodeGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
