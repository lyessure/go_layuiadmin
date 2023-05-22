package admin

import (
	"fmt"
	"layuiadmin/dao"
	"layuiadmin/util"
	"time"

	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type subMenuItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Jump  string `json:"jump"`
}

type MenuItem struct {
	Title string        `json:"title"`
	Icon  string        `json:"icon"`
	List  []subMenuItem `json:"list"`
}

func Getmenuitem(c *gin.Context) {
	uid := c.GetString("uid")
	isadmin := c.GetInt("isadmin")
	fmt.Println("uid:", uid, "isadmin:", isadmin)
	var menus []MenuItem

	var amenu MenuItem

	if isadmin == 1 {
		amenu = MenuItem{
			Title: "系统管理",
			Icon:  "layui-icon-home",
			List: []subMenuItem{
				{Name: "SubItem 1", Title: "帐号列表", Jump: "/user/administrators/list"},
			},
		}
		menus = append(menus, amenu)
	}

	amenu = MenuItem{
		Title: "订单管理",
		Icon:  "layui-icon-home",
		List: []subMenuItem{
			{Name: "SubItem 2", Title: "查看订单", Jump: "/app/workorder/list"},
		},
	}
	menus = append(menus, amenu)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": menus,
	})
}

func Userlogin(c *gin.Context) {
	//code=1001失败
	username := c.PostForm("username")
	password := c.PostForm("password")
	hash := md5.Sum([]byte(password))
	md5password := hex.EncodeToString(hash[:])
	var i int
	sql := "select 1 from users where username=? and password=?"
	err := dao.Db.Get(&i, sql, username, md5password)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "用户名或者密码错误",
			"data": nil,
		})
		return
	}
	session := util.Getrandomstr(32)
	sql = "update users set session=?,sessiontime=? where username=?"
	dao.Db.Exec(sql, session, int(time.Now().Unix()), username)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": gin.H{
			"access_token": session,
		},
	})
}

func Userlogout(c *gin.Context) {
	sql := "update users set session='' where uid=?"
	dao.Db.Exec(sql, c.GetString("uid"))
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "退出成功",
		"data": nil,
	})
}

func Userdetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"role":     c.GetInt("isadmin"),
			"username": c.GetString("username"),
			"sex":      "妖",
		},
	})
}

func Changepassword(c *gin.Context) {
	oldpassword := c.PostForm("oldPassword")
	password := c.PostForm("password")
	if len(password) < 6 || len(oldpassword) < 6 {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "密码长度不能小于6位",
			"data": nil,
		})
		return
	}
	hash := md5.Sum([]byte(oldpassword))
	md5oldpassword := hex.EncodeToString(hash[:])
	hash = md5.Sum([]byte(password))
	md5password := hex.EncodeToString(hash[:])
	sql := "select 1 from users where uid=? and password=?"
	var i int
	err := dao.Db.Get(&i, sql, c.GetString("uid"), md5oldpassword)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "原密码错误",
			"data": nil,
		})
		return
	}
	sql = "update users set password=? where uid=?"
	dao.Db.Exec(sql, md5password, c.GetString("uid"))
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "修改成功",
		"data": nil,
	})
}

func Adminlist(c *gin.Context) {
	var users []struct {
		Uid      string `db:"uid" json:"id"`
		Username string `db:"username" json:"username"`
		Isadmin  int    `db:"isadmin" json:"isadmin"`
		Cadmin   string `json:"cadmin"`
		Realname string `db:"realname" json:"realname"`
	}
	sql := "select uid,username,isadmin,realname from users"
	err := dao.Db.Select(&users, sql)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "没有记录",
			"data": nil,
		})
		return
	}
	for i := 0; i < len(users); i++ {
		if users[i].Isadmin == 1 {
			users[i].Cadmin = "超级管理员"
		} else {
			users[i].Cadmin = "普通用户"
		}
	}
	var count int
	sql = "select count(*) from users"
	dao.Db.Get(&count, sql)

	c.JSON(200, gin.H{
		"code":  0,
		"count": count,
		"msg":   "success",
		"data":  users,
	})
}

func Addadmin(c *gin.Context) {
	if c.GetInt("isadmin") != 1 {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "没有权限",
			"data": nil,
		})
		return
	}
	username := c.PostForm("username")
	password := c.PostForm("password")
	realname := c.PostForm("realname")
	isadmin := c.PostForm("isadmin")
	if len(password) < 6 {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "密码长度不能小于6位",
			"data": nil,
		})
		return
	}
	hash := md5.Sum([]byte(password))
	md5password := hex.EncodeToString(hash[:])
	sql := "select 1 from users where username=?"
	var i int
	err := dao.Db.Get(&i, sql, username)
	if err == nil {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "用户名已存在",
			"data": nil,
		})
		return
	}
	sql = "insert into users(username,password,realname,isadmin) values(?,?,?,?)"
	dao.Db.Exec(sql, username, md5password, realname, isadmin)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "添加成功",
		"data": nil,
	})
}

func Editadmin(c *gin.Context) {
	if c.GetInt("isadmin") != 1 {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "没有权限",
			"data": nil,
		})
		return
	}
	uid := c.PostForm("id")
	realname := c.PostForm("realname")
	isadmin := c.PostForm("isadmin")
	password := c.PostForm("password")
	if len(password) > 0 && len(password) < 6 {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "密码长度不能小于6位",
			"data": nil,
		})
		return
	}
	if len(password) > 0 {
		hash := md5.Sum([]byte(password))
		md5password := hex.EncodeToString(hash[:])
		sql := "update users set password=? where uid=?"
		dao.Db.Exec(sql, md5password, uid)
	}
	sql := "update users set realname=?,isadmin=? where uid=?"
	dao.Db.Exec(sql, realname, isadmin, uid)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "修改成功",
		"data": nil,
	})
}

func Rmadmin(c *gin.Context) {
	if c.GetInt("isadmin") != 1 {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "没有权限",
			"data": nil,
		})
		return
	}
	uid := c.PostForm("id")
	if uid == c.GetString("uid") {
		c.JSON(200, gin.H{
			"code": 4001,
			"msg":  "不能删除自己",
			"data": nil,
		})
		return
	}
	sql := "delete from users where uid=?"
	dao.Db.Exec(sql, uid)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除成功",
		"data": nil,
	})
}
