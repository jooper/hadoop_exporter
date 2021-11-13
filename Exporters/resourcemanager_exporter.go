package main

/*
fun: resourcemanager(yarn) exporter
auth:jwp
date：2021/11/4
*/
import (
	"encoding/json"
	"flag"
	"github.com/wyukawa/hadoop_exporter/Utiles"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/log"
)

const (
	namespace = "resourcemanager"
)

var yml = Utiles.Yml()
var (
	listenAddress      = flag.String("web.listen-address", ":"+yml.ResourceManagerExporterPort, "Address on which to expose metrics and web interface.")
	metricsPath        = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	resourceManagerUrl = flag.String("resourcemanager.url", yml.ResourceManagerJmx, "Hadoop ResourceManager URL.")
)

/*
fun:定义Gauge类型的指标项
auth:jwp
date：2021/11/4
*/
type Exporter struct {
	url                   string
	activeNodes           prometheus.Gauge
	rebootedNodes         prometheus.Gauge
	decommissionedNodes   prometheus.Gauge
	unhealthyNodes        prometheus.Gauge
	lostNodes             prometheus.Gauge
	totalNodes            prometheus.Gauge
	totalVirtualCores     prometheus.Gauge
	availableMB           prometheus.Gauge
	reservedMB            prometheus.Gauge
	appsKilled            prometheus.Gauge
	appsFailed            prometheus.Gauge
	appsRunning           prometheus.Gauge
	appsPending           prometheus.Gauge
	appsCompleted         prometheus.Gauge
	appsSubmitted         prometheus.Gauge
	allocatedMB           prometheus.Gauge
	reservedVirtualCores  prometheus.Gauge
	availableVirtualCores prometheus.Gauge
	allocatedVirtualCores prometheus.Gauge
	containersAllocated   prometheus.Gauge
	containersReserved    prometheus.Gauge
	containersPending     prometheus.Gauge
	totalMB               prometheus.Gauge
}

/*
fun:构建新exporter
auth:jwp
date：2021/11/4
*/
func NewExporter(url string) *Exporter {
	return &Exporter{
		url: url,
		activeNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "activeNodes",
			Help:      "activeNodes",
		}),
		rebootedNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "rebootedNodes",
			Help:      "rebootedNodes",
		}),
		decommissionedNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "decommissionedNodes",
			Help:      "decommissionedNodes",
		}),
		unhealthyNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "unhealthyNodes",
			Help:      "unhealthyNodes",
		}),
		lostNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "lostNodes",
			Help:      "lostNodes",
		}),
		totalNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "totalNodes",
			Help:      "totalNodes",
		}),
		totalVirtualCores: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "totalVirtualCores",
			Help:      "totalVirtualCores",
		}),
		availableMB: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "availableMB",
			Help:      "availableMB",
		}),
		reservedMB: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "reservedMB",
			Help:      "reservedMB",
		}),
		appsKilled: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "appsKilled",
			Help:      "appsKilled",
		}),
		appsFailed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "appsFailed",
			Help:      "appsFailed",
		}),
		appsRunning: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "appsRunning",
			Help:      "appsRunning",
		}),
		appsPending: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "appsPending",
			Help:      "appsPending",
		}),
		appsCompleted: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "appsCompleted",
			Help:      "appsCompleted",
		}),
		appsSubmitted: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "appsSubmitted",
			Help:      "appsSubmitted",
		}),
		allocatedMB: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "allocatedMB",
			Help:      "allocatedMB",
		}),
		reservedVirtualCores: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "reservedVirtualCores",
			Help:      "reservedVirtualCores",
		}),
		availableVirtualCores: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "availableVirtualCores",
			Help:      "availableVirtualCores",
		}),
		allocatedVirtualCores: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "allocatedVirtualCores",
			Help:      "allocatedVirtualCores",
		}),
		containersAllocated: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "containersAllocated",
			Help:      "containersAllocated",
		}),
		containersReserved: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "containersReserved",
			Help:      "containersReserved",
		}),
		containersPending: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "containersPending",
			Help:      "containersPending",
		}),
		totalMB: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "totalMB",
			Help:      "totalMB",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.activeNodes.Describe(ch)
	e.rebootedNodes.Describe(ch)
	e.decommissionedNodes.Describe(ch)
	e.unhealthyNodes.Describe(ch)
	e.lostNodes.Describe(ch)
	e.totalNodes.Describe(ch)
	e.totalVirtualCores.Describe(ch)
	e.availableMB.Describe(ch)
	e.reservedMB.Describe(ch)
	e.appsKilled.Describe(ch)
	e.appsFailed.Describe(ch)
	e.appsRunning.Describe(ch)
	e.appsPending.Describe(ch)
	e.appsCompleted.Describe(ch)
	e.appsSubmitted.Describe(ch)
	e.allocatedMB.Describe(ch)
	e.reservedVirtualCores.Describe(ch)
	e.availableVirtualCores.Describe(ch)
	e.allocatedVirtualCores.Describe(ch)
	e.containersAllocated.Describe(ch)
	e.containersReserved.Describe(ch)
	e.containersPending.Describe(ch)
	e.totalMB.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	resp, err := http.Get(e.url + "/ws/v1/cluster/metrics")
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
	cm := m["clusterMetrics"].(map[string]interface{})
	e.activeNodes.Set(cm["activeNodes"].(float64))
	e.rebootedNodes.Set(cm["rebootedNodes"].(float64))
	e.decommissionedNodes.Set(cm["decommissionedNodes"].(float64))
	e.unhealthyNodes.Set(cm["unhealthyNodes"].(float64))
	e.lostNodes.Set(cm["lostNodes"].(float64))
	e.totalNodes.Set(cm["totalNodes"].(float64))
	e.totalVirtualCores.Set(cm["totalVirtualCores"].(float64))
	e.availableMB.Set(cm["availableMB"].(float64))
	e.reservedMB.Set(cm["reservedMB"].(float64))
	e.appsKilled.Set(cm["appsKilled"].(float64))
	e.appsFailed.Set(cm["appsFailed"].(float64))
	e.appsRunning.Set(cm["appsRunning"].(float64))
	e.appsPending.Set(cm["appsPending"].(float64))
	e.appsCompleted.Set(cm["appsCompleted"].(float64))
	e.appsSubmitted.Set(cm["appsSubmitted"].(float64))
	e.allocatedMB.Set(cm["allocatedMB"].(float64))
	e.reservedVirtualCores.Set(cm["reservedVirtualCores"].(float64))
	e.availableVirtualCores.Set(cm["availableVirtualCores"].(float64))
	e.allocatedVirtualCores.Set(cm["allocatedVirtualCores"].(float64))
	e.containersAllocated.Set(cm["containersAllocated"].(float64))
	e.containersReserved.Set(cm["containersReserved"].(float64))
	e.containersPending.Set(cm["containersPending"].(float64))
	e.totalMB.Set(cm["totalMB"].(float64))

	e.activeNodes.Collect(ch)
	e.rebootedNodes.Collect(ch)
	e.decommissionedNodes.Collect(ch)
	e.unhealthyNodes.Collect(ch)
	e.lostNodes.Collect(ch)
	e.totalNodes.Collect(ch)
	e.totalVirtualCores.Collect(ch)
	e.availableMB.Collect(ch)
	e.reservedMB.Collect(ch)
	e.appsKilled.Collect(ch)
	e.appsFailed.Collect(ch)
	e.appsRunning.Collect(ch)
	e.appsPending.Collect(ch)
	e.appsCompleted.Collect(ch)
	e.appsSubmitted.Collect(ch)
	e.allocatedMB.Collect(ch)
	e.reservedVirtualCores.Collect(ch)
	e.availableVirtualCores.Collect(ch)
	e.allocatedVirtualCores.Collect(ch)
	e.containersAllocated.Collect(ch)
	e.containersReserved.Collect(ch)
	e.containersPending.Collect(ch)
	e.totalMB.Collect(ch)

	pushToFaas(ConvertMetrics(data))

}

