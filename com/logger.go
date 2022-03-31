package com

import (
	log "github.com/sirupsen/logrus"
)

// Log teamkit 公共日志输出
var Log *log.Entry

func init() {
	log.SetReportCaller(true)
	// 日志前缀设置
	Log = log.WithField("pkg", "appstoreconnect")
}
