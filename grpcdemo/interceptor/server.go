package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

/**
服务端的拦截器示例
*/

func RequestInfoPrint(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 1. 开始执行请求前相关操作
	// 1-1：记录消息
	start := time.Now()
	msg := fmt.Sprintf("请求:%s", info.FullMethod)
	log.Println(msg)

	// 1-2：打印metadata
	// 从传入的context中获取metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("请求Header:")
		// 遍历metadata，并输出key和value
		for k, v := range md {
			msg := fmt.Sprintf("\t\t\t===> %s: \t %s", k, v)
			fmt.Println(msg)
		}
	}

	// 2. 调用handler处理请求
	resp, err = handler(ctx, req)
	log.Println(msg, "耗时:", time.Since(start))

	// 3. 返回请求响应的结果
	return resp, err
}
