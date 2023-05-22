package router

import (
	"layuiadmin/admin"

	"github.com/gin-gonic/gin"
)

func Setupadminrouter(r *gin.Engine) {
	r.GET("/getmenuitem", Validsession, admin.Getmenuitem)
	r.POST("/userlogin", admin.Userlogin)
	r.GET("/userlogout", Validsession, admin.Userlogout)
	r.GET("/userdetail", Validsession, admin.Userdetail)
	r.POST("/changepassword", Validsession, admin.Changepassword)
	r.GET("/adminlist", Validsession, admin.Adminlist)
	r.POST("/addadmin", Validsession, admin.Addadmin)
	r.POST("/editadmin", Validsession, admin.Editadmin)
	r.POST("/rmadmin", Validsession, admin.Rmadmin)
}
