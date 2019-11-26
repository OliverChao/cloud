package handlerFuncs

import (
	"cloud/forever"
	"github.com/gin-gonic/gin"
)

func SealawIndex(c *gin.Context) {
	m := forever.GetKindsFromRedis()
	c.HTML(200, "index.html", gin.H{
		"options": m,
	})
}

func SealawKinds(c *gin.Context) {
	//ret := echo.NewRetResult()
	//ret.Code = -1
	//defer c.JSON(200, ret)
	kind := c.Param("kind")
	m, e := forever.GetAllArticleFromRedis(kind)
	if e != nil {
		//ret.Data = e.Error()
		// todo : redirect to index page
		return
	}
	c.HTML(200, "articles.html", gin.H{
		"options": m,
	})
	//ret.Data = m
}

//func SealawSearch()
