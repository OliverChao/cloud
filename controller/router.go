package controller

import (
	"cloud/controller/handlerFuncs"
	"cloud/controller/middleware"
	"cloud/forever"
	"cloud/shadow"
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
	engine.LoadHTMLGlob("templates/*")
	engine.StaticFS("/resource", http.Dir("resource"))

	engine.Any("/", func(c *gin.Context) {
		c.JSON(200, "Yoo~~~ Hello~~~ iFei~~~")
	})

	engine.GET("/cloud/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	engine.POST("/cloud/login", handlerFuncs.Login)

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
	engine.POST("/upload", handlerFuncs.UploadOneFile)
	engine.POST("/multi/upload", handlerFuncs.UploadMultiFiles)

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

	engine.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		c.JSON(200, orSearch.Docs)
	})

	sealaw := engine.Group("/sealaw")
	sealaw.GET("/index", handlerFuncs.SealawIndex)

	//sealaw.GET("/index", func(c *gin.Context) {
	//	m := forever.GetKindsFromRedis()
	//	c.HTML(200, "index.html", gin.H{
	//		"options": m,
	//	})
	//})

	sealaw.GET("/kinds/:kind", handlerFuncs.SealawKinds)

	//sealaw.GET("/search",)

	api := engine.Group("/api")
	api.Use(middleware.LoginCheck)

	api.GET("/", func(c *gin.Context) {
		defer c.String(200, "this is api group")
		session := sessions.Default(c)
		get, ok := session.Get("token").(string)
		if !ok {
			return
		}

		logrus.Info(get)
		data, sign := shadow.UnParseToken(get)
		logrus.Info("[data] is ", string(data))
		err := shadow.RsaVerify(data, sign)
		if err != nil {
			logrus.Error("verify failed")
		} else {

		}
	})

	api.POST("/delete/article", handlerFuncs.DeleteArticle)

	api.POST("/delete/kind", handlerFuncs.DeleteKindFunc)

	return engine
}
