package handlerFuncs

import (
	"bufio"
	"cloud/forever"
	"cloud/utils"
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
	m := forever.GetAllKindFromRedis()

	form, _ := c.MultipartForm()

	files := form.File["mupload"]

	if len(files) == 0 {
		c.HTML(200, "hint_upload_none.html", gin.H{})
		c.HTML(http.StatusOK, "admin_onload.html", gin.H{
			"name":    "admin",
			"options": m,
		})
		return
	}

	if len(files) == 1 && !utils.UploadFilter(files[0].Filename) {
		c.HTML(200, "hint_upload_noexpand.html", gin.H{})
		c.HTML(http.StatusOK, "admin_onload.html", gin.H{
			"name":    "admin",
			"options": m,
		})
		return
	}

	kindName := c.PostForm("select")
	//logrus.Info(kindName)

	for _, file := range files {
		isvalid := utils.UploadFilter(file.Filename)
		//扩展名不符合要求
		if !isvalid {
			continue
			// return
		}
		//logrus.Info(file.Filename,"----",file.Size)
		open, e := file.Open()
		if e != nil {
			// ret.Msg = "get file hash failed"
			c.HTML(200, "hint_upload_fail.html", gin.H{})
			c.HTML(http.StatusOK, "admin_onload.html", gin.H{
				"name":    "admin",
				"options": m,
			})
			return
		}

		hashHex, filePath, topic := _fb(open, kindName, file.Filename)
		open.Close()

		// save size to redis
		forever.SaveSizeWhenUpload(filePath,file.Size)

		e = c.SaveUploadedFile(file, filePath)
		if e != nil {
			c.HTML(200, "hint_upload_fail.html", gin.H{})
			c.HTML(http.StatusOK, "admin_onload.html", gin.H{
				"name":    "admin",
				"options": m,
			})
			// ret.Msg = "upload failed"
		} else {
			forever.AddFileToKind(kindName, filePath, topic, hashHex, 1)
		}
		logrus.Info(file.Filename)
	}
	// ret.Code = 1
	// ret.Msg = "upload successfully"
	m2 := forever.GetAllKindFromRedis()
	c.HTML(200, "hint_upload.html", gin.H{})
	c.HTML(http.StatusOK, "admin_onload.html", gin.H{
		"name":    "admin",
		"options": m2,
	})
	return
}

func UploadOneFile(c *gin.Context) {
	// ret := echo.NewRetResult()
	// ret.Code = -1
	m := forever.GetAllKindFromRedis()

	file, err := c.FormFile("upload")
	if file == nil {

		//没有文件
		c.HTML(200, "hint_upload_none.html", gin.H{})
		c.HTML(http.StatusOK, "admin_onload.html", gin.H{
			"name":    "admin",
			"options": m,
		})
		// logrus.Info("file is none")
		return
	}

	isvalid := utils.UploadFilter(file.Filename)
	//扩展名不符合要求
	if !isvalid {
		c.HTML(200, "hint_upload_noexpand.html", gin.H{})
		c.HTML(http.StatusOK, "admin_onload.html", gin.H{
			"name":    "admin",
			"options": m,
		})
		return
	}

	kindName := c.PostForm("select")
	if err != nil {
		// ret.Msg = "no file upload"
		c.HTML(200, "hint_upload_none.html", gin.H{})
		c.HTML(http.StatusOK, "admin_onload.html", gin.H{
			"name":    "admin",
			"options": m,
		})
		return
	}

	//打开文件, 计算 hashcode
	open, e := file.Open()
	if e != nil {
		return
	}
	defer open.Close()

	hashHex, filePath, topic := _fb(open, kindName, file.Filename)

	// save size to redis
	forever.SaveSizeWhenUpload(filePath,file.Size)
	chanceNum := forever.CheckDoubleUploadFiles(kindName, topic, hashHex)

	switch chanceNum {
	case 2:
		// 2 哈希一样, 名字一样 => 直接 卡掉, 不执行操作.
		return
	case 1:
		// 1  哈希不一样, 名字一样 =>  先删除原来的,文件, 后正常上传,       1
		_ = forever.DeleteArticleFunc(topic, kindName)
	}

	e = c.SaveUploadedFile(file, filePath)
	if e != nil {
		// ret.Msg = "upload failed"
		c.HTML(200, "hint_upload_fail.html", gin.H{})
		c.HTML(http.StatusOK, "admin_onload.html", gin.H{
			"name":    "admin",
			"options": m,
		})
	} else {
		//ret.Data = hashHex
		// ret.Code = 1	sealaw.SealawKinds(
		// ret.Msg = "upload successfully"
		// ret.Data = filePath
		//	to//do : change mysql and redis
		forever.AddFileToKind(kindName, filePath, topic, hashHex, 1)

	}

	c.HTML(200, "hint_upload.html", gin.H{})
	c.HTML(http.StatusOK, "admin_onload.html", gin.H{
		"name":    "admin",
		"options": m,
	})

}

type DeleteData struct {
	Name string `form:"name" json:"name" binding:"required"`
	Kind string `form:"kind" json:"kind" binding:"required"`
}

func DeleteArticle(c *gin.Context) {
	// ret := echo.NewRetResult()
	// ret.Code = -1
	// defer c.JSON(200, ret)
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
		// ret.Msg = err.Error()
		return
	}
	logrus.Info(data)
	if err := forever.DeleteArticleFunc(data.Name, data.Kind); err != nil {
		// ret.Msg = err.Error()
		return
	}
	// ret.Msg = "dec.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	// ret.Code = 1)c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	//成功则执行以下，重新渲染页面
	kind := c.Param("kind")
	m_2 := forever.GetKindsFromRedis()
	c.HTML(200, "list_parent.html", gin.H{
		"options": m_2,
		"k":       kind,
	})
}

type deleteKind struct {
	Kind string `form:"kind" json:"kind" binding:"required"`
}

func DeleteKindFunc(c *gin.Context) {

	// ret := echo.NewRetResult()
	// ret.Code = -1
	// defer c.JSON(200, ret)
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
		// ret.Msg = err.Error()
		// c.HTML(200, "admin_createkind.html", gin.H{
		// 	"options": m,
		// })

		return
	}
	logrus.Info(data)
	if err := forever.DeleteKindFunc(data.Kind); err != nil {
		// ret.Msg = err.Error()
		return
	}
	// ret.Msg = "delete successfully"
	// ret.Code = 1
	m := forever.GetKindsFromRedis()
	// defer c.HTML(200, "hint_delete_succeful.html", gin.H{})
	defer c.HTML(200, "admin_createkind.html", gin.H{
		"options": m,
	})

}

func Show_index(c *gin.Context) {
	articles := forever.GetLastTenFilesInfo()
	numMap := forever.GetKindsFromRedis()

	// c.HTML(http.StatusOK, "admin_index.html", gin.H{})
	// logrus.Info(numMap)
	// c.JSON(200, numMap)
	// c.JSON(200, articles)

	c.HTML(200, "admin_index.html", gin.H{
		"articles": articles,
		"numMap":   numMap,
	})
}
