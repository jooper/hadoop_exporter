# Hadoop Exporter for Prometheus
Exports hadoop metrics via HTTP for Prometheus consumption.

How to build
```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/log

set GOARCH=amd64  设置目标可执行程序操作系统构架，包括 386，amd64，arm
set GOOS=linux    设置可执行程序运行操作系统，支持 darwin，freebsd，linux，
go install github.com/mitchellh/gox@v1.0.1 交叉编译（不通平台，windows，linux，macos等）

build到当前目录
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build namenode_exporter.go
build到指定目录
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./target/nm_exporter namenode_exporter.go
```

How to run
```
上传config.yml和nm_exporter到服务器
sudo -s
chmod +x nm_exporter
nohup ./nm_exporter &
活
./start_exporter 
```

Help on flags of namenode_exporter:
```
-namenode.jmx.url string
    Hadoop JMX URL. (default "http://localhost:50070/jmx")
-web.listen-address string
    Address on which to expose metrics and web interface. (default ":9070")
-web.telemetry-path string
    Path under which to expose metrics. (default "/metrics")
```

Help on flags of resourcemanager_exporter:
```
-resourcemanager.url string
    Hadoop ResourceManager URL. (default "http://localhost:8088")
-web.listen-address string
    Address on which to expose metrics and web interface. (default ":9088")
-web.telemetry-path string
    Path under which to expose metrics. (default "/metrics")
```

go lang 1.14.4 glide
