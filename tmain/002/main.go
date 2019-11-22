package main

import "cloud/forever"

func main() {
	forever.MysqlRegister()
	forever.MysqlUnRegister()
}
