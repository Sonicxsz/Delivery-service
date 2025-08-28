package logger

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
)

var (
	Log *Logger
)

type Logger struct {
	logger *logrus.Logger
	file   *os.File
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Infoln(args)
}
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warnln(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Errorln(args)
}

func (l *Logger) Close() {
	if l.file != nil {
		err := l.file.Close()
		log.Fatal(err)
	}
}

func newLogger() *Logger {
	return &Logger{
		logger: logrus.New(),
	}
}

func Init(level string, logDir string) error {
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	Log = newLogger()
	Log.logger.SetLevel(parsedLevel)
	Log.logger.SetFormatter(&logrus.TextFormatter{})

	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(logDir, "log.txt"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	Log.file = file
	Log.logger.SetOutput(file)

	return nil
}
