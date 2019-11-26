package handlerFuncs

import (
	"bufio"
	"cloud/echo"
	"cloud/forever"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"strings"
)

func _fb(f multipart.File, kindName, filename string) (hashHex, filePath, topic string) {
	reader := bufio.NewReader(f)
	sha := sha256.New()
	_, _ = reader.WriteTo(sha)
	sum := sha.Sum(nil)

	hashHex = hex.EncodeToString(sum)

	//kindName := c.PostForm("select")
	filePath = strings.Join([]string{"resource", kindName, filename}, "/")

	split := strings.Split(filename, ".")
	if len(split) > 1 {
		topic = strings.Join(split[:len(split)-1], ".")
	} else {
		topic = filename
	}
	return
}

func UploadMultiFiles(c *gin.Context) {
	ret := echo.NewRetResult()
	ret.Code = -1
	defer c.JSON(http.StatusOK, ret)

	form, _ := c.MultipartForm()
	files := form.File["upload"]
	kindName := c.PostForm("select")
	//logrus.Info(kindName)

	for _, file := range files {
		open, e := file.Open()
		if e != nil {
			ret.Msg = "get file hash failed"
			return
		}
		hashHex, filePath, topic := _fb(open, kindName, file.Filename)
		open.Close()
		e = c.SaveUploadedFile(file, filePath)
		if e != nil {
			ret.Msg = "upload failed"
		} else {
			forever.AddFileToKind(kindName, filePath, topic, hashHex, 1)
		}
		logrus.Info(file.Filename)
	}
	ret.Code = 1
	ret.Msg = "upload successfully"
	return
}

func UploadOneFile(c *gin.Context) {
	ret := echo.NewRetResult()
	ret.Code = -1
	defer c.JSON(http.StatusOK, ret)

	file, err := c.FormFile("upload")
	kindName := c.PostForm("select")
	if err != nil {
		ret.Msg = "no file upload"
		return
	}

	//打开文件, 计算 hashcode
	open, e := file.Open()
	if e != nil {
		ret.Msg = "get file hash error"
		return
	}
	defer open.Close()

	hashHex, filePath, topic := _fb(open, kindName, file.Filename)

	e = c.SaveUploadedFile(file, filePath)
	if e != nil {
		ret.Msg = "upload failed"
	} else {
		//ret.Data = hashHex
		ret.Code = 1
		ret.Msg = "upload successfully"
		ret.Data = filePath
		//	to//do : change mysql and redis
		forever.AddFileToKind(kindName, filePath, topic, hashHex, 1)
	}
}

type DeleteData struct {
	Name string `form:"name" json:"name" binding:"required"`
	Kind string `form:"kind" json:"kind" binding:"required"`
}

func DeleteArticle(c *gin.Context) {
	ret := echo.NewRetResult()
	ret.Code = -1
	defer c.JSON(200, ret)
	var data DeleteData
	var err error
	contentType := c.Request.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		err = c.BindJSON(&data)
	case "application/x-www-form-urlencoded":
		err = c.MustBindWith(&data, binding.FormPost)
	}
	if err != nil {
		logrus.Error(err)
		ret.Msg = err.Error()
		return
	}
	logrus.Info(data)
	if err := forever.DeleteArticleFunc(data.Name, data.Kind); err != nil {
		ret.Msg = err.Error()
		return
	}
	ret.Msg = "delete successfully"
	ret.Code = 1
}

type deleteKind struct {
	Kind string `form:"kind" json:"kind" binding:"required"`
}

func DeleteKindFunc(c *gin.Context) {
	ret := echo.NewRetResult()
	ret.Code = -1
	defer c.JSON(200, ret)
	var data deleteKind
	var err error
	contentType := c.Request.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		err = c.BindJSON(&data)
	case "application/x-www-form-urlencoded":
		err = c.MustBindWith(&data, binding.FormPost)
	}

	if err != nil {
		logrus.Error(err)
		ret.Msg = err.Error()
		return
	}
	logrus.Info(data)
	if err := forever.DeleteKindFunc(data.Kind); err != nil {
		ret.Msg = err.Error()
		return
	}
	ret.Msg = "delete successfully"
	ret.Code = 1

}
