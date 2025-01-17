package Utiles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/log"
	"io/ioutil"
	"net/http"
)

//const apiUrl = "http://10.231.143.223:18991/hit/xrobot"

var yml = Yml()

var apiUrl = yml.FaasApiUrl

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

var ExporterTenatid = yml.ExportProject
var ExporterAppid = yml.ExportAppId
var ExporterSkill = yml.ExportSkill

func PushNameNodeMetricsToFaas(v map[string]interface{}) {
	pushToFaas(ExporterTenatid, ExporterAppid, ExporterSkill, "hadoop_namenode", v)
}

func PushDataNodeMetricsToFaas(v map[string]interface{}) {
	pushToFaas(ExporterTenatid, ExporterAppid, ExporterSkill, "hadoop_datanode", v)
}

func PushResourceManagerMetricsToFaas(v map[string]interface{}) {
	pushToFaas(ExporterTenatid, ExporterAppid, ExporterSkill, "hadoop_resourcemanager", v)
}

func PushJournalNodeMetricsToFaas(v map[string]interface{}) {
	pushToFaas(ExporterTenatid, ExporterAppid, ExporterSkill, "hadoop_journalnode", v)
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

	//alertMsg("MarketCenter_CRMTeam", "AlterMsgForCrm", "AlterMsgReceiveAndSend-AlterMsg", content, tels, "test", "test")
	alertMsg(yml.AlertProject, yml.AlertAppId, yml.AlertSkill, content, tels, "test", "test")
}

func alertMsg(tenatid string, appid string, skill string, content string, tels string, reginon string, key string) {
	para2 := AutoGeneratedMsg{
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

	fmt.Printf("\n进入短信发送模块\n")
	paraJson1, _ := json.Marshal(para2)

	fmt.Printf("%+v\n", para2)
	// 注意上传config.yml
	resp1, _ := http.Post(apiUrl, contentType, bytes.NewReader(paraJson1))

	log.Info("resp1", apiUrl, contentType, resp1)
	//这里要判空，否则打包后会报空指针
	if resp1 != nil {
		body1, _ := ioutil.ReadAll(resp1.Body)
		fmt.Printf("短信发送成功\n")
		fmt.Println(string(body1))
	}

}
