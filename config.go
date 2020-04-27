package glog

type Config func(*Logger)

func WithLevel(level Level) Config {
	return func(l *Logger) {
		l.level = level
	}
}

func WithPrefix(prefix string) Config {
	return func(l *Logger) {
		l.prefix = prefix
	}
}

func WithFlags(flag int) Config {
	return func(l *Logger) {
		l.flag = flag
	}
}

func WithCallDepth(calldepth int) Config {
	return func(l *Logger) {
		l.calldepth = calldepth
	}
}
