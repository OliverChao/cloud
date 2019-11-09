package main

import (
	"fmt"
	"cloud/config/mysqlCon"
)

func main() {
	//config := IFCon.BaseConfig{
	//	Host:     "127.0.0.1",
	//	MysqlUri: "",
	//	LogLevel: uint32(logrus.DebugLevel),
	//}
	//fmt.Println(config)
	IFMysqlCon := mysqlCon.LoadMysqlConfig()
	fmt.Println(IFMysqlCon.Models)

}
