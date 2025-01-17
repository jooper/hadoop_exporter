package Utiles

import (
	"github.com/robfig/cron"
)

func StartScheduler(cronStr string) {
	yml := Yml()
	c := cron.New()
	if cronStr == "" {
		cronStr = "0/5 * * * * ?"
	}
	//每5秒执行一次
	c.AddFunc(cronStr, func() {
		GetJMxMsg(yml.NameNodeExporterIp, yml.NameNodeExporterPort)
	})
	c.Start()
}

func StartSchedulerWithCron(ip string, port string, cronStr string) {
	c := cron.New()
	if cronStr == "" {
		cronStr = "0/5 * * * * ?"
	}
	//每5秒执行一次
	c.AddFunc(cronStr, func() {
		GetJMxMsg(ip, port)
	})
	c.Start()
}
