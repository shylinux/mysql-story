package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/ice/devops"

	_ "shylinux.com/x/mysql-story/src/clickhouse"
	_ "shylinux.com/x/mysql-story/src/client"
	_ "shylinux.com/x/mysql-story/src/elasticsearch"
	_ "shylinux.com/x/mysql-story/src/mongodb"
	_ "shylinux.com/x/mysql-story/src/postgresql"
	_ "shylinux.com/x/mysql-story/src/server"
	_ "shylinux.com/x/mysql-story/src/sqlite"
	_ "shylinux.com/x/mysql-story/src/studio"
)

func main() { print(ice.Run()) }

func init() { ice.Info.NodeIcon = "src/studio/mysql.png" }
