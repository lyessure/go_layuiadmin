package router

import (
	"layuiadmin/dao"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	Setupadminrouter(r)
	return r
}

func Validsession(c *gin.Context) {
	var session string

	if c.Request.Method == http.MethodGet {
		session = c.Query("access_token")
	} else {
		session = c.PostForm("access_token")
	}
	if session == "" {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "未登录",
			"data": nil,
		})
		c.Abort()
		return
	}
	sql := "select uid,sessiontime,isadmin,username from users where session=?"
	var userinfo struct {
		Uid         string `db:"uid"`
		Sessiontime int    `db:"sessiontime"`
		Isadmin     int    `db:"isadmin"`
		Username    string `db:"username"`
	}
	err := dao.Db.Get(&userinfo, sql, session)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "未登录",
			"data": nil,
		})
		c.Abort()
		return
	}
	if userinfo.Sessiontime+30*60 < int(time.Now().Unix()) { //session有效期30分钟
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "未登录",
			"data": nil,
		})
		c.Abort()
		return
	}
	sql = "update users set sessiontime=? where uid=?"
	dao.Db.Exec(sql, int(time.Now().Unix()), userinfo.Uid)
	c.Set("uid", userinfo.Uid)
	c.Set("isadmin", userinfo.Isadmin)
	c.Set("username", userinfo.Username)
	c.Next()
}
