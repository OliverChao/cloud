package controller

import (
	"cloud/config/baseCon"
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
	//engine := gin.Default()
	baseConf := baseCon.LoadBaseConfig()
	ipAddress := baseConf.IP
	engine := gin.New()

	// logger ..
	engine.Use(gin.Logger())

	engine.Use(gin.Recovery())
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
	engine.Static("/pan", "./pan")
	engine.StaticFS("/resource", http.Dir("resource"))

	engine.GET("/", func(c *gin.Context) {
		//c.Request.URL.Path = "/sealaw/list"
		//engine.HandleContext(c)
		c.Redirect(http.StatusSeeOther, "/sealaw/list")
		return
	})

	engine.GET("/admin", func(c *gin.Context) {
		//c.Request.URL.Path = "/auth/admin_index"
		//engine.HandleContext(c)
		c.Redirect(http.StatusSeeOther, ipAddress+"/auth/admin_index")
		c.Abort()
		return
	})

	// logout router
	engine.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Options(sessions.Options{
			Path:   "/",
			MaxAge: -1,
		})
		session.Clear()
		if err := session.Save(); nil != err {
			logrus.Errorf("saves session failed: " + err.Error())
		}
		//	`调到登`入页
		c.Redirect(http.StatusSeeOther,ipAddress+"/cloud/login")
		c.Abort()
		return
	})

	engine.GET("/connect_me", func(c *gin.Context) {
		c.HTML(http.StatusOK, "connect_me.html", gin.H{})
	})

	engine.GET("/cloud/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	engine.POST("/cloud/login", handlerFuncs.Login)

	engine.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		c.JSON(200, orSearch.Docs)
	})

	engine.POST("/search", func(c *gin.Context) {
		query := c.PostForm("query")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		docs := orSearch.Docs
		m := forever.GetSearchResultFromRiot(docs)
		// c.JSON(200, orSearch.Docs)
		c.HTML(200, "detail.html", gin.H{
			"options": m,
		})
	})
	engine.GET("/search_parent", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		c.JSON(200, orSearch.Docs)
	})

	engine.POST("/search_parent", func(c *gin.Context) {
		query := c.PostForm("query")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		docs := orSearch.Docs
		m := forever.GetSearchResultFromRiot(docs)
		// c.JSON(200, orSearch.Docs)
		c.HTML(200, "detail_parent.html", gin.H{
			"options": m,
		})
	})

	auth := engine.Group("/auth")
	auth.Use(middleware.LoginCheck)


	auth.Any("/change/pass",handlerFuncs.ChangePasswdFun)

	//**************************
	// 数据库,
	auth.GET("/test_db", func(c *gin.Context) {
		articles := forever.GetLastTenFilesInfo()
		c.JSON(200, articles)
	})

	//********************

	auth.GET("/upload", func(c *gin.Context) {
		m := forever.GetAllKindFromRedis()
		//logrus.Info(m)
		c.HTML(http.StatusOK, "admin_onload.html", gin.H{
			"name":    "admin",
			"options": m,
		})
	})

	// 渲染页面
	//auth.GET("/", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "admin_index.html", gin.H{})
	//})
	auth.GET("/admin_index", handlerFuncs.Show_index)

	auth.GET("/admin_del", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_del.html", gin.H{})
	})
	auth.GET("/download", func(c *gin.Context) {
		c.HTML(200, "admin_download.html", gin.H{})
	})

	//engine.POST("/upload", func(c *gin.Context) {
	//	c.String(200, "get post request")
	//})
	auth.POST("/upload", handlerFuncs.UploadOneFile)
	auth.POST("/multi/upload", handlerFuncs.UploadMultiFiles)

	auth.GET("/createkind", func(c *gin.Context) {
		//c.String(200,c.Request.Header.Get("Content-Type"))
		m := forever.GetKindsFromRedis()
		c.HTML(200, "admin_createkind.html", gin.H{
			"options": m,
		})
	})

	auth.POST("/createkind", func(c *gin.Context) {
		kind := c.PostForm("kind")
		if kind == "" {
			c.String(200, "asdasd")
			return
		}
		b := forever.AddNewKind(kind)
		if b {
			m := forever.GetKindsFromRedis()
			c.HTML(200, "admin_createkind.html", gin.H{
				"options": m,
			})
		}
	})

	sealaw := engine.Group("/sealaw")
	sealaw.GET("/list", handlerFuncs.SealawIndex)
	sealaw.GET("/list_parent", handlerFuncs.SealawIndex_2)
	sealaw.GET("/docs/:kind/:article", func(c *gin.Context) {
		kind := c.Param("kind")
		article := c.Param("article")

		path := forever.GetOneArticleFromRedisByTitle(kind, article)
		s, e := forever.CheckDocAndSave(path)
		if e != nil {
			// doc 跳到
			c.HTML(200, "show_detail_doc.html", gin.H{
				"path": path,
				"k":    kind,
				"art":  article,
			})
			// c.String(200, s)
		} else {
			//pdf  docx txt -----
			c.HTML(200, "show_detail.html", gin.H{
				"text_s": s,
				"k":      kind,
				"art":    article,
			})
			// c.String(200, s)
		}
	})

	sealaw.POST("/list", func(c *gin.Context) {
		query := c.PostForm("query")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		docs := orSearch.Docs
		m := forever.GetSearchResultFromRiot(docs)
		// c.JSON(200, orSearch.Docs)
		c.HTML(200, "detail.html", gin.H{
			"options": m,
		})
	})

	//管理界面的list表POSt
	sealaw.POST("/list_praent", func(c *gin.Context) {
		query := c.PostForm("query")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		docs := orSearch.Docs
		m := forever.GetSearchResultFromRiot(docs)
		// c.JSON(200, orSearch.Docs)
		c.HTML(200, "detail_parent.html", gin.H{
			"options": m,
		})
	})

	sealaw.GET("/show_detail", func(c *gin.Context) {
		c.HTML(http.StatusOK, "show_detail.html", gin.H{
			"name": "oliver",
		})
	})

	//sealaw.GET("/index", func(c *gin.Context) {
	//	m := forever.GetKindsFromRedis()
	//	c.HTML(200, "index.html", gin.H{
	//		"options": m,
	//	})
	//})

	sealaw.GET("/kinds/:kind", handlerFuncs.SealawKinds)

	sealaw.POST("/kinds/:kind", func(c *gin.Context) {
		query := c.PostForm("query")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		docs := orSearch.Docs
		m := forever.GetSearchResultFromRiot(docs)
		// c.JSON(200, orSearch.Docs)
		c.HTML(200, "detail.html", gin.H{
			"options": m,
		})
	})
	//管理界面的kinds表POST

	sealaw.GET("/kinds_parent/:kind", handlerFuncs.SealawKinds_2)

	//管理界面的kinds表POST
	sealaw.POST("/kinds_parent/:kind", func(c *gin.Context) {
		query := c.PostForm("query")
		if query == "" {
			//	返回所有article
			return
		}
		orSearch := forever.OrSearch(query)
		docs := orSearch.Docs
		m := forever.GetSearchResultFromRiot(docs)
		// c.JSON(200, orSearch.Docs)
		c.HTML(200, "detail_parent.html", gin.H{
			"options": m,
		})
	})

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
