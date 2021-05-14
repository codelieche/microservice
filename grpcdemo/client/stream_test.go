package client

import (
	"context"
	"fmt"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

func getNewsStoreClient() pb.NewsStoreClient {
	addr := "127.0.0.1:9081"
	//conn, err := grpc.Dial(addr) // grpc: no transport security set
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		return nil
	}

	// 2. get news store rpc client
	newsStoreClient := pb.NewNewsStoreClient(conn)

	// 3. 返回
	return newsStoreClient
}

func TestNewsStoreService_GetNewsStream(t *testing.T) {
	// 1. connect
	addr := "127.0.0.1:9081"
	//conn, err := grpc.Dial(addr) // grpc: no transport security set
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		t.Error(err)
		return
	}

	// 2. get news store rpc client
	newsStoreClient := pb.NewNewsStoreClient(conn)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// 3. 开始使用服务器端stream
	r := &pb.NewsRequest{Data: "All"}
	if streamClient, err := newsStoreClient.GetNewsStream(ctx, r); err != nil {
		t.Error(err)
		return
	} else {
		for {
			if item, err := streamClient.Recv(); err != nil {
				if err == io.EOF {
					log.Println("读取完毕了")
				} else {
					t.Error("读取响应结果出错：", err)
				}
				break
			} else {
				log.Println("收到新闻：", item.Data)
			}
		}
		log.Println("Done")
	}
}

func TestNewsStoreService_PutNewsStream(t *testing.T) {
	// 1. get client
	client := getNewsStoreClient()
	if client == nil {
		t.Error("get news store client is nil：")
		return
	}

	// 2. context
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 3. start send news
	if putStream, err := client.PutNewsStream(ctx); err != nil {
		t.Error("get put stream server error: ", err)
		return
	} else {
		// send news
		i := 0
		for {
			i += 1
			if i > 10 {
				putStream.Send(nil)
				break
			}
			item := pb.NewsRequest{
				Data: fmt.Sprintf("Put News %d", i),
			}
			if err := putStream.Send(&item); err != nil {
				t.Error("Put News Error: ", err)
				return
			}
			time.Sleep(time.Second)
		}
		log.Println("Done")
	}
}

func TestNewsStoreService_GetPutNewsStream(t *testing.T) {
	// 1. get client
	client := getNewsStoreClient()
	if client == nil {
		t.Error("Get Client Is nil")
		return
	}

	// 2. context
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 3. start get and push stream
	stream, err := client.GetPutNewsStream(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	// 4. run 2 gorutine
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 4-1: Get Stream
	go func() {
		defer wg.Done()
		for {
			if item, err := stream.Recv(); err != nil {
				log.Println(err)
				break
			} else {
				log.Println("Receive News: ", item.Data)
				if item.Data == "close" {
					log.Println("Please close receive")
					break
				}
			}
		}
		log.Println("Receive gorutine done")
	}()

	// 4-2: Client Push Stream
	go func() {
		defer wg.Done()
		i := 0
		for {
			i += 1
			if i > 10 {
				break
			}
			time.Sleep(time.Second)
			item := pb.NewsRequest{
				Data: fmt.Sprintf("Client Push News Item %d", i),
			}
			if err := stream.Send(&item); err != nil {
				t.Error(err)
				return
			}
		}
		log.Println("Client Push gorutine done")
	}()

	// 5. wait gorutine done
	wg.Wait()
	log.Println("Done")
}
