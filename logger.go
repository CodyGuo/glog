package glog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	Lmsglevel                     // log level: [INFO]
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
	LglogFlags    = LstdFlags | Lmicroseconds | Lshortfile | Lmsgprefix | Lmsglevel
)

const (
	currCallDepth = 3
)

var glog = New(os.Stdout, WithCallDepth(currCallDepth+2))

type Logger struct {
	once        *sync.Once
	mu          sync.Mutex
	level       Level
	levelLength uint8
	prefix      string
	flag        int
	calldepth   int
	buf         []byte
	stdLog      *log.Logger
}

func New(out io.Writer, config ...Config) *Logger {
	l := &Logger{
		once:        &sync.Once{},
		level:       INFO,
		levelLength: levelMaxLength,
		prefix:      "",
		flag:        LstdFlags,
	}
	// std log calldepth 2, The caller's calldepth needs to be increased by 1
	l.calldepth = currCallDepth + 1
	for _, c := range config {
		c(l)
	}
	l.stdLog = log.New(out, l.prefix, l.flag)
	return l
}

func (l *Logger) Trace(v ...interface{}) {
	l.log(TRACE, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(DEBUG, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.log(INFO, v...)
}

func (l *Logger) Notice(v ...interface{}) {
	l.log(NOTICE, v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.log(WARNING, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.log(ERROR, v...)
}

func (l *Logger) Critical(v ...interface{}) {
	l.log(CRITICAL, v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.logf(TRACE, format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(DEBUG, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(INFO, format, v...)
}
func (l *Logger) Noticef(format string, v ...interface{}) {
	l.logf(NOTICE, format, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.logf(WARNING, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.logf(CRITICAL, format, v...)
}

func (l *Logger) log(level Level, v ...interface{}) {
	l.output(level, fmt.Sprint(v...))
}

func (l *Logger) logf(level Level, format string, v ...interface{}) {
	l.output(level, fmt.Sprintf(format, v...))
}

func (l *Logger) formatHeader(buf *[]byte, level Level) {
	if l.flag&Lmsglevel != 0 {
		s := level.String()
		end := level.Len()
		if l.levelLength <= levelMaxLength && levelMinLength <= l.levelLength {
			if l.levelLength < end {
				end = l.levelLength
			}
			s = s[:end]
		}
		*buf = append(*buf, '[')
		*buf = append(*buf, s...)
		*buf = append(*buf, ']')
		*buf = append(*buf, ' ')
	}
}

func (l *Logger) output(level Level, s string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level > level {
		return nil
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, level)
	l.buf = append(l.buf, s...)
	return l.stdLog.Output(l.calldepth, string(l.buf))
}

func (l *Logger) Level() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) LevelLength() uint8 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.levelLength
}

func (l *Logger) SetLevelLength(length uint8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.levelLength = length
}

func (l *Logger) Prefix() string {
	return l.stdLog.Prefix()
}

func (l *Logger) SetPrefix(prefix string) {
	l.stdLog.SetPrefix(prefix)
}

func (l *Logger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
	l.stdLog.SetFlags(l.flag)
}

func (l *Logger) CallDepth() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.calldepth
}

func (l *Logger) SetCallDepth(calldepath int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calldepth = calldepath
}

func (l *Logger) AutoCallDepth() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.once.Do(func() {
		l.calldepth = l.calldepth + 1
	})
}

func (l *Logger) Output() io.Writer {
	return l.stdLog.Writer()
}

func (l *Logger) SetOutput(w io.Writer) {
	l.stdLog.SetOutput(w)
}

func Trace(v ...interface{}) {
	glog.Trace(v...)
}

func Debug(v ...interface{}) {
	glog.Debug(v...)
}

func Info(v ...interface{}) {
	glog.Info(v...)
}

func Notice(v ...interface{}) {
	glog.Notice(v...)
}

func Warning(v ...interface{}) {
	glog.Warning(v...)
}

func Error(v ...interface{}) {
	glog.Error(v...)
}

func Critical(v ...interface{}) {
	glog.Critical(v...)
}

func Tracef(format string, v ...interface{}) {
	glog.Tracef(format, v...)
}

func Debugf(format string, v ...interface{}) {
	glog.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	glog.Infof(format, v...)
}

func Noticef(format string, v ...interface{}) {
	glog.Noticef(format, v...)
}

func Warningf(format string, v ...interface{}) {
	glog.Warningf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	glog.Errorf(format, v...)
}

func Criticalf(format string, v ...interface{}) {
	glog.Criticalf(format, v...)
}

func GetLevel() Level {
	return glog.Level()
}

func SetLevel(level Level) {
	glog.SetLevel(level)
}

func LevelLength() uint8 {
	return glog.LevelLength()
}

func SetLevelLength(length uint8) {
	glog.SetLevelLength(length)
}

func Prefix() string {
	return glog.Prefix()
}

func SetPrefix(prefix string) {
	glog.SetPrefix(prefix)
}

func Flags() int {
	return glog.Flags()
}

func SetFlags(flag int) {
	glog.SetFlags(flag)
}

func CallDepth() int {
	return glog.CallDepth()
}

func SetCallDepth(calldepth int) {
	glog.SetCallDepth(calldepth)
}

func AutoCallDepth() {
	glog.AutoCallDepth()
}

func ResetCallDepth() {
	glog.SetCallDepth(currCallDepth + 2)
}
