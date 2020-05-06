package service

import (
	"context"
	"fmt"

	"github.com/pinezapple/LibraryProject20201/docmanager/core"
	"github.com/pinezapple/LibraryProject20201/skeleton/booting"
	"github.com/pinezapple/LibraryProject20201/skeleton/logger"
	"github.com/pinezapple/LibraryProject20201/skeleton/model"
	"github.com/pinezapple/LibraryProject20201/skeleton/model/docmanagerModel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// GRPCServer booting grpc server by configuration
func GRPCServer(ctx context.Context) (daemon model.Daemon, err error) {
	lg := core.GetLogger()
	logger.LogInfo(lg, ".grpc")

	// get configuration
	conf := core.GetGRPCServerConfig()
	if conf == nil {
		err = fmt.Errorf("gRPC Server configuration not initialized")
		return
	}

	return booting.GRPCService(ctx,
		core.ServiceName,
		core.GetEtcdClient(),
		*conf,
		func(s *grpc.Server) {
			docmanagerModel.RegisterDocmanagerServer(s, &docmanagerServer{lg: lg})
		},
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: conf.MaxConnectionIdle}),
	)
}
