package glog

const (
	TRACE Level = iota
	DEBUG
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
)

var (
	levelMinLength uint8 = 1
	levelMaxLength uint8 = CRITICAL.Len()
)

var levelName = []string{
	"TRACE",
	"DEBUG",
	"INFO",
	"NOTICE",
	"WARNING",
	"ERROR",
	"CRITICAL",
}

type Level uint32

func (l Level) String() string {
	if l > CRITICAL {
		return "INVALID"
	}
	return levelName[l]
}

func (l Level) Len() uint8 {
	return uint8(len(l.String()))
}
