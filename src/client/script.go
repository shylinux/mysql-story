package client

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/lex"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type script struct {
	client client
	list   string `name:"list sess path file auto" help:"脚本"`
}

func (s script) Engine(m *ice.Message, arg ...string) {
	s.client.open(m, m.Option(aaa.SESS), m.Option("database"), func(db *Driver) {
		ls := []string{}
		kit.For(strings.Split(m.Cmdx(nfs.CAT, path.Join(arg[2], arg[1])), lex.NL), func(text string) {
			ls = append(ls, strings.Split(text, "--")[0])
		})
		kit.For(strings.Split(strings.Join(ls, lex.NL), ";"), func(stm string) {
			if stm == "" {
				return
			} else if kit.HasPrefix(stm, SHOW, SELECT) {
				msg := db.Query(m.Spawn(), stm)
				m.Push("stm", stm).Push("err", msg.TableEcho())
				m.Push("lastInsertId", "")
				m.Push("rowsAffected", "")
			} else {
				db.Exec(m, stm)
			}
		})
	})
}
func (s script) List(m *ice.Message, arg ...string) {
	if kit.HasPrefixList(arg, ctx.ACTION) {
		m.Cmdy(code.VIMER, arg)
	} else if len(arg) == 0 {
		m.Cmdy(s.client)
	} else if len(arg) == 1 || arg[1] == nfs.USR {
		m.Cmdy(nfs.DIR, arg[1:])
	} else if len(arg) == 2 {
		m.Option(nfs.DIR_REG, kit.ExtReg("sql"))
		m.DirDeepAll(arg[1], ".", nil).RenameAppend(nfs.PATH, nfs.FILE)
	} else {
		m.Cmdy(nfs.CAT, path.Join(arg[1], arg[2])).Display("")
	}
}

func init() { ice.CodeModCmd(script{}) }

type Script struct{ script }

func init() { ice.CodeModCmd(Script{}) }
