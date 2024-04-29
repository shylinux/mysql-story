package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/mysql-story/src/client"
	_ "shylinux.com/x/mysql-story/src/server"
	_ "shylinux.com/x/mysql-story/src/sqlite"
	_ "shylinux.com/x/mysql-story/src/studio"
)

func main() { print(ice.Run()) }

func init() { ice.Info.NodeIcon = "src/studio/mysql.png" }
