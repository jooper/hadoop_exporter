package main

import (
	"encoding/json"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/log"
	"github.com/wyukawa/hadoop_exporter/Utiles"
	"io/ioutil"
	"net/http"
)

const (
	namespace = "namenode" //指标前缀 如：namenode_BlocksTot
)

var (
	listenAddress  = flag.String("web.listen-address", ":"+Utiles.Yml().NameNodeExporterPort, "Address on which to expose metrics and web interface.")
	metricsPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	namenodeJmxUrl = flag.String("namenode.jmx.url", Utiles.Yml().NameNodeJmx, "Hadoop JMX URL.")
)

type Exporter struct {
	url                      string
	MissingBlocks            prometheus.Gauge
	CapacityTotal            prometheus.Gauge
	CapacityUsed             prometheus.Gauge
	CapacityRemaining        prometheus.Gauge
	CapacityUsedNonDFS       prometheus.Gauge
	BlocksTotal              prometheus.Gauge
	FilesTotal               prometheus.Gauge
	CorruptBlocks            prometheus.Gauge
	ExcessBlocks             prometheus.Gauge
	StaleDataNodes           prometheus.Gauge
	pnGcCount                prometheus.Gauge
	pnGcTime                 prometheus.Gauge
	cmsGcCount               prometheus.Gauge
	cmsGcTime                prometheus.Gauge
	heapMemoryUsageCommitted prometheus.Gauge
	heapMemoryUsageInit      prometheus.Gauge
	heapMemoryUsageMax       prometheus.Gauge
	heapMemoryUsageUsed      prometheus.Gauge
	isActive                 prometheus.Gauge
}

