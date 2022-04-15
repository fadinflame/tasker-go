package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tasker/pb"
	task_ep "tasker/pkg/endpoint"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

func NewHttpServer(tm task_ep.Endpoints, options []kithttp.ServerOption, logger log.Logger) http.Handler {
	var (
		r            = mux.NewRouter()
		errorLogger  = kithttp.ServerErrorLogger(logger)
		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
	)

	options = append(options, errorLogger, errorEncoder)
	r.Methods("POST").Path("/create-task").Handler(kithttp.NewServer(
		tm.CreateTask,
		decodeTaskCreateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/task/{id}").Handler(kithttp.NewServer(
		tm.GetTask,
		decodeGetTaskRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/task/{id}").Handler(kithttp.NewServer(
		tm.UpdateTask,
		decodeUpdateTaskRequest,
		encodeResponse,
		options...,
	))

	r.Methods("DELETE").Path("/task/{id}").Handler(kithttp.NewServer(
		tm.DeleteTask,
		decodeDeleteTaskRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeTaskCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req pb.TaskCreateRequest
	er := json.NewDecoder(r.Body).Decode(&req)
	return &req, er
}

func decodeGetTaskRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := getIdFromRequest(r)
	if err != nil {
		return nil, ErrBadRouting
	}
	return &pb.TaskGetRequest{TaskId: id}, nil
}

func decodeUpdateTaskRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req pb.TaskUpdateRequest
	id, err := getIdFromRequest(r)
	if err != nil {
		return nil, ErrBadRouting
	}
	req.TaskId = id
	er := json.NewDecoder(r.Body).Decode(&req)
	return &req, er
}

func decodeDeleteTaskRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req pb.TaskDeleteRequest
	id, err := getIdFromRequest(r)
	if err != nil {
		return nil, ErrBadRouting
	}
	req.TaskId = id
	return &req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic(any("shit"))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 0, 64)
	if err != nil {
		return 0, ErrBadRouting
	}
	return id, nil
}
