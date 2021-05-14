#!/bin/bash

# 生成的代码存放目录
TARGET_DIR=pb

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

# 处理函数
function generate() {
    # 1. 提取变量
    FILE=$1
    GATEWAY=$2
    echo "`date '+%F %T'`: 开始处理${FILE}.proto: (${FILE}, gateway ${GATEWAY})"

    # 2. 生成grpc proto相关代码
    protoc -I=. --go_out=paths=source_relative:${TARGET_DIR} \
         --go-grpc_out=require_unimplemented_servers=false:${TARGET_DIR} --go-grpc_opt=paths=source_relative \
        ./${FILE}.proto

    # 3. 生产grpc gateway相关代码
    if [[ $GATEWAY == "true" ]];then
        echo "`date '+%F %T'`: 开始处理${FILE} gateway"
        protoc -I=. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=./yaml/${FILE}.yaml:${TARGET_DIR} ./${FILE}.proto
    fi

}

# 处理hello.proto文件
generate hello true
