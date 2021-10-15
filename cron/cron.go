// Package cron 定时任务
package cron

import (
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promSyncLabels = []string{
		"jobname",
	}
	promSyncError = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cron",
			Name:      "sync_error",
			Help:      "cron sync job error",
		}, promSyncLabels,
	)
)

// RunCronJobs 启动定时任务
func RunCronJobs(async bool) {
	timezone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logging.Errorf(nil, "RunCronJobs time LoadLocation error:%v, using Local timezone as default", err.Error())
		timezone, _ = time.LoadLocation("Local")
	}
	logging.Debugf(nil, "cron timezone:%v", timezone)
	sched := gocron.NewScheduler(timezone)
	// 同步基金净值列表和4433列表
	sched.Cron("0 18 * * 1-5").Do(SyncFund)
	// 同步东方财富行业列表
	sched.Cron("0 4 1 * *").Do(SyncIndustryList)
	if async {
		sched.StartAsync()
	} else {
		sched.StartBlocking()
	}
}
