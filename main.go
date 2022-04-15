package main

import (
	"database/sql"
	"flag"
	"fmt"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tasker/pb"
	"tasker/pkg/consts"
	"tasker/pkg/endpoint"
	"tasker/pkg/service"
	"tasker/pkg/transport"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {

	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("gRPC.addr", ":8082", "gRPC listen address")

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	db, err := sql.Open("sqlite3", consts.DatabasePath)
	if err != nil {
		level.Error(logger)
	}

	ts := service.NewTaskService(db, logger)
	if err := ts.Migrate(); err != nil {
		level.Error(logger)
	}

	endpoints := endpoint.MakeEndpoints(ts)
	httpServer := transport.NewHttpServer(endpoints, []kithttp.ServerOption{}, logger)
	gRPCServer := transport.NewGRPCServer(endpoints, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: httpServer,
		}
		errs <- server.ListenAndServe()
	}()

	go func() {
		level.Info(logger).Log("transport", "GRPC", "addr", *grpcAddr)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:8082"))
		if err != nil {
			level.Error(logger).Log("failed to listen: %v", err)
		}

		baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		pb.RegisterTaskServiceServer(baseServer, gRPCServer)
		errs <- baseServer.Serve(lis)
	}()

	level.Error(logger).Log("exit", <-errs)
}
