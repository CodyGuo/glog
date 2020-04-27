package main

import (
	"os"

	"github.com/CodyGuo/glog"
)

func main() {
	glog.Debug("hello debug")
	glog.Info("hello info")

	customLog := glog.New(os.Stdout,
		glog.WithLevel(glog.DEBUG),
		glog.WithFlags(glog.LstdFlags),
		glog.WithPrefix("[customLog] "))

	customLog.Debug("hello debug")
	customLog.Info("hello info")
}
