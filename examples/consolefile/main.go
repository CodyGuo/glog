package main

import (
	"os"

	"github.com/CodyGuo/glog"
)

func main() {
	fileOpt := glog.WithFile("glog.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	log := glog.New(os.Stderr, fileOpt, glog.WithFlags(glog.LglogFlags))
	defer log.Close()

	log.Debug("hello debug")
	log.Info("hello info")
	log.Notice("hello notice")
	log.Error("hello error")
}
