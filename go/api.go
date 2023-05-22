package main

import (
	"layuiadmin/dao"
	"layuiadmin/router"
)

func main() {
	dao.Initmysql()
	router.SetupRouter().Run(":8080")
}
