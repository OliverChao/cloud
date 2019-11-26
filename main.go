package main

import (
	"cloud/config/baseCon"
	"cloud/controller"
	"cloud/forever"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	//logrus.SetLevel(logrus.DebugLevel)
	//iFcon := baseCon.InitIFcon()
	logrus.SetLevel(logrus.DebugLevel)

}

func main() {
	IFconBase := baseCon.LoadBaseConfig()
	//logrus.SetLevel(IFconBase.LogLevel)
	forever.MysqlRegister()
	forever.RedisRegister()
	//forever.InitResourceDirs()
	forever.RiotRegister()
	forever.InitRiot()
	//gin.SetMode(gin.DebugMode)
	router := controller.RegisterRouterMap()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", IFconBase.Host, IFconBase.Port),
		Handler: router,
	}
	ExitServerHandler(server)
	if err := server.ListenAndServe(); err != nil {
		logrus.Errorf("serve listens failed: %v", err)
	}
}

func ExitServerHandler(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logrus.Infof("[signal %v]exiting IFei now...", s)
		if err := server.Close(); nil != err {
			logrus.Errorf("server close failed:%v", err)
		}

		// unregister here
		forever.MysqlUnRegister()
		forever.RedisUnRegister()
		forever.RiotUnregister()

		os.Exit(0)
	}()
}
