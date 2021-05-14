package interceptors

import (
	"context"
	"github.com/codelieche/microservice/usercenter/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type authInterceptor struct {
	authrizationHeader string
	tokenPrefix        string
	j                  *core.Jwt
}

func NewAuthInterceptor(j *core.Jwt) grpc.UnaryServerInterceptor {
	// 实例化
	i := &authInterceptor{
		authrizationHeader: "authorization",
		tokenPrefix:        "Bearer ",
		j:                  j,
	}

	// 返回一个函数
	return i.HandleRequest
}

func (i *authInterceptor) GetTokenFromContext(ctx context.Context) (string, error) {
	if m, ok := metadata.FromIncomingContext(ctx); !ok {
		return "", status.Error(codes.Unauthenticated, "请传入Token")
	} else {
		// m类型：map[string][]string
		token := ""
		for _, v := range m[i.authrizationHeader] {
			log.Println(v)
			if strings.HasPrefix(v, i.tokenPrefix) {
				token = strings.TrimPrefix(v, i.tokenPrefix)
			}
		}

		if token == "" {
			return "", status.Error(codes.Unauthenticated, "请传入Token")
		} else {
			return token, nil
		}
	}

}

func (i *authInterceptor) HandleRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("收到了请求：", req, info)
	// 1. 从context获取token
	tkn, err := i.GetTokenFromContext(ctx)
	if err != nil {
		if info.FullMethod != "/usercenter.UserService/Login" {
			log.Println("未登录请求:", info.FullMethod)
			//return nil, status.Error(codes.Unauthenticated, "Token过期了")
			return nil, err
		}
	} else {
		// 校验token
		if claims, err := i.j.ParseToken(tkn); err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		} else {
			log.Println("校验token成功：", claims)
		}
	}

	return handler(ctx, req)
}
