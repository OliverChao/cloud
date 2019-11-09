package handlerFuncs

import (
	"cloud/echo"
	"cloud/forever"
	"cloud/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func DealToMarkdown(c *gin.Context) {

	//done : rewrite this return data
	ret := echo.NewRetResult()
	ret.Code = -1
	defer c.JSON(http.StatusOK, ret)
	markdownEngine, e := forever.GetMarkdownEngine()
	if e != nil {
		ret.Msg = "cannot get markdown engine"
		return
	}

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		ret.Msg = "post data format error"
		return
	}

	s, ok := arg["text"]
	if !ok {
		ret.Msg = "post data format error"
		return
	}
	mtext := s.(string)

	html, err := markdownEngine.MarkdownStr("", mtext)
	if err != nil {
		logrus.Error("[MarkdownStr err] " + err.Error())
		html = "MarkdownStr err"
		return
	}
	//successfully return data
	ret.Code = 1
	ret.Data = html

}

func PostArticle(c *gin.Context) {

	ret := echo.NewRetResult()
	ret.Code = -1
	defer c.JSON(http.StatusOK, ret)

	//arg := map[string]interface{}{}
	var arg echo.BindPostArticle
	if err := c.BindJSON(&arg); nil != err {
		ret.Msg = "post data format error"
		return
	}
	//logrus.Info(arg)
	article := model.Article{
		UUID:        strings.Replace(uuid.NewV4().String(), "-", "", -1),
		PushedAt:    time.Now(),
		Title:       arg.Title,
		Content:     arg.Content,
		Tags:        arg.Tags,
		Commentable: arg.Commentable,
		Stared:      arg.Stared,
		Topped:      arg.Topped,

		// Todo: 需要从session 的信息中获取 userid, 这里需要改
		UserID: 1,
	}
	forever.AddArticle(&article)

	ret.Code = 1
	ret.Msg = "success"

}

func TestRedis(c *gin.Context) {
	//c.JSON(http.StatusOK, gin.H{
	//	"name":"oliver",
	//})
	s := forever.GetArticleInfo()
	if s != "" {
		c.JSON(http.StatusOK, s)
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"err": "errorororororororororo",
			"msg": "oliver loves annabelle",
		})
	}

}

func TestLuteEngine(c *gin.Context) {
	markdownEngine, e := forever.GetMarkdownEngine()
	if e != nil {
		c.JSON(404, nil)

		return
	}
	bytes, err := ioutil.ReadAll(c.Request.Body)
	//logrus.Info(string(bytes))
	if err != nil {
		logrus.Error("[ioutil.ReadAll] " + err.Error())
	}
	html, e := markdownEngine.Markdown("", bytes)
	//logrus.Debug("html")
	c.JSON(200, string(html))
}

func GetAllTitles(c *gin.Context) {

}

func GetPosts(c *gin.Context) {
	name := c.Param("user")
	d, _ := c.Get("token_data")
	//sign, exists2 := c.Get("token_sign")
	data, ok := d.([]byte)
	if !ok {
		c.AbortWithStatusJSON(200, map[string]string{"msg": "token unparse error", "code": "-1"})
		return
	}
	var m map[string]interface{}
	_ = json.Unmarshal(data, &m)
	if name != m["name"].(string) {
		c.AbortWithStatusJSON(200, map[string]string{"msg": "token unparse error", "code": "-1"})
		return
	}
	// todo :add other verify methods
	//

	//page := c.DefaultQuery("page","1")
	//count := c.DefaultQuery("count","7")

	c.JSON(200, m)
}
