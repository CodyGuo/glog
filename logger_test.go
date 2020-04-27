package glog

import (
	"bytes"
	"log"
	"testing"

	"github.com/sirupsen/logrus"
)

type callDepth struct {
	name      string
	calledpth int
	want      int
}

var calldepthes = []callDepth{
	{"calldepth 3", 3, 3},
	{"calldepth 4", 4, 4},
	{"calldepth 5", 5, 5},
}

func TestCallDepth(t *testing.T) {
	for _, testcase := range calldepthes {
		t.Run(testcase.name, func(t *testing.T) {
			SetCallDepth(testcase.calledpth)
			if got := CallDepth(); got != testcase.want {
				t.Errorf("calldepth: expected %d, got %d", testcase.want, got)
			}
		})
	}
}

func BenchmarkLogPrintln(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := log.New(&buf, "", log.LstdFlags)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkGLogInfo(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
}

func BenchmarkLogPrintlnNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := log.New(&buf, "", 0)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkGLogInfoNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf, WithFlags(0))
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
}

func BenchmarkLogRusInfo(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		logrus.Info(testString)
	}
}
