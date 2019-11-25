package middleware

import (
	"cloud/shadow"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoginCheck(c *gin.Context) {
	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {
		// todo : redirect ... to login page
		c.AbortWithStatusJSON(200, map[string]string{"msg": "token expired..need login..", "code": "-1"})
		return
	} else {
		logrus.Info("[LoginCheck] get token...")
	}
	data, sign := shadow.UnParseToken(token)
	e := shadow.RsaVerify(data, sign)
	if e != nil {
		// todo : redirect ... to login page
		c.AbortWithStatusJSON(200, map[string]string{"msg": "verify failed", "code": "-1"})
		return
	} else {
		logrus.Debug("login check correct")
	}
	c.Next()
}
