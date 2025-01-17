package main

import (
	"encoding/json"
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wyukawa/hadoop_exporter/Utiles"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/log"

	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/log"
	"fmt"
	"os"
)

const (
	namespace = "datanode"
)

var yml = Utiles.Yml()

var (
	listenAddress  = flag.String("web.listen-address", ":"+yml.DataNodeExporterPort, "Address on which to expose metrics and web interface.")
	metricsPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	datanodeJmxUrl = flag.String("datanode.jmx.url", yml.DataNodeJmx, "Hadoop JMX URL.")
)

type Exporter struct {
	url                    string
	WritesFromRemoteClient prometheus.Gauge
	WritesFromLocalClient  prometheus.Gauge
	WriteBlockOpNumOps     prometheus.Gauge
	VolumeFailures         prometheus.Gauge
	TotalWriteTime         prometheus.Gauge
	ThreadsBlocked         prometheus.Gauge
	ReadBlockOpAvgTime     prometheus.Gauge
	HeartbeatsNumOps       prometheus.Gauge
	HeartbeatsAvgTime      prometheus.Gauge
	GcTimeMillis           prometheus.Gauge
	GcCount                prometheus.Gauge
	DatanodeNetworkErrors  prometheus.Gauge
	BytesWritten           prometheus.Gauge
	//BlocksWritten		prometheus.Gauge
	BlocksReplicated    prometheus.Gauge
	BlockReportsNumOps  prometheus.Gauge
	BlockReportsAvgTime prometheus.Gauge
}

func NewExporter(url string) *Exporter {
	return &Exporter{
		url: url,
		WritesFromRemoteClient: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "WritesFromRemoteClient",
			Help:      "WritesFromRemoteClient",
		}),
		WritesFromLocalClient: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "WritesFromLocalClient",
			Help:      "WritesFromLocalClient",
		}),
		WriteBlockOpNumOps: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "WriteBlockOpNumOps",
			Help:      "WriteBlockOpNumOps",
		}),
		VolumeFailures: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "VolumeFailures",
			Help:      "VolumeFailures",
		}),
		TotalWriteTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "TotalWriteTime",
			Help:      "TotalWriteTime",
		}),
		ThreadsBlocked: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ThreadsBlocked",
			Help:      "ThreadsBlocked",
		}),
		ReadBlockOpAvgTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ReadBlockOpAvgTime",
			Help:      "ReadBlockOpAvgTime",
		}),
		HeartbeatsNumOps: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "HeartbeatsNumOps",
			Help:      "HeartbeatsNumOps",
		}),
		HeartbeatsAvgTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "HeartbeatsAvgTime",
			Help:      "HeartbeatsAvgTime",
		}),
		GcTimeMillis: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "GcTimeMillis",
			Help:      "GcTimeMillis",
		}),
		GcCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "GcCount",
			Help:      "GcCount",
		}),
		DatanodeNetworkErrors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "DatanodeNetworkErrors",
			Help:      "DatanodeNetworkErrors",
		}),
		BytesWritten: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BytesWritten",
			Help:      "BytesWritten",
		}),
		BlocksReplicated: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlocksReplicated",
			Help:      "BlocksReplicated",
		}),
		BlockReportsNumOps: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlockReportsNumOps",
			Help:      "BlockReportsNumOps",
		}),
		BlockReportsAvgTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlockReportsAvgTime",
			Help:      "BlockReportsAvgTime",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.WritesFromRemoteClient.Describe(ch)
	e.WritesFromLocalClient.Describe(ch)
	e.WriteBlockOpNumOps.Describe(ch)
	e.VolumeFailures.Describe(ch)
	e.TotalWriteTime.Describe(ch)
	e.ThreadsBlocked.Describe(ch)
	e.ReadBlockOpAvgTime.Describe(ch)
	e.HeartbeatsNumOps.Describe(ch)
	e.HeartbeatsAvgTime.Describe(ch)
	e.GcTimeMillis.Describe(ch)
	e.GcCount.Describe(ch)
	e.DatanodeNetworkErrors.Describe(ch)
	e.BytesWritten.Describe(ch)
	//e.BlocksWritten.Describe(ch)
	e.BlocksReplicated.Describe(ch)
	e.BlockReportsNumOps.Describe(ch)
	e.BlockReportsAvgTime.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	hostName, osErr := os.Hostname()
	if osErr != nil {
		fmt.Printf("%s", osErr)
	}

	resp, err := http.Get(e.url)
	if err != nil {
		//log.Error(err)
		fmt.Errorf("erro:", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Error(err)
		fmt.Errorf("erro:", err)
	}
	var f interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		//log.Error(err)
		fmt.Errorf("erro:", err)
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error...")
		}
	}()
	m := f.(map[string]interface{})
	var nameList = m["beans"].([]interface{})
	for _, nameData := range nameList {
		nameDataMap := nameData.(map[string]interface{})
		//log.Info(nameDataMap["name"])
		if nameDataMap["name"] == "Hadoop:service=DataNode,name=DataNodeActivity-"+hostName+"-50010" {
			e.WritesFromRemoteClient.Set(nameDataMap["WritesFromRemoteClient"].(float64))
			e.WritesFromLocalClient.Set(nameDataMap["WritesFromLocalClient"].(float64))
			e.WriteBlockOpNumOps.Set(nameDataMap["WriteBlockOpNumOps"].(float64))
			e.VolumeFailures.Set(nameDataMap["VolumeFailures"].(float64))
			e.TotalWriteTime.Set(nameDataMap["TotalWriteTime"].(float64))
			e.ReadBlockOpAvgTime.Set(nameDataMap["ReadBlockOpAvgTime"].(float64))
			e.HeartbeatsNumOps.Set(nameDataMap["HeartbeatsNumOps"].(float64))
			e.HeartbeatsAvgTime.Set(nameDataMap["HeartbeatsAvgTime"].(float64))
			e.DatanodeNetworkErrors.Set(nameDataMap["DatanodeNetworkErrors"].(float64))
			e.BytesWritten.Set(nameDataMap["BytesWritten"].(float64))
			//e.BlocksWritten.Set(nameDataMap["BlocksWritten"].(float64))
			e.BlocksReplicated.Set(nameDataMap["BlocksReplicated"].(float64))
			e.BlockReportsNumOps.Set(nameDataMap["BlockReportsNumOps"].(float64))
			e.BlockReportsAvgTime.Set(nameDataMap["BlockReportsAvgTime"].(float64))

		}
		if nameDataMap["name"] == "Hadoop:service=DataNode,name=JvmMetrics" {
			e.GcTimeMillis.Set(nameDataMap["GcTimeMillis"].(float64))
			e.GcCount.Set(nameDataMap["GcCount"].(float64))
			e.ThreadsBlocked.Set(nameDataMap["ThreadsBlocked"].(float64))
		}

		ConvertMetrics(nameDataMap)
	}
	e.WritesFromRemoteClient.Collect(ch)
	e.WritesFromLocalClient.Collect(ch)
	e.WriteBlockOpNumOps.Collect(ch)
	e.VolumeFailures.Collect(ch)
	e.TotalWriteTime.Collect(ch)
	e.ThreadsBlocked.Collect(ch)
	e.ReadBlockOpAvgTime.Collect(ch)
	e.HeartbeatsNumOps.Collect(ch)
	e.HeartbeatsAvgTime.Collect(ch)
	e.GcTimeMillis.Collect(ch)
	e.GcCount.Collect(ch)
	e.DatanodeNetworkErrors.Collect(ch)
	e.BytesWritten.Collect(ch)
	//e.BlocksWritten.Collect(ch)
	e.BlocksReplicated.Collect(ch)
	e.BlockReportsNumOps.Collect(ch)
	e.BlockReportsAvgTime.Collect(ch)
}

