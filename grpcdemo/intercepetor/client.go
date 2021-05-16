package intercepetor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

/**
客户端的拦截器示例
*/

func RerequestBeforePrint(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	msg := fmt.Sprintf("我现在开始发起请求：%s", method)
	log.Println(msg)
	// 发起请求
	err := invoker(ctx, method, req, reply, cc, opts...)
	// 打印耗时
	log.Println("耗时：", time.Since(start))

	return err
}
