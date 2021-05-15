package server

import (
	"context"
	"fmt"
	"github.com/codelieche/microservice/grpcdemo/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"sync"
	"time"
)

type NewsStoreService struct {
	logger *zap.Logger // 日志
}

func NewNewsStoreService(logger *zap.Logger) pb.NewsStoreServer {
	return &NewsStoreService{logger: logger}
}

func (s *NewsStoreService) Ping(ctx context.Context, empty *emptypb.Empty) (*pb.Pong, error) {
	pong := &pb.Pong{
		Status:  true,
		Message: "pong",
	}
	return pong, nil
}

func (s *NewsStoreService) GetNewsStream(request *pb.NewsRequest, server pb.NewsStore_GetNewsStreamServer) error {
	// 发送十条新闻
	index := 0
	category := request.Data

	for {
		index += 1
		if index > 10 {
			break
		}
		// 发送新闻
		response := pb.NewsResponse{
			Data: fmt.Sprintf("%s: News Item %d", category, index),
		}
		if err := server.Send(&response); err != nil {
			s.logger.Fatal("send news error", zap.Error(err))
			return err
		}
		// 睡眠1秒
		time.Sleep(time.Second)
	}
	return nil
}

func (s *NewsStoreService) PutNewsStream(server pb.NewsStore_PutNewsStreamServer) error {
	// 从客户端推流中不断获取请求
	for {
		if item, err := server.Recv(); err != nil {
			if err == io.EOF {
				s.logger.Info("从客户端推流中读取完毕")
				break
			} else {
				s.logger.Error("Receive Client Stream Error", zap.Error(err))
				return err
			}
		} else {
			s.logger.Info("Receive News Request", zap.String("news", item.Data))
		}
	}
	s.logger.Info("Receive Done")
	return nil
}

func (s *NewsStoreService) GetPutNewsStream(server pb.NewsStore_GetPutNewsStreamServer) error {
	// 1. 并发处理服务端和客户端的推流
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 2. 处理客户端推流
	go func() {
		defer wg.Done()
		for {
			if item, err := server.Recv(); err != nil {
				s.logger.Info("Receive Item Error", zap.Error(err))
				break
			} else {
				s.logger.Info("Receive News", zap.String("news", item.Data))

				if item.Data == "close" {
					log.Println("Please close receive")
					break
				}
			}
		}
		log.Println("Receive gorutine done")
	}()

	// 3. 处理服务端推流
	go func() {
		defer wg.Done()
		i := 0
		for {
			i += 1
			if i > 10 {
				server.Send(&pb.NewsResponse{Data: "close"})
				break
			}
			// 推送新闻
			item := pb.NewsResponse{
				Data: fmt.Sprintf("News Item %d", i),
			}
			if err := server.Send(&item); err != nil {
				s.logger.Error("Send News Error", zap.Error(err))
				break
			}
		}
		log.Println("Server push gorutine done")

	}()

	// 4. 等待2个协程的结束
	wg.Wait()
	s.logger.Info("GetPutStream Done")
	return nil
}
