package logtrack

import (
	"github.com/dankozitza/dkutils"
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
	Log_file  string
	To_Stdout bool
}

var stat statdist.Stat = statdist.Stat{
	statdist.GetId(),
	"INIT",
	seestack.Short(),
	"package initialized",
	""}

func init() {
	statdist.Handle(stat, true)
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
func (l *LogTrack) Pv(v float64, msg ...interface{}) {

	conf := sconf.Inst()

	fix_lvl_range()

	if conf["logtrack_verbosity_level"].(float64) >= v {

		if conf["logtrack_verbosity_level"].(float64) >= 3 {
			prefix := "[" + seestack.ShortExclude(1) + "] "
			msg = append(msg, 0)
			copy(msg[1:], msg[0:])
			msg[0] = prefix
		}

		var file_path string
		if l.Log_file != "" {
			file_path = l.Log_file
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
	conf := sconf.Inst()

	cpy := conf["logtrack_default_log_file"]
	err := dkutils.ForceType(&cpy, "logtrack_"+seestack.LastFile()+".log")
	conf["logtrack_default_log_file"] = cpy

	if err != nil {
		stat.ShortStack = seestack.Short()
		stat.Status = "WARN"
		stat.Message = "log_file was set to the default: " +
			conf["logtrack_default_log_file"].(string)
		statdist.Handle(stat, false)
	}

	//if conf["logtrack_default_log_file"] == nil {

	//	// set the default log file if not set
	//	conf["logtrack_default_log_file"] = "logtrack_" +
	//		seestack.LastFile() + ".log"

	//	stat.ShortStack = seestack.Short()
	//	stat.Status = "WARN"
	//	stat.Message = "log_file was set to the default: " +
	//		conf["logtrack_default_log_file"].(string)
	//	statdist.Handle(stat, false)
	//}
	return
}

// fix_lvl_range
//
//   Ensures that logtrack_verbosity_level is defined and usable
//
func fix_lvl_range() {
	conf := sconf.Inst()

	cpy := conf["logtrack_verbosity_level"]
	err := dkutils.ForceType(&cpy, float64(3))
	conf["logtrack_verbosity_level"] = cpy

	if err != nil {
		stat.ShortStack = seestack.Short()
		stat.Status = "WARN"
		stat.Message = "logtrack_verbosity_level was set to the default value 3. "
		stat.Message += err.Error()
		statdist.Handle(stat, false)
	}

	if conf["logtrack_verbosity_level"].(float64) < 0 {
		stat.ShortStack = seestack.Short()
		stat.Status = "WARN"
		stat.Message = "logtrack_verbosity_level is out of range, "
		stat.Message += "defaulting to 0"
		conf["logtrack_verbosity_level"] = 0
		statdist.Handle(stat, false)

	} else if conf["logtrack_verbosity_level"].(float64) > 5 {
		stat.ShortStack = seestack.Short()
		stat.Status = "WARN"
		stat.Message = "logtrack_verbosity_level is out of range, "
		stat.Message += "defaulting to 5"
		conf["logtrack_verbosity_level"] = 5
		statdist.Handle(stat, false)
	}
	return
}
