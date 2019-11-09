package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
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

	engine.StaticFS("/resource",http.Dir("resource"))
	engine.Any("/", func(c *gin.Context) {
		c.JSON(200, "Yoo~~~ Hello~~~ iFei~~~")
	})

	api := engine.Group("/api")
	//api.Use(middleware.LoginCheck)

	api.GET("/", func(c *gin.Context) {
		c.String(200, "this is api group")
	})

	return engine
}
