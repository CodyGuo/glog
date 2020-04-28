package main

import (
	"os"

	"github.com/CodyGuo/glog"
)

func main() {
	Info("hello info")
	Infof("hello info")

	glog.SetFlags(glog.LglogFlags)
	glog.Debug("hello debug")
	glog.Info("hello info")

	customLog := glog.New(os.Stdout,
		glog.WithLevel(glog.TRACE),
		glog.WithLevelLength(4),
		glog.WithFlags(glog.LglogFlags),
		glog.WithPrefix("[customLog] "))

	customLog.Trace("hello trace")
	customLog.Debug("hello debug")
	customLog.Info("hello info")
	customLog.Notice("hello notice")
	customLog.Warning("hello warning")
	customLog.Error("hello error")
	customLog.Critical("hello critical")
}
