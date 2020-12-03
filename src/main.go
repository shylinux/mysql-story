package main

import (
	ice "github.com/shylinux/icebergs"
	_ "github.com/shylinux/icebergs/base"
	_ "github.com/shylinux/icebergs/core"
	_ "github.com/shylinux/icebergs/misc"

	_ "github.com/shylinux/mysql-story/src/client"
	_ "github.com/shylinux/mysql-story/src/server"
)

func main() { print(ice.Run()) }
