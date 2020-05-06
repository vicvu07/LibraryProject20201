package microservice

import (
	"context"
	"strconv"
	"sync/atomic"

	"github.com/pinezapple/LibraryProject20201/portal/core"
	"github.com/pinezapple/LibraryProject20201/skeleton/model"
	"github.com/pinezapple/LibraryProject20201/skeleton/model/docmanagerModel"
	grpc "google.golang.org/grpc"
)

type DocmanagerShardServiceWrapper struct {
	Conn       *grpc.ClientConn
	Docmanager docmanagerModel.DocmanagerClient
}

var docmanagerShardServices atomic.Value

func GetDocmanagerShardServices() map[uint64]*DocmanagerShardServiceWrapper {
	if r, stored := docmanagerShardServices.Load().(map[uint64]*DocmanagerShardServiceWrapper); stored {
		return r
	}
	return nil
}

func SetDocmanagerShardService() error {
	tg := make(map[uint64]*DocmanagerShardServiceWrapper)
	numShards := core.GetNumShards()
	// added
	var err error
	var v uint64
	for v = 0; v < numShards; v++ {
		tmp := &DocmanagerShardServiceWrapper{}
		tmp.Conn, err = getGRPCClientConn("docmanager" + strconv.Itoa(int(v)))
		if err != nil {
			return err
		}
		tmp.Docmanager = docmanagerModel.NewDocmanagerClient(tmp.Conn)

		tg[v] = tmp
	}

	return nil
}

func DocmanagerShardServices(ctx context.Context) (fn model.Daemon, err error) {
	if err = SetDocmanagerShardService(); err != nil {
		return
	}

	fn = func() {
		<-ctx.Done()
		tg := GetDocmanagerShardServices()
		for _, conn := range tg {
			if conn != nil && conn.Conn != nil {
				_ = conn.Conn.Close()
			}
		}
	}

	return
}
