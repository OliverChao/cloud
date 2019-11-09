package mysqlCon

import "time"
import "cloud/model"

const (
	MysqlUri = "root:toor@(127.0.0.1:3306)/sealaw?charset=utf8mb4&parseTime=True&loc=Local"
)

//mysql config struct
type MysqlConfig struct {
	MysqlUri string
	//MaxConnNUm int
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Models          []interface{}
}

//return MysqlConfig
//you can change this option to change mysql connection config
func LoadMysqlConfig() (IFMysqlCon *MysqlConfig) {

	IFMysqlCon = &MysqlConfig{
		MysqlUri:        MysqlUri,
		MaxIdleConns:    20,
		MaxOpenConns:    50,
		ConnMaxLifetime: time.Second * 4,
		Models: []interface{}{
			&model.Model{},
			&model.Admin{},
			&model.Article{},
		},
	}
	return
}
