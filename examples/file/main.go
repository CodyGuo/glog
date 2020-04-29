package main

import (
	"os"

	"github.com/CodyGuo/glog"
)

func main() {
	log := glog.New(glog.Discard, glog.WithFile("glog.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644))
	defer log.Close()
	log.Debug("hello debug")
	log.Info("hello info")
	log.Notice("hello notice")
	log.Error("hello error")
}
