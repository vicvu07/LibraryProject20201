package boot

import (
	"context"
	"flag"
	"time"

	"github.com/pinezapple/LibraryProject20201/portal/core"
	"github.com/pinezapple/LibraryProject20201/portal/microservice"
	"github.com/pinezapple/LibraryProject20201/portal/webserver"
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

	webServerPort     = flag.Int("web_port", 11001, "Web Server Endpoints")
	webServerCertFile = flag.String("web_cert", "", "A PEM eoncoded certificate file.")
	webServerKeyFile  = flag.String("web_key", "", "A PEM encoded private key file.")
	clientCAs         = flag.String("client_ca", "", "A CSV list of client's certs trusted by server")

	shardNumber = flag.Int("shard_number", 0, "Docmanager Shards number")
	adm         = flag.Int("adm", 0, "Create admin by default")
)

func Boot() {
	flag.Parse()
	core.InitCore(*shardNumber)

	lg := core.GetLogger()
	logger.LogInfo(lg, "Booting")

	eClient, err := etcd.LoadEtcdClient(*etcdEndpoints, *etcdCertFile, *etcdKeyFile, *etcdCaFile)
	if err != nil {
		logger.LogErr(lg, err)
	}
	core.SetEtcdClient(eClient)

	core.SetHTTPServerConf(&configs.HTTPServerConf{
		Port:      *webServerPort,
		Cert:      *webServerCertFile,
		Key:       *webServerKeyFile,
		ClientCAs: *clientCAs,
	})

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
		webserver.WebServer,
		microservice.DocmanagerShardServices,
	)
}
