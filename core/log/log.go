package log

import (
	"io"
	"log"
	"os"
	"regexp"
)

var (
	logger        *log.Logger
	logWriter     tLogWriter
	logFile       *os.File
	termColorCode *regexp.Regexp
)

type tLogWriter struct{}

func (tLogWriter) Write(p []byte) (n int, err error) {
	_, _ = os.Stdout.Write(p)

	p = termColorCode.ReplaceAll(p, []byte{})
	n, err = logFile.Write(p)

	return
}

func initLogger() {
	var err error
	logFile, err = os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	logger = log.New(logWriter, "", 0)
	termColorCode = regexp.MustCompile(`\[\d+(;\d+)?m`)
}

// Logger возвращает текущий логер приложения и инициализирует его
// при необходимости.
func Logger() *log.Logger {
	if logger == nil {
		initLogger()
	}
	return logger
}

// Writer возвращает io.Writer связанный с текущим логгером.
func Writer() io.Writer {
	if logger == nil {
		initLogger()
	}
	return logWriter
}

// Printf выводит строку с помощью логгера.
func Printf(format string, v ...interface{}) {
	Logger().Printf(format, v)
}
