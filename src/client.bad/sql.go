package client

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/lex"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	SQL = "sql"

	SHOW   = "SHOW"
	CREATE = "CREATE"
	ALTER  = "ALTER"
	DROP   = "DROP"

	INSERT = "INSERT"
	DELETE = "DELETE"
	SELECT = "SELECT"
	UPDATE = "UPDATE"
)

type sql struct {
	ice.Code
	ice.Lang
	Client
	sess     string `data:"biz"`
	database string `data:"demo"`

	list string `name:"list sess@key database@key path auto" help:"脚本"`
}

func (s sql) Init(m *ice.Message, arg ...string) {
	s.Lang.Init(m, nfs.SCRIPT, m.Resource(""))
}
func (s sql) Render(m *ice.Message, arg ...string) {
	ctx.ProcessField(m.Message, m.PrefixKey(), []string{m.Config(aaa.SESS), m.Config(DATABASE), path.Join(m.Option(nfs.PATH), m.Option(nfs.FILE))}, arg...)
}
func (s sql) Engine(m *ice.Message, arg ...string) {
	ctx.OptionFromConfig(m.Message, aaa.SESS, DATABASE)
	msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(s.Client), "", mdb.HASH, m.OptionSimple(aaa.SESS), kit.Dict(ice.MSG_FIELDS, "username,password,host,port"))
	s.Code.Xterm(m, "", []string{mdb.TYPE, kit.Format("%s -h%s -P%s -u%s -p%s %s", MYSQL,
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD), m.Option(DATABASE)),
		mdb.NAME, m.Option(DATABASE),
	}, arg...)
}
func (s sql) List(m *ice.Message, arg ...string) {
	if len(arg) < 2 || arg[0] == "" || arg[1] == "" {
		m.Cmdy(s.Client, kit.Slice(arg, 0, 2))
	} else if len(arg) < 3 {
		m.Cmdy(nfs.DIR, nfs.SRC, kit.Dict(nfs.DIR_REG, kit.ExtReg(SQL), nfs.DIR_DEEP, ice.TRUE))
	} else {
		s.open(m, arg[0], arg[1], func(db *Driver) {
			ls := []string{}
			kit.For(strings.Split(m.Cmdx(nfs.CAT, arg[2]), lex.NL), func(text string) {
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
}

func init() { ice.CodeModCmd(sql{}) }
