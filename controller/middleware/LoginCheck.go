package middleware

import (
	"cloud/config/baseCon"
	"cloud/shadow"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func TestCookie(c *gin.Context){
	baseConfig := baseCon.LoadBaseConfig()
	ipAddress := baseConfig.IP
	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {
		// todo : redirect ... to login page
		//c.AbortWithStatusJSON(200, map[string]string{"msg": "token expired..need login..", "code": "-1"})

		//c.Request.URL.Path = "/auth/login"
		c.Redirect(http.StatusSeeOther, ipAddress+"/cloud/login")
		c.Abort()
		return
	} else {
		logrus.Info("[LoginCheck] get token...")
	}
	data, sign := shadow.UnParseToken(token)
	e := shadow.RsaVerify(data, sign)
	if e != nil {
		// todo : redirect ... to login page
		//c.AbortWithStatusJSON(200, map[string]string{"msg": "verify failed", "code": "-1"})
		c.Redirect(http.StatusSeeOther, ipAddress+"/cloud/login")
		c.Abort()
		return
	} else {
		logrus.Debug("login check correct")
	}
}


func LoginCheck(c *gin.Context) {
	baseConfig := baseCon.LoadBaseConfig()
	ipAddress := baseConfig.IP
	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {
		// todo : redirect ... to login page
		//c.AbortWithStatusJSON(200, map[string]string{"msg": "token expired..need login..", "code": "-1"})

		//c.Request.URL.Path = "/auth/login"
		c.Redirect(http.StatusSeeOther, ipAddress+"/cloud/login")
		c.Abort()
		return
	} else {
		logrus.Info("[LoginCheck] get token...")
	}
	data, sign := shadow.UnParseToken(token)
	e := shadow.RsaVerify(data, sign)
	if e != nil {
		// todo : redirect ... to login page
		//c.AbortWithStatusJSON(200, map[string]string{"msg": "verify failed", "code": "-1"})
		c.Redirect(http.StatusSeeOther, ipAddress+"/cloud/login")
		c.Abort()
		return
	} else {
		logrus.Debug("login check correct")
	}
	c.Next()
}
