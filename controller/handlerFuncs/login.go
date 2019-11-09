package handlerFuncs

//todo : add Code Enum
//current:
//{
//	1 : successfully verify & set cookie
//	-1: post data format error
//	-2 : verify failed
//	2 : verify successfully but set cookie failed
//}
//const (
//	expireUser    = 14400 * time.Second
//	expireArticle = 14400 * time.Second
//)
//
//func LoginIn(c *gin.Context) {
//	//ret := echo.BindPostArticle{}
//	ret := echo.NewRetResult()
//	ret.Code = -1
//	defer c.JSON(http.StatusOK, ret)
//
//	//todo :add RSA decrypt if encrypt by client
//
//	var arg echo.BindLogin
//	if err := c.BindJSON(&arg); err != nil {
//		ret.Msg = "post data format error"
//		return
//	}
//	//
//	sha := sha256.New()
//	sha.Write([]byte(arg.Password))
//	sum := sha.Sum(nil)
//	passwd := hex.EncodeToString(sum)
//
//	user, e := forever.VerifyUser(arg.Username, passwd)
//	if e != nil {
//		ret.Code = -2
//		ret.Msg = "verify failed"
//		return
//	}
//
//	// to//do: // currently redis store article content, but maybe store article number can more  effectively.
//	go RedisSaveUserInfo(user)
//
//	// return some information about this user
//	//ret.Data = user.RetData()
//	dataMap := user.RetData()
//	dataMap["gen_time"] = time.Now().Unix()
//	//dataMap["expire"] =
//	byteData, _ := json.Marshal(dataMap)
//	sign, _ := shadow.RsaSign(byteData)
//	token := fmt.Sprintf("%x.%x", byteData, sign)
//
//	ret.Data = token
//	ret.Msg = "verify successfully"
//
//	session := sessions.Default(c)
//	session.Set("token", token)
//	e = session.Save()
//
//	if e != nil {
//		logrus.Error("[LoginIn] cookie save error")
//		ret.Code = 2
//		return
//	} else {
//		logrus.Info("[LoginIn] cookie save successfully")
//	}
//
//	ret.Code = 1
//
//}
//
//func GetUserArticleInfo(user *model.User) []model.Article {
//	db := forever.GetGlobalGormDB()
//	// user already has been found, so there use Related
//	//todo :add error handle
//	db.Model(&user).Related(&user.Articles)
//	//logrus.Info("[RedisSetUserInfo]",user.Articles[0].Title)
//	if len(user.Articles) > 10 {
//		return user.Articles[:10]
//
//	} else {
//		return user.Articles
//	}
//}
//
//func RedisSaveUserInfo(user *model.User) {
//	client := forever.GetGlobalRedisClient()
//
//	keyName := "user:" + user.Name + ":info"
//	exists := client.Exists(keyName)
//	if exists.Val() == 1 {
//		logrus.Info("redis has saved ", keyName)
//		return
//	}
//
//	//redis for user info
//	retStructData := &echo.UserInfoInRedis{
//		Name:              user.Name,
//		TotalArticleCount: user.TotalArticleCount,
//		ID:                user.ID,
//	}
//	retData := utils.Struct2Map(retStructData)
//	err := client.HMSet(keyName, retData).Err()
//	if err != nil {
//		logrus.Error("[RedisSaveUserInfo]", err)
//	} else {
//		logrus.Debug("[RedisSaveUserInfo] successfully")
//	}
//	client.Expire(keyName, expireUser)
//
//	//redis for article
//	articles := GetUserArticleInfo(user)
//	pipeline := client.Pipeline()
//	for _, v := range articles {
//		m := v.RetArticle()
//		key := fmt.Sprintf("article:%v:%v", v.UserID, v.ID)
//		//client.HMSet(key,m)
//		pipeline.HMSet(key, m)
//		pipeline.Expire(key, expireArticle)
//	}
//
//	_, err = pipeline.Exec()
//	err = pipeline.Close()
//
//	//marshalData, _ := json.Marshal(articles)
//}
