package glog

import (
	"io"
	"os"
)

type Option func(*Logger)

func WithFile(name string, flag int, perm os.FileMode) Option {
	return func(l *Logger) {
		f, err := os.OpenFile(name, flag, perm)
		if err != nil {
			panic(err)
		}
		l.closers = append(l.closers, f)
		l.out = io.MultiWriter(l.out, f)
	}
}

func WithMultiWriter(writers ...io.Writer) Option {
	return func(l *Logger) {
		writers = append(writers, l.out)
		l.out = io.MultiWriter(writers...)
	}
}

func WithWriteCloser(writeCloser io.WriteCloser) Option {
	return func(l *Logger) {
		l.closers = append(l.closers, writeCloser)
		l.out = io.MultiWriter(l.out, writeCloser)
	}
}

func WithMultiWriteCloser(writeClosers ...io.WriteCloser) Option {
	return func(l *Logger) {
		for _, writeCloser := range writeClosers {
			l.closers = append(l.closers, writeCloser)
			l.out = io.MultiWriter(l.out, writeCloser)
		}
	}
}

func WithFlags(flag int) Option {
	return func(l *Logger) {
		l.flag = flag
	}
}

func WithPrefix(prefix string) Option {
	return func(l *Logger) {
		l.prefix = prefix
	}
}

func WithLevel(level Level) Option {
	return func(l *Logger) {
		l.level = level
	}
}

func WithLevelLength(levelLength uint8) Option {
	return func(l *Logger) {
		l.levelLength = levelLength
	}
}

func WithCallDepth(calldepth int) Option {
	return func(l *Logger) {
		l.callDepth = calldepth
	}
}

func WithAutoCallDepth() Option {
	return func(l *Logger) {
		l.AutoCallDepth()
	}
}
