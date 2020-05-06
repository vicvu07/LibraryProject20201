package boot

import (
	"context"
	"flag"
	"time"

	"github.com/pinezapple/LibraryProject20201/docmanager/core"
	grpcService "github.com/pinezapple/LibraryProject20201/docmanager/service"
	"github.com/pinezapple/LibraryProject20201/skeleton/booting"
	"github.com/pinezapple/LibraryProject20201/skeleton/configs"
	etcd "github.com/pinezapple/LibraryProject20201/skeleton/libs/etcd"
	"github.com/pinezapple/LibraryProject20201/skeleton/logger"
)

var (
	etcdEndpoints = flag.String("etcd_endpoints", "http://localhost:2379", "ETCD Endpoints")
	etcdCertFile  = flag.String("etcd_cert", "", "A PEM eoncoded certificate file.")
	etcdKeyFile   = flag.String("etcd_key", "", "A PEM encoded private key file.")
	etcdCaFile    = flag.String("etcd_ca", "", "A PEM eoncoded CA's certificate file.")

	grpcPublicIP     = flag.String("grpc_public_ip", "", "Public IP of gRPC Server")
	grpcPort         = flag.Int("grpc_port", 10000, "gRPC Server Endpoints")
	grpcCertFile     = flag.String("grpc_cert", "", "A PEM eoncoded certificate file.")
	grpcKeyFile      = flag.String("grpc_key", "", "A PEM encoded private key file.")
	grpcClientCaFile = flag.String("grpc_client_ca", "", "A PEM eoncoded client CA's file.")
	shardID          = flag.Int("shard", 0, "Shard's number'")
)

func Boot() {
	flag.Parse()
	core.InitCore(*shardID)

	lg := core.GetLogger()
	logger.LogInfo(lg, "Booting")

	eClient, err := etcd.LoadEtcdClient(*etcdEndpoints, *etcdCertFile, *etcdKeyFile, *etcdCaFile)
	if err != nil {
		logger.LogErr(lg, err)
	}
	core.SetEtcdClient(eClient)

	// set configurations
	core.SetGRPCServerConfig(&configs.GRPCServerConfig{
		PublicIP:          *grpcPublicIP,
		Port:              *grpcPort,
		Cert:              *grpcCertFile,
		Key:               *grpcKeyFile,
		ClientCA:          *grpcClientCaFile,
		MaxConnectionIdle: time.Hour,
	})
	// get config from etcd
	if err = core.GetConfigFromEtcd(); err != nil {
		logger.LogErr(lg, err)
		panic("failed to get config from etcd")
	}

	booting.BootstrapDaemons(context.Background(),
		core.WatchConfig,
		grpcService.GRPCServer,
	)

}
