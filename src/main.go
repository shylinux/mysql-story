package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/mysql-story/src/client"
	_ "shylinux.com/x/mysql-story/src/server"
	_ "shylinux.com/x/mysql-story/src/studio"
)

func main() { print(ice.Run()) }

func init() {
	ice.Info.NodeIcon = "src/studio/studio.png"
	ice.Info.NodeMain = "web.code.mysql.studio"
}