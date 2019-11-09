package forever

func MysqlUnRegister() {
	//db.DropTable(&model.Model{}, &model.User{})
	DisconnectDB()

}

func RedisUnRegister() {
	DisConnectRedis()

}
