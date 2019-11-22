package controller

import (
	"bufio"
	"cloud/echo"
	"cloud/forever"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func RegisterRouterMap() *gin.Engine {
	engine := gin.Default()
	//engine := gin.New()
	//engine.Use(gin.Recovery())
	//engine.Use(gin.Logger())
	//sessionRedis.NewStore()
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		Secure:   strings.HasPrefix("http://127.0.0.1", "https"),
		HttpOnly: true,
	})
	engine.Use(sessions.Sessions("law", store))

	engine.StaticFS("/resource", http.Dir("resource"))
	engine.Any("/", func(c *gin.Context) {
		c.JSON(200, "Yoo~~~ Hello~~~ iFei~~~")
	})

	engine.LoadHTMLGlob("templates/*")

	engine.GET("/upload", func(c *gin.Context) {
		m := forever.GetAllKindFromRedis()
		//logrus.Info(m)
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"name":    "oliver",
			"options": m,
		})
	})

	//engine.POST("/upload", func(c *gin.Context) {
	//	c.String(200, "get post request")
	//})

	engine.POST("/upload", func(c *gin.Context) {
		ret := echo.NewRetResult()
		ret.Code = -1
		defer c.JSON(http.StatusOK, ret)

		file, err := c.FormFile("upload")
		if err != nil {
			ret.Msg = "no file upload"
			return
		}
		open, e := file.Open()
		defer open.Close()
		if e != nil {
			ret.Msg = "get file hash error"
			return
		}
		reader := bufio.NewReader(open)
		sha := sha256.New()
		reader.WriteTo(sha)
		sum := sha.Sum(nil)

		hashHex := hex.EncodeToString(sum)
		ret.Data = hashHex

		//logrus.Info(c.PostForm("select"))
		kindName := c.PostForm("select")
		filePath := strings.Join([]string{"resource", kindName, file.Filename}, "/")
		split := strings.Split(file.Filename, ".")
		topic := strings.Join(split[:len(split)-1], ".")
		e = c.SaveUploadedFile(file, filePath)
		if e != nil {
			ret.Msg = "upload failed"
		} else {
			ret.Code = 1
			ret.Msg = "upload successfully"
			ret.Data = filePath
			//	todo : change mysql and redis
			forever.AddFileToKind(kindName, filePath, topic, hashHex)
		}
	})

	engine.POST("/multi/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload"]
		kindName := c.PostForm("select")
		logrus.Info(kindName)
		for _, file := range files {
			//filePath := strings.Join([]string{"resource",kindName,file.Filename},"/")
			//split := strings.Split(file.Filename, ".")
			//topic := strings.Join(split[:len(split)-1],".")
			//e := c.SaveUploadedFile(file, filePath)
			logrus.Info(file.Filename)
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, "%d files uploaded!", len(files))
	})

	engine.GET("/createkind", func(c *gin.Context) {
		//c.String(200,c.Request.Header.Get("Content-Type"))
		c.HTML(200, "test.html", gin.H{})
	})

	engine.POST("/createkind", func(c *gin.Context) {
		kind := c.PostForm("kind")
		if kind == "" {
			c.String(200, "asdasd")
			return
		}
		b := forever.AddNewKind(kind)
		if b {
			c.String(200, "OK")
		}
	})

	api := engine.Group("/api")
	//api.Use(middleware.LoginCheck)

	api.GET("/", func(c *gin.Context) {
		c.String(200, "this is api group")
	})

	api.GET("/getkinds", func(c *gin.Context) {
		s, e := forever.GetKindsFromRedis()
		if e != nil {
			c.String(200, "[getkinds] get redis error %v", e)
		} else {
			c.String(200, s)
		}
	})

	return engine
}
