package Utiles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//const apiUrl = "http://10.231.143.223:18991/hit/xrobot"
var apiUrl = Yml().FaasApiUrl

const contentType = "application/json"

//定义请求体格式
type AutoGenerated struct {
	Intent  string  `json:"intent"`
	Request Request `json:"request"`
}
type Request struct {
	Tenantid     string                 `json:"tenantid"`
	Appid        string                 `json:"appid"`
	Skill        string                 `json:"skill"`
	EsIndex      string                 `json:"es_index"`
	CurTimestmap string                 `json:"cur_timestmap"`
	CdcTimestamp string                 `json:"cdc_timestamp"`
	CdcTableName string                 `json:"cdc_table_name"`
	DelayValue   int                    `json:"delay_value"`
	Content      map[string]interface{} `json:"content"`
}

func PushNameNodeMetricsToFaas(v map[string]interface{}) {
	pushToFaas("MarketCenter_SaleTeam", "streamsets_monitor", "adbStreamsets-sinkMsgToEs", "hadoop_namenode", v)
}

func PushDataNodeMetricsToFaas(v map[string]interface{}) {
	pushToFaas("MarketCenter_SaleTeam", "streamsets_monitor", "adbStreamsets-sinkMsgToEs", "hadoop_datanode", v)
}

func PushResourceManagerMetricsToFaas(v map[string]interface{}) {
	pushToFaas("MarketCenter_SaleTeam", "streamsets_monitor", "adbStreamsets-sinkMsgToEs", "hadoop_resourcemanager", v)
}

func PushJournalNodeMetricsToFaas(v map[string]interface{}) {
	pushToFaas("MarketCenter_SaleTeam", "streamsets_monitor", "adbStreamsets-sinkMsgToEs", "hadoop_journalnode", v)
}

func pushToFaas(tenatid string, appid string, skill string, esindex string, v map[string]interface{}) {
	para := AutoGenerated{
		Intent: "API",
		Request: Request{
			Tenantid: tenatid,
			Appid:    appid,
			Skill:    skill,
			EsIndex:  esindex,
			Content:  v,
		}}

	paraJson, _ := json.Marshal(para)
	resp, _ := http.Post(apiUrl, contentType, bytes.NewReader(paraJson))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func MapToJson(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func GetJMxMsg(ip string, port string) string {
	var url = fmt.Sprintf("http://%s:%s/metrics", ip, port)
	data := Get(url)
	return data
}

type AutoGeneratedMsg struct {
	Intent  string     `json:"intent"`
	Request RequestMsg `json:"request"`
}
type RequestMsg struct {
	Tenantid string `json:"tenantid"`
	Appid    string `json:"appid"`
	Skill    string `json:"skill"`
	Content  string `json:"content"`
	Tels     string `json:"tels"`
	Region   string `json:"region"`
	Key      string `json:"key"`
}

func AlertMsg(tels string, content string) {
	alertMsg("MarketCenter_CRMTeam", "AlterMsgForCrm", "AlterMsgReceiveAndSend-AlterMsg", content, tels, "test", "test")
}

func alertMsg(tenatid string, appid string, skill string, content string, tels string, reginon string, key string) {
	para := AutoGeneratedMsg{
		Intent: "API",
		Request: RequestMsg{
			Tenantid: tenatid,
			Appid:    appid,
			Skill:    skill,
			Content:  content,
			Tels:     tels,
			Region:   reginon,
			Key:      key,
		}}

	paraJson, _ := json.Marshal(para)
	resp, _ := http.Post(apiUrl, contentType, bytes.NewReader(paraJson))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
