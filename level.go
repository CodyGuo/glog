package glog

const (
	DEBUG Level = iota
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
)

var levelNames = []string{
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
		return "None"
	}
	return levelNames[l]
}
