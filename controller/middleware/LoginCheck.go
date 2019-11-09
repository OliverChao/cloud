package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"iFei/shadow"
)

//var (
//	IFconBase = baseCon.LoadBaseConfig()
//	redirectUrl = IFconBase.Host+":"+string(IFconBase.Port)+"/ifei"
//)

func LoginCheck(c *gin.Context) {

	//session := sessions.Default(c)
	////session.Set("cotest","test~")
	////_ = session.Save()
	//token, ok := session.Get("token").(string)
	//if !ok {
	//	//c.AbortWithStatus(302)
	//	c.AbortWithStatusJSON(200, map[string]string{"msg": "token expired..need login..", "code": "-1"})
	//	return
	//} else {
	//	logrus.Info("[LoginCheck] get token...")
	//}
	//
	////split := strings.Split(token, ".")
	////data, _ := hex.DecodeString(split[0])
	////sign, _ := hex.DecodeString(split[1])
	//data, sign := shadow.UnParseToken(token)
	//e := shadow.RsaVerify(data, sign)
	//if e != nil {
	//	c.AbortWithStatusJSON(200, map[string]string{"msg": "verify failed", "code": "-1"})
	//	return
	//} else {
	//	logrus.Debug("login check correct")
	//}
	//c.Set("token_data", data)
	//c.Set("token_sign", sign)
	//c.Set("test", "test")
	//c.Next()
}
