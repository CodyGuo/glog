package glog

import (
	"bytes"
	"testing"
)

func TestWithLevel(t *testing.T) {
	var buf bytes.Buffer
	for _, testcase := range levelTests {
		t.Run(testcase.name, func(t *testing.T) {
			buf.Reset()
			l := New(&buf, WithLevel(testcase.level))
			if got := l.Level(); got.String() != testcase.want {
				t.Errorf("level %s: expected %s, got %s", testcase.level, testcase.want, got)
			}
		})
	}
}

func TestWithLevelLength(t *testing.T) {
	const want uint8 = 4
	var buf bytes.Buffer
	for _, testcase := range levelTests {
		t.Run(testcase.name, func(t *testing.T) {
			buf.Reset()
			l := New(&buf, WithLevelLength(4))
			if got := l.LevelLength(); got != want {
				t.Errorf("LevelLength %s: expected %d, got %d", testcase.name, want, got)
			}
		})
	}
}

func TestWithPrefix(t *testing.T) {
	want := "[testPrefix]"
	var buf bytes.Buffer
	l := New(&buf, WithPrefix("[testPrefix]"))
	if got := l.Prefix(); got != want {
		t.Errorf("prefix [testPrefix]: expected %s, got %s", want, got)
	}
}

func TestWithFlags(t *testing.T) {
	want := LglogFlags
	var buf bytes.Buffer
	l := New(&buf, WithFlags(LstdFlags|Lmicroseconds|Lshortfile|Lmsgprefix|Lmsglevel))
	if got := l.Flags(); got != want {
		t.Errorf("flag %d: expected: %d, got: %d",
			LstdFlags|Lmicroseconds|Lshortfile|Lmsgprefix, want, got)
	}
}

func TestWithCallDepth(t *testing.T) {
	want := 5
	var buf bytes.Buffer
	l := New(&buf, WithCallDepth(5))
	if got := l.CallDepth(); got != want {
		t.Errorf("calldepth 5: expected %d, got %d", want, got)
	}
}

func TestWithAutoCallDepth(t *testing.T) {
	want := 4
	var buf bytes.Buffer
	l := New(&buf, WithAutoCallDepth())
	l.AutoCallDepth()
	if got := l.CallDepth(); got != want {
		t.Errorf("autoCallDepth 5: expected %d, got %d", want, got)
	}
}
