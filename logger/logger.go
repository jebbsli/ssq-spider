package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"ssq-spider/configure"
)


var Logger *logrus.Logger

func checkFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

func InitLogger() error {
	var file *os.File
	var err error

	if !checkFileExist("log") {
		if err := os.Mkdir("log", 0666); err != nil {
			panic("os.Mkdir error: " + err.Error())
		}
	}

	filePath := filepath.Join("log", configure.GlobalConfig.Log.Filename)
	if !checkFileExist(filePath) {
		file, err = os.Create(filePath)
		if err != nil {
			panic("os.Create file error: " + err.Error())
		}
	} else {
		file, err = os.Open(filePath)
		if err != nil {
			panic("os.Open file error: " + err.Error())
		}
	}

	Logger = logrus.New()
	Logger.Out = file

	Logger.Info("Init log success ...")

	return nil
}
