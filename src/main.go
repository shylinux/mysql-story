package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/mysql-story/src/client"
	_ "shylinux.com/x/mysql-story/src/elasticsearch"
	_ "shylinux.com/x/mysql-story/src/server"
)

func main() { print(ice.Run()) }
