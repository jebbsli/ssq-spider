package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"ssq-spider/configure"
	"ssq-spider/utils"
)

var Logger *logrus.Logger

func InitLogger() error {
	var err error

	Logger = logrus.New()

	logDir := filepath.Join(configure.RootDir, "log")

	if !utils.CheckFileExist(logDir) {
		if err := os.Mkdir(logDir, 0666); err != nil {
			panic("os.Mkdir error: " + err.Error())
		}
	}

	filePath := filepath.Join(logDir, configure.GlobalConfig.Log.Filename)
	if !utils.CheckFileExist(filePath) {
		file, err := os.Create(filePath)
		if err != nil {
			panic("os.Create file error: " + err.Error())
		}
		_ = file.Close()
	}

	Logger.Out, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic("os.Open file error: " + err.Error())
	}

	return nil
}
