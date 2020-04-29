package glog

import "testing"

type levelTester struct {
	name  string
	level Level
	want  string
}

var levelTests = []levelTester{
	{"trace", TRACE, "TRACE"},
	{"debug", DEBUG, "DEBUG"},
	{"info", INFO, "INFO"},
	{"notice", NOTICE, "NOTICE"},
	{"warning", WARNING, "WARNING"},
	{"error", ERROR, "ERROR"},
	{"critical", CRITICAL, "CRITICAL"},
	{"fatal", FATAL, "FATAL"},
}

func TestLevel(t *testing.T) {
	for _, testcase := range levelTests {
		t.Run(testcase.name, func(t *testing.T) {
			SetLevel(testcase.level)
			got := GetLevel()
			if got.String() != testcase.want {
				t.Errorf("level: expected %s, got %s", testcase.want, got)
			}
		})
	}
}