/*
fun:获取指标数据，并推送到faas平台
auth:jwp
date：2021/11/4
*/
func ConvertMetrics(data []byte) map[string]interface{} {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		log.Error(err)
	}
	m := f.(map[string]interface{})
	cm := m["clusterMetrics"].(map[string]interface{})

	faasMap := make(map[string]interface{})

	faasMap["activeNodes"] = cm["activeNodes"].(float64)
	faasMap["rebootedNodes"] = cm["rebootedNodes"].(float64)
	faasMap["decommissionedNodes"] = cm["decommissionedNodes"].(float64)
	faasMap["unhealthyNodes"] = cm["unhealthyNodes"].(float64)
	faasMap["lostNodes"] = cm["lostNodes"].(float64)
	faasMap["totalNodes"] = cm["totalNodes"].(float64)
	faasMap["totalVirtualCores"] = cm["totalVirtualCores"].(float64)
	faasMap["availableMB"] = cm["availableMB"].(float64)
	faasMap["reservedMB"] = cm["reservedMB"].(float64)
	faasMap["appsKilled"] = cm["appsKilled"].(float64)
	faasMap["appsFailed"] = cm["appsFailed"].(float64)
	faasMap["appsRunning"] = cm["appsRunning"].(float64)
	faasMap["appsPending"] = cm["appsPending"].(float64)
	faasMap["appsCompleted"] = cm["appsCompleted"].(float64)
	faasMap["appsSubmitted"] = cm["appsSubmitted"].(float64)
	faasMap["allocatedMB"] = cm["allocatedMB"].(float64)
	faasMap["reservedVirtualCores"] = cm["reservedVirtualCores"].(float64)
	faasMap["availableVirtualCores"] = cm["availableVirtualCores"].(float64)
	faasMap["allocatedVirtualCores"] = cm["allocatedVirtualCores"].(float64)
	faasMap["containersAllocated"] = cm["containersAllocated"].(float64)
	faasMap["containersReserved"] = cm["containersReserved"].(float64)
	faasMap["containersPending"] = cm["containersPending"].(float64)
	faasMap["totalMB"] = cm["totalMB"].(float64)
	return faasMap
}

func pushToFaas(faasMap map[string]interface{}) {
	Utiles.PushResourceManagerMetricsToFaas(faasMap)
	str, _ := json.Marshal(faasMap)
	log.Info(string(str))
}

func main() {
	//开启调度,需在http服务前，调度crontab表达式在配置文件中
	Utiles.StartSchedulerWithCron(yml.ResourceManagerExporterIp, yml.ResourceManagerExporterPort, yml.CronStr)

	flag.Parse()

	exporter := NewExporter(*resourceManagerUrl)
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>ResourceManager Exporter</title></head>
		<body>
		<h1>ResourceManager Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
	})
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
