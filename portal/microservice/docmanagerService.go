package microservice

import (
	"context"
	"fmt"
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
	fmt.Println("in get docmanager shard service")
	if r, stored := docmanagerShardServices.Load().(map[uint64]*DocmanagerShardServiceWrapper); stored {
		/*
			fmt.Print("in stored")
			ser, ok := r[1]
			if !ok {
				panic(ok)
			}
			req := &docmanagerModel.SelectAllDocReq{}
			resp, err := ser.Docmanager.SelectAllDoc(context.Background(), req)
			fmt.Println(resp, err)
		*/
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
		fmt.Println(tmp)
	}
	docmanagerShardServices.Store(tg)
	return nil
}

func DocmanagerShardServices(ctx context.Context) (fn model.Daemon, err error) {
	//	fmt.Println("In docmanger service")
	if err = SetDocmanagerShardService(); err != nil {
		return
	}

	//	_ = GetDocmanagerShardServices()

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
