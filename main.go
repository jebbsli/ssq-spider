package main

import "ssq-spider/logger"

func main() {
	if err := logger.InitLogger(); err != nil {
		panic("Init logger error: " + err.Error())
	}

	logger.Logger.Info("ssqserver start ...")
}
