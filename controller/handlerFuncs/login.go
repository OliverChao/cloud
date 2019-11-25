package handlerFuncs

import (
	"cloud/echo"
	"cloud/forever"
	"cloud/shadow"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Passwd   string `form:"passwd" json:"passwd" binding:"required"`
}

func Login(c *gin.Context) {
	ret := echo.NewRetResult()
	ret.Code = -1
	defer c.JSON(200, ret)
	var user User
	if e := c.MustBindWith(&user, binding.FormPost); e != nil {
		return
	}

	sha := sha256.New()
	sha.Write([]byte(user.Passwd))
	sum := sha.Sum(nil)
	passwdSha256 := hex.EncodeToString(sum)
	if e := forever.VerifyUser(user.Username, passwdSha256); e != nil {
		ret.Msg = "login failed"
		return
	}
	dataBytes, _ := json.Marshal(user)
	sign, _ := shadow.RsaSign(dataBytes)
	token := fmt.Sprintf("%x.%x", dataBytes, sign)
	// todo : redirect to login successfully page

	session := sessions.Default(c)
	session.Set("token", token)
	_ = session.Save()

	ret.Code = 1
	ret.Data = user
	ret.Msg = token
}
