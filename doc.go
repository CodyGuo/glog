/*
Package glog implements a log infrastructure for Go:

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

Output:
  2020/04/27 17:12:21 [INFO] hello info
  [customLog] 2020/04/27 17:12:21 [DEBUG] hello debug
  [customLog] 2020/04/27 17:12:21 [INFO] hello info
For a full guide visit https://github.com/CodyGuo/glog
*/
package glog
