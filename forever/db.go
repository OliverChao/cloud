package forever

import (
	"cloud/config/mysqlCon"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

//db is the variable in this package
var db *gorm.DB

//connection mysql to make db can be used
func Connect(con *mysqlCon.MysqlConfig) {
	var err error
	db, err = gorm.Open("mysql", con.MysqlUri)
	if err != nil {
		logrus.Fatalf("[mysql] connect mysql error" + err.Error())

		return
	} else {
		logrus.Info("[mysql] connect successfully")
	}

	if err = db.AutoMigrate(con.Models...).Error; err != nil {
		logrus.Fatal("[mysql] auto migrate tables failed: " + err.Error())
	}

	db.DB().SetMaxIdleConns(con.MaxIdleConns)
	db.DB().SetMaxOpenConns(con.MaxOpenConns)
	db.DB().SetConnMaxLifetime(con.ConnMaxLifetime)
}

func DisconnectDB() {
	if err := db.Close(); err != nil {
		logrus.Errorf("[mysql] Disconnect failed: " + err.Error())
	} else {
		logrus.Info("[mysql] Disconnect successful")
	}
}

func GetGlobalGormDB() *gorm.DB {
	return db
}

func DropAll(con *mysqlCon.MysqlConfig) {
	db.DropTableIfExists(con.Models...)
}

//func CreateDemo() {
//	comment := model.Comment{
//		ArticleID: 1,
//		UserID:    1,
//		Content:   "yes~ you can",
//		ReplyID:   0,
//		PushedAt:  time.Now(),
//	}
//	article := model.Article{
//		PushedAt:     time.Now(),
//		Title:        "together forever",
//		Content:      "oliver loves annabelle forever~",
//		Tags:         "love",
//		Commentable:  true,
//		ViewCount:    0,
//		CommentCount: 0,
//		Path:         "",
//		Stared:       false,
//		Topped:       false,
//		UserID:       1,
//		Comments:     []model.Comment{comment},
//	}
//	//uuids, _ := uuid.NewV4()
//	article.UUID = strings.Replace(uuid.NewV4().String(), "-", "", -1)
//
//	user := model.User{
//		Name:              "oliver",
//		cloudKey:           "520annabelle",
//		Password:          "ce5ca673d13b36118d54a7cf13aeb0ca012383bf771e713421b4d1fd841f539a",
//		TotalArticleCount: 2,
//		Articles:          []model.Article{article},
//	}
//
//	db.Create(&user)
//
//}

//func QueryDemo() {
//	u := model.User{
//		Model: model.Model{ID: 1},
//	}
//	db.Model(&u).Find(&u)
//	db.Model(&u).Related(&u.Articles)
//	//db.Model(&u.Articles).Related(&u.Articles)
//
//	for _, v := range u.Articles {
//		db.Model(&v).Related(&v.Comments)
//		//fmt.Println(v.Comments)
//		for _, c := range v.Comments {
//			db.Model(&c).Related(&c.User)
//			fmt.Println("user:", c.User.Name)
//		}
//	}
//
//	bytes, _ := json.Marshal(u)
//	fmt.Println(string(bytes))
//
//}
