package logtrack

import (
	"github.com/dankozitza/logdist"
	"github.com/dankozitza/sconf"
	"github.com/dankozitza/seestack"
	"github.com/dankozitza/statdist"
)

// logtrack
//
// This package is used to manage the verbosity level and append/prepend
// information to log messages. Each message will be printed to the log file
// given and to stdout.
//
// the verbosity level can be set in the sconf settings under
// 'logtrack_verbosity_level'. It is an integer ranging from 0 to 5
//
//    verbosity levels:
//       5 - log everything
//       4 - log everything other than level 5 messages
//       3 - default, log messages that are expected to be seen
//       2 - do not log the stack in each message
//       1 - do not log anything other than level 0 and level 1 messages
//       0 - do not log anything
//
type LogTrack struct {
	log_file  string
	To_Stdout bool
}

var conf sconf.Sconf = sconf.Inst()
var stat statdist.Stat = statdist.Stat{
	statdist.GetId(),
	"INIT",
	seestack.Short(),
	"package initialized",
	""}

func init() {
	statdist.Handle(stat)
}

func (l *LogTrack) Set_log_file_path(file_path string) {
	l.log_file = file_path
	return
}

func New() LogTrack {

	l := LogTrack{"", true}

	fix_ldlf_path()
	fix_lvl_range()

	return l
}

// Pv - Print Verbose
//
// Prints the message if the logshare_verbosity_level is high enough.
//
//    v levels:
//
//       5 - low priority messages that are only printed when
//           logtrack_verbosity_level is set to 5
//       4 - ...
//       3 - normal priority messages
//       2 - messages that can be printed with no ShortStack
//       1 - ...
//       0 - high priority messages that will always be printed
//
func (l *LogTrack) Pv(v int, msg ...interface{}) {

	fix_lvl_range()

	if conf["logtrack_verbosity_level"].(int) >= v {

		if conf["logtrack_verbosity_level"].(int) >= 3 {
			prefix := "[" + seestack.ShortExclude(1) + "] "
			msg = append(msg, 0)
			copy(msg[1:], msg[0:])
			msg[0] = prefix
		}

		var file_path string
		if l.log_file != "" {
			file_path = l.log_file
		} else {
			fix_ldlf_path()
			file_path = conf["logtrack_default_log_file"].(string)
		}

		logdist.Message(file_path, l.To_Stdout, msg...)
	}
	return
}

// P - Print
//
// calls Pv with default verbose level 3
//
func (l *LogTrack) P(msg ...interface{}) {
	l.Pv(3, msg...)
	return
}

// fix_ldlf_path
//
// Ensures that logtrack_default_log_file is defined and usable
//
func fix_ldlf_path() {
	if conf["logtrack_default_log_file"] == nil {

		// set the default log file if not set
		conf["logtrack_default_log_file"] = "logtrack_" +
			seestack.LastFile() + ".log"

		stat.ShortStack = seestack.Short()
		stat.Status = "WARN"
		stat.Message = "log_file was set to the default: " +
			conf["logtrack_default_log_file"].(string)
		statdist.Handle(stat)
	}
	return
}

// fix_lvl_range
//
//   Ensures that logtrack_verbosity_level is defined and usable
//
func fix_lvl_range() {
	if conf["logtrack_verbosity_level"] == nil {
		conf["logtrack_verbosity_level"] = 3

	} else {
		if conf["logtrack_verbosity_level"].(int) < 0 {
			stat.ShortStack = seestack.Short()
			stat.Status = "WARN"
			stat.Message = "logtrack_verbosity_level is out of range, "
			stat.Message += "defaulting to 0"
			conf["logtrack_verbosity_level"] = 0
			statdist.Handle(stat)

		} else if conf["logtrack_verbosity_level"].(int) > 5 {
			stat.ShortStack = seestack.Short()
			stat.Status = "WARN"
			stat.Message = "logtrack_verbosity_level is out of range, "
			stat.Message += "defaulting to 5"
			conf["logtrack_verbosity_level"] = 5
			statdist.Handle(stat)
		}
	}
	return
}
