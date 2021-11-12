package main

/*
获取altermanager webhook传递过来的告警json数据
*/

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wyukawa/hadoop_exporter/Utiles"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const serverPort = "8081"

type AutoGenerated struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []Alerts          `json:"alerts"`
	GroupLabels       GroupLabels       `json:"groupLabels"`
	CommonLabels      CommonLabels      `json:"commonLabels"`
	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
}
type Labels struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Severity  string `json:"severity"`
}
type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
type Alerts struct {
	Status       string      `json:"status"`
	Labels       Labels      `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
	Fingerprint  string      `json:"fingerprint"`
}
type GroupLabels struct {
	Alertname string `json:"alertname"`
}
type CommonLabels struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Severity  string `json:"severity"`
}
type CommonAnnotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
type Cmd struct {
	ReqType  int
	FileName string
}

func getPrometheusAlterManager() {
	//http://127.0.0.1:8080/webhook?

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

		if r.Method == "POST" {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("Read failed:", err)
			}
			defer r.Body.Close()

			prometheusAlterMsgjson := &AutoGenerated{}
			err = json.Unmarshal(b, prometheusAlterMsgjson)
			if err != nil {
				log.Println("json format error:", err)
			}

			log.Println("获取到到webhook数据:", prometheusAlterMsgjson)

			log.Println("报警信息：", prometheusAlterMsgjson.CommonAnnotations)

			//调用报警平台
			//sendMsgToFassAlterApi(nil)
		} else {

			log.Println("ONly support Post")
			fmt.Fprintf(w, "Only support post")
		}

	})
}

/*
调用报警Fass报警平台
*/
func sendMsgToFassAlterApi(faasMap map[string]interface{}) {
	Utiles.PushResourceManagerMetricsToFaas(faasMap)
	str, _ := json.Marshal(faasMap)
	log.Println(string(str))
}

func startListener() {
	var listenAddress = flag.String("web.listen-address", ":8081", "Address on which to expose metrics and web interface.")
	log.Printf("Starting Server: %s", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	getPrometheusAlterManager()
	startListener()
}