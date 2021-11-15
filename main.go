package main

import (
	"net/http"
	"ssq-spider/logger"
	"ssq-spider/server"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		panic("Init logger error: " + err.Error())
	}

	go server.SSQSpider()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("敬请期待"))
	})

	logger.Logger.Info("ssqserver start ...")
	if err := http.ListenAndServe(":5623", nil); err != nil {
		panic("start server error: " + err.Error())
	}
}
