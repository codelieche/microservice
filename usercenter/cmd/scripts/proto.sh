#!/bin/bash


cd ../../proto/

if [[ $PWD =~ "proto" ]]
then
  ls -alh
  files=`ls -alh`
  if ! [[ $files =~ ".proto" ]]
  then
    echo "`date '+%F %T'`: 未找到proto文件，程序退出"; exit 1;
  else
    echo "`date '+%F %T'`: 当前目录：$PWD"
  fi
else
  echo "`date '+%F %T'`: 请进入scripts目录执行脚本"
  exit 1
fi


echo "`date '+%F %T'`: 开始执行protoc命令"
# 生成grpc proto相关代码
protoc -I=. --go_out=paths=source_relative:userpb \
     --go-grpc_out=require_unimplemented_servers=false:userpb --go-grpc_opt=paths=source_relative \
    ./user.proto

# 生产grpc gateway相关代码
protoc -I=. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=./yaml/user.yaml:userpb ./user.proto
