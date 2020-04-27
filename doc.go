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

Output:
  2020/04/27 23:15:24 [INFO] hello info
  [customLog] 2020/04/27 23:15:24.391850 main.go:19: [TRAC] hello trace
  [customLog] 2020/04/27 23:15:24.391870 main.go:20: [DEBU] hello debug
  [customLog] 2020/04/27 23:15:24.391880 main.go:21: [INFO] hello info
  [customLog] 2020/04/27 23:15:24.391887 main.go:22: [NOTI] hello notice
  [customLog] 2020/04/27 23:15:24.391894 main.go:23: [WARN] hello warning
  [customLog] 2020/04/27 23:15:24.391904 main.go:24: [ERRO] hello error
  [customLog] 2020/04/27 23:15:24.391915 main.go:25: [CRIT] hello critical
For a full guide visit https://github.com/CodyGuo/glog
*/
package glog
