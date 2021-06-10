package service

import "github.com/gin-gonic/gin"

func Svc() {
	r := gin.Default()
	r.GET("/ping", Pong)
	r.GET("/echo1", Echo)
	r.GET("/echo", ExecShell) //{"type":"input","input":"pwd\n"}
	r.Static("/static", "../static")
	r.Run("127.0.0.1:8082")
	r.Run()
}
