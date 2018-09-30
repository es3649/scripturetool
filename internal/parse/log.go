package parse

import "github.com/sirupsen/logrus"

var log = logrus.New()

// SetVerbosity sets the verbosity of the logger
func SetVerbosity(lvl int) {
	switch {
	case lvl <= 0:
		log.SetLevel(logrus.FatalLevel)
	case lvl == 1:
		log.SetLevel(logrus.ErrorLevel)
	case lvl == 2:
		log.SetLevel(logrus.WarnLevel)
	case lvl == 3:
		log.SetLevel(logrus.InfoLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}
