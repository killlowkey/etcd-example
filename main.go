package main

import (
	"context"
	"flag"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

var (
	addr = flag.String("addr", "localhost:2379", "-addr=localhost:2379")
	cli  *clientv3.Client
	key  = "sample_key"
)

func init() {
	flag.Parse()
	connectEtcd()
}

func connectEtcd() {
	var err error
	// 创建 etcd v3 客户端
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{*addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		_ = cli.Close()
	}()

	//watchExample()
	//txExample()
	//delExample()
	contextExample()
}

// insertExample etcd 数据插入
func insertExample() {
	// 插入数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Put(ctx, key, "sample_value")
	log.Println("insert: ", resp)
	cancel()
	if err != nil {
		panic(err)
	}
}

// getExample etcd 数据获取
func getExample() {
	// 获取数据
	res, err := cli.Get(context.Background(), key)
	if err != nil {
		panic(err)
	}
	log.Println(res.Kvs)
}

// delExample etcd 数据删除
func delExample() {
	res, err := cli.Put(context.Background(), "del-key", "1")
	if err != nil {
		panic(res)
	}
	log.Printf("insert %v data successfully\n", res)

	delRes, err := cli.Delete(context.Background(), "del-key")
	if err != nil {
		panic(err)
	}
	log.Printf("delete data：%v\n", delRes)
}

// contextExample 使用上下文来控制 GRPC 通信
func contextExample() {
	// time.Microsecond 必然超时
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancelFunc()

	res, err := cli.Put(ctx, "context-example", "context")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}

// txExample 事务例子
// https://juejin.cn/post/7089804562818662413
func txExample() {
	// 先插入一些测试数据
	_, _ = cli.Put(context.Background(), "k1", "10")
	_, _ = cli.Put(context.Background(), "k2", "20")

	// 开启一个新的事务
	txn, err := cli.Txn(context.Background()).If(
		// 判断是否需要执行事务
		clientv3.Compare(clientv3.Value("k1"), "=", "10"),
		clientv3.Compare(clientv3.Value("k2"), "=", "20"),
	).Then(
		// If 通过之后，则执行
		clientv3.OpPut("txn", "success"),
	).Else(
		// If 不通过后，则执行
		clientv3.OpPut("txn", "error"),
	).Commit()
	if err != err {
		panic(err)
	}

	log.Println("Transaction：", txn.Succeeded)
}

// watchExample etcd watch 机制
// 一般用于实现注册中心和配置中心，key 被修改后，可以动态的感知
func watchExample() {
	var wg sync.WaitGroup
	wg.Add(2)

	// 先插入数据，保存 etcd 是存在这个数据的
	insertExample()

	// watch key
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		watchChan := cli.Watch(ctx, key)
		select {
		case resp := <-watchChan:
			for _, event := range resp.Events {
				log.Printf("event type：%d, key[%s] value[%s]\n",
					event.Type,
					event.Kv.Key,
					event.Kv.Value,
				)
			}
		}
		wg.Done()
	}()

	go func() {
		resp, err := cli.Put(context.Background(), key, "ray")
		if err != nil {
			panic(err)
		}

		log.Printf("change %s value to ray, resp: %+v\n", key, resp)
		wg.Done()
	}()

	wg.Wait()
}
