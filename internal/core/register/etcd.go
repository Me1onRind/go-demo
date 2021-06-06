package register

import (
	"context"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

var client *clientv3.Client

const (
	prefix = "service"
)

func init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
}

func Register(ctx context.Context, serviceName, addr string) error {
	// 创建一个租约
	lease := clientv3.NewLease(client)
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	leaseResp, err := lease.Grant(cancelCtx, 3)
	if err != nil {
		return err
	}

	leaseChannel, err := lease.KeepAlive(ctx, leaseResp.ID) // 长链接, 不用设置超时时间
	if err != nil {
		return err
	}

	em, err := endpoints.NewManager(client, prefix)
	if err != nil {
		return err
	}

	cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	if err := em.AddEndpoint(cancelCtx, fmt.Sprintf("%s/%s/%s", prefix, serviceName, uuid.NewV4().String()), endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case resp := <-leaseChannel:
				if resp != nil {
					//log.Println("keep alive success.")
				} else {
					log.Println("keep alive failed.")
					time.Sleep(time.Second)
					continue
				}
			case <-ctx.Done():
				log.Println("close service register")

				cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
				defer cancel()
				em.DeleteEndpoint(cancelCtx, serviceName)

				lease.Close()
				client.Close()
				return
			}
		}
	}()

	return nil
}
