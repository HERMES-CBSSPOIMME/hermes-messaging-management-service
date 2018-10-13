package logruswrapper

import (
	io "io"

	logrus "github.com/sirupsen/logrus"
)

var (
	logrusLogger = logrus.New()
)

// Init a logger
// output parameter can be anything of type io.Writer such as a File or os.Stdout
// format parameter can be either JSON or TXT
func Init(output io.Writer, format string) {

	logrusLogger.Out = output

	switch format {
	case "JSON":
		logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	case "TXT":
		logrusLogger.SetFormatter(&logrus.TextFormatter{})
	default:
		logrusLogger.SetFormatter(&logrus.TextFormatter{})
	}
}

// Info log an information
// logEntryInfos parameter is a LogEntryInfos struct instance
func Info(logEntryInfos *LogEntryInfos) {

	logrusLogger.WithFields(logrus.Fields{
		"logID":    logEntryInfos.LogID,
		"session":  logEntryInfos.Session,
		"service":  logEntryInfos.Service,
		"endpoint": logEntryInfos.Endpoint,
		"code":     logEntryInfos.Code,
	}).Info(logEntryInfos.Message)
}

// Error log an error
// logEntryInfos parameter is a LogEntryInfos struct instance
func Error(logEntryInfos *LogEntryInfos) {

	logrusLogger.WithFields(logrus.Fields{
		"logID":    logEntryInfos.LogID,
		"session":  logEntryInfos.Session,
		"service":  logEntryInfos.Service,
		"endpoint": logEntryInfos.Endpoint,
		"code":     logEntryInfos.Code,
	}).Error(logEntryInfos.Message)
}