/*
fun:获取指标数据，并推送到faas平台
auth:jwp
date：2021/11/10
*/
func ConvertMetrics(nameDataMap map[string]interface{}) {
	faasMap := make(map[string]interface{})
	hostName, _ := os.Hostname()
	//hostName = "hdp03"

	if nameDataMap["name"] == "Hadoop:service=DataNode,name=DataNodeActivity-"+hostName+"-50010" {
		faasMap["WritesFromRemoteClient"] = nameDataMap["WritesFromRemoteClient"].(float64)
		faasMap["WritesFromLocalClient"] = nameDataMap["WritesFromLocalClient"].(float64)
		faasMap["WriteBlockOpNumOps"] = nameDataMap["WriteBlockOpNumOps"].(float64)
		faasMap["VolumeFailures"] = nameDataMap["VolumeFailures"].(float64)
		faasMap["TotalWriteTime"] = nameDataMap["TotalWriteTime"].(float64)
		faasMap["ReadBlockOpAvgTime"] = nameDataMap["ReadBlockOpAvgTime"].(float64)
		faasMap["HeartbeatsNumOps"] = nameDataMap["HeartbeatsNumOps"].(float64)
		faasMap["HeartbeatsAvgTime"] = nameDataMap["HeartbeatsAvgTime"].(float64)
		faasMap["DatanodeNetworkErrors"] = nameDataMap["DatanodeNetworkErrors"].(float64)
		faasMap["BytesWritten"] = nameDataMap["BytesWritten"].(float64)
		//e.BlocksWritten.Set(nameDataMap["BlocksWritten"].(float64))
		faasMap["BlocksReplicated"] = nameDataMap["BlocksReplicated"].(float64)
		faasMap["BlockReportsNumOps"] = nameDataMap["BlockReportsNumOps"].(float64)
		faasMap["BlockReportsAvgTime"] = nameDataMap["BlockReportsAvgTime"].(float64)

		//str, _ := json.Marshal(nameDataMap)
		//log.Info(string(str))
	}
	if nameDataMap["name"] == "Hadoop:service=DataNode,name=JvmMetrics" {
		faasMap["GcTimeMillis"] = nameDataMap["GcTimeMillis"].(float64)
		faasMap["GcCount"] = nameDataMap["GcCount"].(float64)
		faasMap["ThreadsBlocked"] = nameDataMap["ThreadsBlocked"].(float64)
	}

	if nameDataMap["name"] == "Hadoop:service=DataNode,name=FSDatasetState" {
		faasMap["Capacity"] = nameDataMap["Capacity"].(float64)
		faasMap["DfsUsed"] = nameDataMap["DfsUsed"].(float64)
	}
	if len(faasMap) > 0 {
		faasMap["datanode_name"] = hostName
		Utiles.PushDataNodeMetricsToFaas(faasMap)
		str, _ := json.Marshal(faasMap)
		log.Info(string(str))
	}
}

func main() {

	//开启调度,需在http服务前，调度crontab表达式在配置文件中
	Utiles.StartSchedulerWithCron(yml.DataNodeExporterIp, yml.DataNodeExporterPort, yml.CronStr)

	flag.Parse()
	exporter := NewExporter(*datanodeJmxUrl)
	prometheus.MustRegister(exporter)

	fmt.Printf("Starting Server: %s", *listenAddress)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>DataNode Exporter</title></head>
		<body>
		<h1>DataNode Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
	})
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		//log.Fatal(err)
		fmt.Errorf("erro:", err)
	}
}
