#!/usr/bin/env bash

# 启动单独的exporter
# use case  ./start_daemon.sh nm_exporter
exporter_nm=$1

if [ -n "$exporter_nm" ];then
  nohup ./$exporter_nm 2>&1 &
  tail -f nohup.out
else
  echo '输入正确到exporter名称'
fi


