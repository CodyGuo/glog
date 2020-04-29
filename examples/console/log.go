package main

import (
	"os"

	"github.com/CodyGuo/glog"
)

func Info(v ...interface{}) {
	glog.AutoCallDepth()
	defer glog.ResetCallDepth()
	glog.SetFlags(glog.LglogFlags)
	glog.Info(v...)
}

func Infof(format string, v ...interface{}) {
	l := glog.New(os.Stdout, glog.WithFlags(glog.LglogFlags), glog.WithAutoCallDepth())
	l.Infof(format, v...)
}
