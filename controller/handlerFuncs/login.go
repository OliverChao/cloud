package handlerFuncs

import (
	"cloud/config/baseCon"
	"cloud/forever"
	"cloud/shadow"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Passwd   string `form:"passwd" json:"passwd" binding:"required"`
}

func Login(c *gin.Context) {
	// ret := echo.NewRetResult()
	// ret.Code = -1
	// defer c.JSON(200, ret)
	var user User
	if e := c.MustBindWith(&user, binding.FormPost); e != nil {
		return
	}

	sha := sha256.New()
	sha.Write([]byte(user.Passwd))
	sum := sha.Sum(nil)
	passwdSha256 := hex.EncodeToString(sum)
	if e := forever.VerifyUser(user.Username, passwdSha256); e != nil {
		//弹窗
		c.HTML(200, "hint.html", gin.H{})
		c.HTML(200, "login.html", gin.H{})
		// ret.Msg = "login failed"
		return
	}
	dataBytes, _ := json.Marshal(user)
	sign, _ := shadow.RsaSign(dataBytes)
	token := fmt.Sprintf("%x.%x", dataBytes, sign)
	// todo : redirect to login successfully page

	session := sessions.Default(c)
	session.Set("token", token)
	_ = session.Save()

	c.HTML(200, "admin_index.html", gin.H{})
	// ret.Code = 1
	// ret.Data = user
	// ret.Msg = token
}

func ChangePasswdFun(c *gin.Context){
	oldPasswd := c.PostForm("oldPasswd")
	newPasswd := c.PostForm("newPasswd")

	var user User
	session := sessions.Default(c)
	token, _ := session.Get("token").(string)
	data, _ := shadow.UnParseToken(token)
	_ = json.Unmarshal(data, &user)
	if oldPasswd != user.Passwd || newPasswd == user.Passwd{
		c.String(200,"failed")
		return
	}

	sha := sha256.New()
	sha.Write([]byte(user.Passwd))
	sum := sha.Sum(nil)
	passwdSha256 := hex.EncodeToString(sum)
	forever.UpdateAdminPasswd(user.Username, passwdSha256)
	//
	baseConfig := baseCon.LoadBaseConfig()
	ipAddress := baseConfig.IP
	c.Redirect(http.StatusSeeOther, ipAddress+"/logout")
	//c.String(200,"successfully")
}