func NewExporter(url string) *Exporter {
	return &Exporter{
		url: url,
		MissingBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "MissingBlocks",
			Help:      "MissingBlocks",
		}),
		CapacityTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityTotal",
			Help:      "CapacityTotal",
		}),
		CapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityUsed",
			Help:      "CapacityUsed",
		}),
		CapacityRemaining: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityRemaining",
			Help:      "CapacityRemaining",
		}),
		CapacityUsedNonDFS: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityUsedNonDFS",
			Help:      "CapacityUsedNonDFS",
		}),
		BlocksTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlocksTotal",
			Help:      "BlocksTotal",
		}),
		FilesTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "FilesTotal",
			Help:      "FilesTotal",
		}),
		CorruptBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CorruptBlocks",
			Help:      "CorruptBlocks",
		}),
		ExcessBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ExcessBlocks",
			Help:      "ExcessBlocks",
		}),
		StaleDataNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "StaleDataNodes",
			Help:      "StaleDataNodes",
		}),
		pnGcCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ParNew_CollectionCount",
			Help:      "ParNew GC Count",
		}),
		pnGcTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ParNew_CollectionTime",
			Help:      "ParNew GC Time",
		}),
		cmsGcCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ConcurrentMarkSweep_CollectionCount",
			Help:      "ConcurrentMarkSweep GC Count",
		}),
		cmsGcTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ConcurrentMarkSweep_CollectionTime",
			Help:      "ConcurrentMarkSweep GC Time",
		}),
		heapMemoryUsageCommitted: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageCommitted",
			Help:      "heapMemoryUsageCommitted",
		}),
		heapMemoryUsageInit: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageInit",
			Help:      "heapMemoryUsageInit",
		}),
		heapMemoryUsageMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageMax",
			Help:      "heapMemoryUsageMax",
		}),
		heapMemoryUsageUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageUsed",
			Help:      "heapMemoryUsageUsed",
		}),
		isActive: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "isActive",
			Help:      "isActive",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.MissingBlocks.Describe(ch)
	e.CapacityTotal.Describe(ch)
	e.CapacityUsed.Describe(ch)
	e.CapacityRemaining.Describe(ch)
	e.CapacityUsedNonDFS.Describe(ch)
	e.BlocksTotal.Describe(ch)
	e.FilesTotal.Describe(ch)
	e.CorruptBlocks.Describe(ch)
	e.ExcessBlocks.Describe(ch)
	e.StaleDataNodes.Describe(ch)
	e.pnGcCount.Describe(ch)
	e.pnGcTime.Describe(ch)
	e.cmsGcCount.Describe(ch)
	e.cmsGcTime.Describe(ch)
	e.heapMemoryUsageCommitted.Describe(ch)
	e.heapMemoryUsageInit.Describe(ch)
	e.heapMemoryUsageMax.Describe(ch)
	e.heapMemoryUsageUsed.Describe(ch)
	e.isActive.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	resp, err := http.Get(e.url)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var f interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		log.Error(err)
	}

	m := f.(map[string]interface{})

	//faasList := list.List{}

	var nameList = m["beans"].([]interface{})
	for _, nameData := range nameList {
		nameDataMap := nameData.(map[string]interface{})

		faasMap := make(map[string]interface{})

		if nameDataMap["name"] == "Hadoop:service=NameNode,name=FSNamesystem" {
			e.MissingBlocks.Set(nameDataMap["MissingBlocks"].(float64))
			e.CapacityTotal.Set(nameDataMap["CapacityTotal"].(float64))
			e.CapacityUsed.Set(nameDataMap["CapacityUsed"].(float64))
			e.CapacityRemaining.Set(nameDataMap["CapacityRemaining"].(float64))

			e.CapacityUsedNonDFS.Set(nameDataMap["CapacityUsedNonDFS"].(float64))
			e.BlocksTotal.Set(nameDataMap["BlocksTotal"].(float64))
			e.FilesTotal.Set(nameDataMap["FilesTotal"].(float64))
			e.CorruptBlocks.Set(nameDataMap["CorruptBlocks"].(float64))
			e.ExcessBlocks.Set(nameDataMap["ExcessBlocks"].(float64))
			e.StaleDataNodes.Set(nameDataMap["StaleDataNodes"].(float64))

			faasMap["MissingBlocks"] = nameDataMap["MissingBlocks"].(float64)
			faasMap["NumLiveDataNodes"] = nameDataMap["NumLiveDataNodes"].(float64)
			faasMap["CapacityTotal"] = nameDataMap["CapacityTotal"].(float64)
			faasMap["CapacityUsed"] = nameDataMap["CapacityUsed"].(float64)
			faasMap["CapacityRemaining"] = nameDataMap["CapacityRemaining"].(float64)
			faasMap["CapacityUsedNonDFS"] = nameDataMap["CapacityUsedNonDFS"].(float64)
			faasMap["BlocksTotal"] = nameDataMap["BlocksTotal"].(float64)
			faasMap["FilesTotal"] = nameDataMap["FilesTotal"].(float64)
			faasMap["CorruptBlocks"] = nameDataMap["CorruptBlocks"].(float64)
			faasMap["ExcessBlocks"] = nameDataMap["ExcessBlocks"].(float64)
			faasMap["StaleDataNodes"] = nameDataMap["StaleDataNodes"].(float64)
		}
		if nameDataMap["name"] == "java.lang:type=GarbageCollector,name=ParNew" {
			e.pnGcCount.Set(nameDataMap["CollectionCount"].(float64))
			e.pnGcTime.Set(nameDataMap["CollectionTime"].(float64))

			faasMap["pnGcCount"] = nameDataMap["CollectionCount"].(float64)
			faasMap["pnGcTime"] = nameDataMap["CollectionTime"].(float64)
		}
		if nameDataMap["name"] == "java.lang:type=GarbageCollector,name=ConcurrentMarkSweep" {
			e.cmsGcCount.Set(nameDataMap["CollectionCount"].(float64))
			e.cmsGcTime.Set(nameDataMap["CollectionTime"].(float64))

			faasMap["cmsGcCount"] = nameDataMap["CollectionCount"].(float64)
			faasMap["cmsGcTime"] = nameDataMap["CollectionTime"].(float64)

		}

		if nameDataMap["name"] == "java.lang:type=Memory" {
			heapMemoryUsage := nameDataMap["HeapMemoryUsage"].(map[string]interface{})
			e.heapMemoryUsageCommitted.Set(heapMemoryUsage["committed"].(float64))
			e.heapMemoryUsageInit.Set(heapMemoryUsage["init"].(float64))
			e.heapMemoryUsageMax.Set(heapMemoryUsage["max"].(float64))
			e.heapMemoryUsageUsed.Set(heapMemoryUsage["used"].(float64))

			faasMap["heapMemoryUsageCommitted"] = heapMemoryUsage["committed"].(float64)
			faasMap["heapMemoryUsageInit"] = heapMemoryUsage["init"].(float64)
			faasMap["heapMemoryUsageMax"] = heapMemoryUsage["max"].(float64)
			faasMap["heapMemoryUsageUsed"] = heapMemoryUsage["used"].(float64)
		}

		if nameDataMap["name"] == "Hadoop:service=NameNode,name=FSNamesystem" {
			if nameDataMap["tag.HAState"] == "active" {
				e.isActive.Set(1)
			} else {
				e.isActive.Set(0)
			}
		}

		if len(faasMap) > 0 {
			//faasMap["timestamp"]=time.Now()//.Format("2006-01-02 15:04:05")

			//faasMap["createtime"]=time.Now()

			//pushMetricsToFaas(faasMap)

			Utiles.PushNameNodeMetricsToFaas(faasMap)
			str, _ := json.Marshal(faasMap)
			log.Info(string(str))

			//faasList.PushBack(faasMap)

		}

	}

	//if faasList.Len()>0 {
	//	log.Info(faasList)
	//	pushMetricsToFaas(faasList)
	//}

	e.MissingBlocks.Collect(ch)
	e.CapacityTotal.Collect(ch)
	e.CapacityUsed.Collect(ch)
	e.CapacityRemaining.Collect(ch)
	e.CapacityUsedNonDFS.Collect(ch)
	e.BlocksTotal.Collect(ch)
	e.FilesTotal.Collect(ch)
	e.CorruptBlocks.Collect(ch)
	e.ExcessBlocks.Collect(ch)
	e.StaleDataNodes.Collect(ch)
	e.pnGcCount.Collect(ch)
	e.pnGcTime.Collect(ch)
	e.cmsGcCount.Collect(ch)
	e.cmsGcTime.Collect(ch)
	e.heapMemoryUsageCommitted.Collect(ch)
	e.heapMemoryUsageInit.Collect(ch)
	e.heapMemoryUsageMax.Collect(ch)
	e.heapMemoryUsageUsed.Collect(ch)
	e.isActive.Collect(ch)

	//log.Info(nameList)
	//pushMetricsToFaas(MapToJson(string(data)))

}

func main() {
	//开启调度,需在http服务前
	Utiles.StartScheduler(Utiles.Yml().CronStr)

	flag.Parse()

	exporter := NewExporter(*namenodeJmxUrl)
	// 注册监控指标
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listenAddress)
	//初始一个http handler
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>NameNode Exporter</title></head>
		<body>
		<h1>NameNode Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))

	})

	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}
