package logtrack

import (
	"fmt"
	"github.com/dankozitza/logdist"
	"github.com/dankozitza/sconf"
	"github.com/dankozitza/stattrack"
	"net/http"
	"syscall"
	"testing"
)

func TestNew(t *testing.T) {
	s := stattrack.New("TestNew: testing New function")

	lt := New()
	if (lt != LogTrack{"", true}) {
		s.Warn("TestNew: object returned from New() is not empty!")
		t.Fail()
	} else {
		s.Pass("TestNew: it worked!")
	}
}

func TestDefaultFilePath(t *testing.T) {
	s := stattrack.New("TestDefaultFilePath: testing default file path")
	lt := New()
	lt.P("test print")

	conf := sconf.Inst()

	if conf["logtrack_default_log_file"] != "logtrack_testing.log" {
		s.Warn("logtrack default log file is: " +
			conf["logtrack_default_log_file"].(string) +
			" it should be logtrack_testing.log")
		t.Fail()

	} else {
		s.Pass("default file path is correct")
	}
}

func TestPv(t *testing.T) {
	s := stattrack.New("TestPv: testing the print verbose function")
	lt := New()
	conf := sconf.Inst()
	lt.P("logtrack_verbosity_level = 3")
	conf["logtrack_verbosity_level"] = 3

	lt.Pv(0, "printing with verbose 0")
	lt.Pv(1, "printing with verbose 1")
	lt.Pv(2, "printing with verbose 2")
	lt.Pv(3, "printing with verbose 3")
	lt.Pv(4, "printing with verbose 4")
	lt.Pv(5, "printing with verbose 5")

	s.Pass("TestPv: did it work?")
}

func TestLogDist(t *testing.T) {
	ldh := logdist.HTTPHandler("stdout")
	http.Handle("/stdout", ldh)
	//http.ListenAndServe("localhost:9000", nil)
}

func TestClean(t *testing.T) {
	file_path := "logtrack_testing.log"
	fmt.Println("TestClean: removing", file_path)
	syscall.Exec("/usr/bin/rm", []string{"rm", file_path}, nil)
}
