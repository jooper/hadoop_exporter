#!/usr/bin/env bash

# 停止单独的exporter
# use case  ./stop_daemon.sh nm_exporter


# 进程名

#process="nm_exporter"
process=$1

if [ -n "$process" ];then
  ps -ef | grep $process
  # 获取进程ID,清除启动时候残留进程
  PID2=$(ps -ef | grep $process | grep -v grep | awk '{print $2}')

  if [ -n "$PID2" ];then
    kill -9 $PID2
    echo 'exporter 已停止,processId:'${PID2}
    exit 1
  else
    echo '服务未启动'
  fi
else
    echo '输入正确到exporter名称'
fi



# 获取进程ID
#PID=$(netstat -nlp|grep 9070 |awk '{print $7}'|awk -F/ '{print $1}')
#
#if [ -n "$PID" ] ; then
#  netstat -nlp|grep 9070 |awk '{print $7}'|awk -F/ '{print $1}'|xargs kill -9
#  echo 'exporter 已停止,processId:'${PID}
#  exit 1
#else
#  echo '服务未启动'
#fi


