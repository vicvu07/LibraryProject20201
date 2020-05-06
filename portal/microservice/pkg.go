package microservice

import (
	"github.com/pinezapple/LibraryProject20201/portal/core"

	etcdnaming "go.etcd.io/etcd/clientv3/naming"
	grpc "google.golang.org/grpc"
)

func getGRPCClientConn(serviceName string) (conn *grpc.ClientConn, err error) {
	balancer := grpc.RoundRobin(&etcdnaming.GRPCResolver{Client: core.GetEtcdClient()})
	conn, err = grpc.Dial(serviceName, grpc.WithInsecure(), grpc.WithBalancer(balancer))
	return
}
