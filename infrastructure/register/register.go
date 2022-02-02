package register

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Me1onRind/go-demo/global/client_singleton"
	uuid "github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	gresolver "google.golang.org/grpc/resolver"
)

const (
	prefix = "service"
)

func Register(ctx context.Context, serviceName, addr string) error {
	log.Println("Try register to etcd ...")
	// 创建一个租约
	lease := clientv3.NewLease(client_singleton.EtcdClient)
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

	em, err := endpoints.NewManager(client_singleton.EtcdClient, prefix)
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
	log.Println("Register etcd success")

	del := func() {
		log.Println("Register close")

		cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		_ = em.DeleteEndpoint(cancelCtx, serviceName)

		lease.Close()
	}
	keepRegister(ctx, leaseChannel, del, serviceName, addr)

	return nil
}

func keepRegister(ctx context.Context, leaseChannel <-chan *clientv3.LeaseKeepAliveResponse, cleanFunc func(), serviceName, addr string) {
	go func() {
		failedCount := 0
		for {
			select {
			case resp := <-leaseChannel:
				if resp != nil {
					//log.Println("keep alive success.")
				} else {
					log.Println("keep alive failed.")
					failedCount++
					for failedCount > 3 {
						cleanFunc()
						if err := Register(ctx, serviceName, addr); err != nil {
							time.Sleep(time.Second)
							continue
						}
						return
					}
					continue
				}
			case <-ctx.Done():
				cleanFunc()
				return
			}
		}
	}()
}

func DialTarget(serviceName string) string {
	return fmt.Sprintf("etcd:///%s/%s", prefix, serviceName)
}
func GrpcResolvers() (gresolver.Builder, error) {
	return resolver.NewBuilder(client_singleton.EtcdClient)
}
