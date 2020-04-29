package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/CodyGuo/glog"
)

func main() {
	log := glog.New(os.Stderr, glog.WithFlags(glog.Lmsgjson|glog.Ldate|glog.Ltime|glog.Lmicroseconds|glog.Lmsgjson|glog.Lmsglevel))

	log.Info("glog: hello json")

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("logrus: hello json")
}
