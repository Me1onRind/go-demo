package initialize

//import (
//"github.com/Me1onRind/go-demo/global/store"
//"github.com/Me1onRind/go-demo/global/store/db_label"
//"github.com/Me1onRind/go-demo/internal/core/config"
//"github.com/Me1onRind/go-demo/internal/lib/client/asynq_client"
//"github.com/Me1onRind/go-demo/internal/lib/client/etcd_client"
//"github.com/Me1onRind/go-demo/internal/lib/client/grpc_client"
//"github.com/Me1onRind/go-demo/internal/lib/client/mysql_client"
//"github.com/Me1onRind/go-demo/internal/lib/client/redis_client"
//"github.com/Me1onRind/go-demo/internal/lib/register"
//)

//func InitGrpcClients() error {
//resolver, err := register.GrpcResolvers()
//if err != nil {
//return err
//}
//if err := grpc_client.InitGoDemoClient(register.DialTarget("go-demo"), resolver); err != nil {
//return err
//}
//return nil
//}

//func CloseGrpcClients() error {
//return grpc_client.CloseGoDemoClient()
//}

//func InitAsynqClient() error {
//asynqConfig := &config.RemoteConfig.Asynq
//asynq_client.InitAsynqClient(asynqConfig)
//return nil
//}

//func CloseAsynqClient() error {
//return asynq_client.AsynqClient.Close()
//}

//func InitEtcdClient() error {
//etcdConfig := &config.LocalConfig.Etcd
//return etcd_client.InitEtcdClient(etcdConfig)
//}

//func CloseEtcdClient() error {
//return etcd_client.EtcdClient.Close()
//}

//func InitMysqlClients() error {
//var err error
//dbs := config.RemoteConfig.DBs
////store.DBs[db_label.DB], err = mysql_client.NewDBClient(&dbs.DB)
////if err != nil {
////return err
////}
//store.DBs[db_label.ConfigDB], err = mysql_client.NewDBClient(&dbs.ConfigDB)
//if err != nil {
//return err
//}
//return nil
//}

//func CloseMysqlClients() error {
//for _, v := range store.DBs {
//db, err := v.DB()
//if err != nil {
//return err
//}
//if err := db.Close(); err != nil {
//return err
//}
//}
//return nil
//}

//func InitRedisClient() error {
//redisConfig := &config.RemoteConfig.Redis
//store.RedisClient = redis_client.NewRedisClient(redisConfig)
//return nil
//}

//func CloseRedisClient() error {
//return store.RedisClient.Close()
//}
