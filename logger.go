package glog

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"
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
	Lmsgjson                      // log json format: {"message":"hello json"}
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
	LglogFlags    = LstdFlags | Lmicroseconds | Lshortfile | Lmsgprefix | Lmsglevel
)

var glog = New(os.Stderr, WithCallDepth(4), WithFlags(LglogFlags))

var Discard io.Writer = ioutil.Discard

type Fields map[string]interface{}

type Logger struct {
	once        *sync.Once
	mu          sync.Mutex
	out         io.Writer
	closers     []io.Closer
	prefix      string
	flag        int
	callDepth   int
	level       Level
	levelLength uint8
	buf         []byte
}

func New(out io.Writer, options ...Option) *Logger {
	l := &Logger{
		out:       out,
		once:      &sync.Once{},
		prefix:    "",
		flag:      LstdFlags,
		callDepth: 3,
		level:     INFO,
	}
	for _, option := range options {
		option(l)
	}
	return l
}
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) AddOutput(writers ...io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	writers = append(writers, l.out)
	l.out = io.MultiWriter(writers...)
}

func (l *Logger) SetFile(name string, flag int, perm os.FileMode) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return err
	}
	l.closers = append(l.closers, f)
	l.out = f
	return nil
}

func (l *Logger) AddFile(name string, flag int, perm os.FileMode) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return err
	}
	l.closers = append(l.closers, f)
	l.out = io.MultiWriter(l.out, f)
	return nil
}

func (l *Logger) SetWriteCloser(writeCloser io.WriteCloser) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.closers = append(l.closers, writeCloser)
	l.out = writeCloser
}

func (l *Logger) AddWriteCloser(writeClosers ...io.WriteCloser) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, writeCloser := range writeClosers {
		l.closers = append(l.closers, writeCloser)
		l.out = io.MultiWriter(l.out, writeCloser)
	}
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.closers) == 0 {
		return os.ErrInvalid
	}
	var errs []error
	for _, closer := range l.closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%+v", errs)
	}
	return nil
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank and Lmsgprefix is unset),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided),
//   * l.prefix (if it's not blank and Lmsgprefix is set).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int, level Level) {
	if l.flag&Lmsgprefix == 0 {
		*buf = append(*buf, l.prefix...)
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
	if l.flag&Lmsgprefix != 0 {
		*buf = append(*buf, l.prefix...)
	}
	if l.flag&Lmsglevel != 0 {
		s := level.String()
		end := level.Len()
		if 0 < l.levelLength && l.levelLength < end {
			end = l.levelLength
			s = s[:end]
		}
		*buf = append(*buf, '[')
		*buf = append(*buf, s...)
		*buf = append(*buf, "] "...)
	}
}

func (l *Logger) jsonFormatHeader(buf *[]byte, t time.Time, file string, line int, level Level, s string) {
	var jsonData = struct {
		Time    string `json:"time,omitempty"`
		Level   string `json:"level,omitempty"`
		File    string `json:"file,omitempty"`
		Message string `json:"message"`
	}{}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
		}
		jsonData.Time = string(*buf)
		*buf = (*buf)[:0]
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		jsonData.File = string(*buf)
		*buf = (*buf)[:0]
	}
	if l.flag&Lmsglevel != 0 {
		s := level.String()
		end := level.Len()
		if 0 < l.levelLength && l.levelLength < end {
			end = l.levelLength
			s = s[:end]
		}
		*buf = append(*buf, s...)
		jsonData.Level = string(*buf)
		*buf = (*buf)[:0]
	}
	*buf = append(*buf, l.prefix...)
	*buf = append(*buf, s...)
	jsonData.Message = string(*buf)
	*buf = (*buf)[:0]

	jsonBytes, err := json.Marshal(&jsonData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "json format failed, error: %v\n", err)
		return
	}
	*buf = append(*buf, jsonBytes...)
}

func (l *Logger) Output(level Level, format string, v ...interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level > level {
		return nil
	}
	now := time.Now()
	var file string
	var line int
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(l.callDepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	var s string
	if format == "" {
		s = fmt.Sprint(v...)
	} else {
		s = fmt.Sprintf(format, v...)
	}
	if l.flag&Lmsgjson != 0 {
		l.jsonFormatHeader(&l.buf, now, file, line, level, s)
	} else {
		l.formatHeader(&l.buf, now, file, line, level)
		l.buf = append(l.buf, s...)
	}
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *Logger) log(level Level, v ...interface{}) {
	l.Output(level, "", v...)
}

func (l *Logger) logf(level Level, format string, v ...interface{}) {
	l.Output(level, format, v...)
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

func (l *Logger) Warn(v ...interface{}) {
	l.log(WARNING, v...)
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

func (l *Logger) Fatal(v ...interface{}) {
	l.log(FATAL, v...)
	l.Close()
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	l.log(PANIC, v...)
	l.Close()
	panic(fmt.Sprint(v...))
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

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(WARNING, format, v...)
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

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logf(FATAL, format, v...)
	l.Close()
	os.Exit(1)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.logf(PANIC, format, v...)
	l.Close()
	panic(fmt.Sprintf(format, v...))
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
}

func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
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

func (l *Logger) CallDepth() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.callDepth
}

func (l *Logger) SetCallDepth(calldepath int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.callDepth = calldepath
}

func (l *Logger) AutoCallDepth() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.once.Do(func() {
		l.callDepth = l.callDepth + 1
	})
}

func (l *Logger) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	glog.SetOutput(w)
}

func AddOutput(writers ...io.Writer) {
	glog.AddOutput(writers...)
}

func SetFile(name string, flag int, perm os.FileMode) error {
	return glog.SetFile(name, flag, perm)
}

func AddFile(name string, flag int, perm os.FileMode) error {
	return glog.AddFile(name, flag, perm)
}

func SetWriteCloser(writeCloser io.WriteCloser) {
	glog.SetWriteCloser(writeCloser)
}

func AddWriteCloser(writeClosers ...io.WriteCloser) {
	glog.AddWriteCloser(writeClosers...)
}

func Close() error {
	return glog.Close()
}

// Flags returns the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func Flags() int {
	return glog.Flags()
}

// SetFlags sets the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func SetFlags(flag int) {
	glog.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return glog.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	glog.SetPrefix(prefix)
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
	glog.SetCallDepth(4)
}

// Writer returns the output destination for the standard logger.
func Writer() io.Writer {
	return glog.Writer()
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

func Warn(v ...interface{}) {
	glog.Warn(v...)
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

func Fatal(v ...interface{}) {
	glog.Fatal(v...)
}

func Panic(v ...interface{}) {
	glog.Panic(v...)
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

func Warnf(format string, v ...interface{}) {
	glog.Warnf(format, v...)
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

func Fatalf(format string, v ...interface{}) {
	glog.Fatalf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	glog.Panicf(format, v...)
}

func Output(level Level, format string, v ...interface{}) error {
	return glog.Output(level, format, v...)
}
