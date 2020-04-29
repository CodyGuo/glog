package glog

const (
	TRACE Level = iota
	DEBUG
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
	FATAL
)

var levelName = []string{
	"TRACE",
	"DEBUG",
	"INFO",
	"NOTICE",
	"WARNING",
	"ERROR",
	"CRITICAL",
	"FATAL",
}

type Level uint32

func (l Level) String() string {
	if l > FATAL {
		return "INVALID"
	}
	return levelName[l]
}

func (l Level) Len() uint8 {
	return uint8(len(l.String()))
}
