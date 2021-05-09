
## protoc基本使用


```bash
# 1. 进入当前目录
cd proto  
# 2. 执行protoc
protoc -I=. --go_out=paths=source_relative:userpb ./user.proto
protoc -I=. --go_out=paths=source_relative:userpb \
     --go-grpc_out=userpb --go-grpc_opt=paths=source_relative \
    ./user.proto
```

- `--go_out=paths=source_relative:gengo`: `source_relative`是使用相对路径，冒号后的是go要生成的文件存放路径
- `--go-grpc_out=userpb --go-grpc_opt=paths=source_relative`: 生成grpc相关的代码

