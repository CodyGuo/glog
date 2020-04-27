package glog

type Config func(*Logger)

func WithLevel(level Level) Config {
	return func(l *Logger) {
		l.level = level
	}
}

func WithLevelLength(levelLength uint8) Config {
	return func(l *Logger) {
		l.levelLength = levelLength
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
