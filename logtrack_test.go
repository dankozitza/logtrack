package logtrack

import (
	"testing"
	//"fmt"
	"github.com/dankozitza/logdist"
	"github.com/dankozitza/sconf"
	"github.com/dankozitza/stattrack"
	//"github.com/dankozitza/seestack"
	"net/http"
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

	lt.Pv("printing with verbose 0", 0)
	lt.Pv("printing with verbose 1", 1)
	lt.Pv("printing with verbose 2", 2)
	lt.Pv("printing with verbose 3", 3)
	lt.Pv("printing with verbose 4", 4)
	lt.Pv("printing with verbose 5", 5)

	//conf := sconf.Inst()

	s.Pass("TestPv: did it work?")
}

func TestLogDist(t *testing.T) {
	ldh := logdist.LogDistHandler("stdout")
	http.Handle("/", ldh)
	http.ListenAndServe("localhost:9000", nil)
}
