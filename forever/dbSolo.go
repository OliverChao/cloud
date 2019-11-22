package forever

import (
	"cloud/model"
	"errors"
	"github.com/sirupsen/logrus"
)

func AddKindToMysql(name string) {
	kind := &model.Kind{
		Name:  name,
		Count: 0,
	}
	db.Create(&kind)
	//db
}

//func WriteNewFileToMysql()

func addArticle(article *model.Article) {
	db.Create(article)
}

func VerifyUser(name, passwd string) (*model.Admin, error) {
	user := &model.Admin{
		Name:     name,
		Password: passwd,
	}

	db.Model(&user).Find(&user)
	if user.Password != passwd {
		logrus.Error("VerifyUser failed")
		return nil, errors.New("VerifyUser failed")
	} else {
		logrus.Info("VerifyUser successfully")
	}
	return user, nil

}

//func GetArticleByID(id int)model.Article{
//
//}
