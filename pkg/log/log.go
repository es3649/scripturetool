package log

import "github.com/sirupsen/logrus"

// Log is the logger we will sue for the whole package
var Log = logrus.New()

// SetVerbosity sets the verbosity of the logger
func SetVerbosity(lvl int) {
	switch {
	case lvl <= 0:
		Log.SetLevel(logrus.FatalLevel)
	case lvl == 1:
		Log.SetLevel(logrus.ErrorLevel)
	case lvl == 2:
		Log.SetLevel(logrus.WarnLevel)
	case lvl == 3:
		Log.SetLevel(logrus.InfoLevel)
	default:
		Log.SetLevel(logrus.DebugLevel)
	}
}
