package handlerFuncs

import (
	"cloud/forever"
	"github.com/gin-gonic/gin"
)

func SealawIndex(c *gin.Context) {
	m := forever.GetKindsFromRedis()
	c.HTML(200, "list.html", gin.H{
		"options": m,
	})
}

//管理界面的list
func SealawIndex_2(c *gin.Context) {
	m := forever.GetKindsFromRedis()
	c.HTML(200, "list_parent.html", gin.H{
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
		c.HTML(200, "hint_nofile.html", gin.H{})
		m_2 := forever.GetKindsFromRedis()
		c.HTML(200, "list.html", gin.H{
			"options": m_2,
			"k":       kind,
		})
		return
	}
	c.HTML(200, "detail.html", gin.H{
		"options": m,
		"k":       kind,
	})
	//ret.Data = m
}

//管理界面的detail
func SealawKinds_2(c *gin.Context) {
	//ret := echo.NewRetResult()
	//ret.Code = -1
	//defer c.JSON(200, ret)
	kind := c.Param("kind")
	m, e := forever.GetAllArticleFromRedis(kind)
	if e != nil {
		//ret.Data = e.Error()
		// todo : redirect to index page
		c.HTML(200, "hint_nofile.html", gin.H{})
		m_2 := forever.GetKindsFromRedis()
		c.HTML(200, "list_parent.html", gin.H{
			"options": m_2,
			"k":       kind,
		})
		return
	}
	c.HTML(200, "detail_parent.html", gin.H{
		"options": m,
		"k":       kind,
	})
	//ret.Data = m
}

//func SealawSearch()
