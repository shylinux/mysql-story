package main

import (
	"github.com/shylinux/icebergs"
	_ "github.com/shylinux/icebergs/base"
	_ "github.com/shylinux/icebergs/core"
	_ "github.com/shylinux/icebergs/misc"
    // add local module
    // _ "20200803-mysql_story/src/demo"
)

func main() { println(ice.Run()) }
