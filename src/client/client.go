package client

import (
	"strings"

	_ "shylinux.com/x/go-sql-mysql"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	MYSQL    = "mysql"
	DATABASE = "database"
	TABLE    = "table"
	DRIVER   = "driver"
)

type Client struct {
	ice.Code
	ice.Hash
	short string `data:"sess"`
	field string `data:"time,sess,username,host,port,database"`

	create string `name:"create sess*=biz username*=root password*=root host*=127.0.0.1 port*=10001 database*=mysql" help:"连接"`
	list   string `name:"list sess@key database@key auto create stmt:textarea" help:"数据库"`
}

func (s Client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case aaa.SESS:
		m.Cmdy(s).Cut(arg[0])
	case aaa.USERNAME:
		m.Cmdy(aaa.USER).Cut("username,usernick")
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
		m.Push(arg[0], "3306")
	case DATABASE:
		m.Cmdy(s, m.Option(aaa.SESS)).Cut(arg[0])
	}
}
func (s Client) Create(m *ice.Message, arg ...string) {
	m.Cmd(mdb.INSERT, ice.GetTypeKey(s), "", mdb.HASH, arg)
}
func (s Client) Remove(m *ice.Message, arg ...string) {
	m.Cmd(mdb.DELETE, ice.GetTypeKey(s), "", mdb.HASH, m.OptionSimple(aaa.SESS))
}
func (s Client) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) < 1 || arg[0] == "" {
		s.Hash.List(m, arg...).Sort(aaa.SESS).PushAction(s.Xterm, s.Remove)
		return m
	}
	s.open(m, arg[0], kit.Select("", arg, 1), func(db *Driver) {
		if len(arg) < 2 {
			db.Query(m, "show databases").ToLowerAppend()
		} else if len(arg) < 3 || arg[2] == "" {
			db.Query(m, "show tables").RenameAppend(kit.Select("", m.Appendv(ice.MSG_APPEND), 0), TABLE).Table(func(value ice.Maps) {
				msg := db.Query(m.Spawn(), kit.Format("show fields from %s", value[TABLE])).ToLowerAppend()
				m.Push(mdb.FIELD, strings.Join(msg.Appendv(mdb.FIELD), ice.FS))
			}).Action(s.ListScript)
		} else if cmd := strings.ToLower(strings.TrimSpace(arg[2])); strings.HasPrefix(cmd, mdb.SELECT) || strings.HasPrefix(cmd, "show") {
			db.Query(m, arg[2])
		} else {
			db.Exec(m, arg[2])
		}
	})

	return m
}
func (s Client) Xterm(m *ice.Message, arg ...string) {
	m.OptionFields("username,password,host,port")
	msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(s), "", mdb.HASH, m.OptionSimple(aaa.SESS))
	p := kit.Path(ice.USR_LOCAL_DAEMON, msg.Append(tcp.PORT), "bin/mysql")
	if m.Warn(!nfs.Exists(m.Message, p), ice.ErrNotFound, ice.BIN, p) {
		return
	}
	s.Code.Xterm(m, []string{mdb.TYPE, kit.Format("%s -h%s -P%s -u%s -p%s", p, msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD))}, arg...)
}
func (s Client) ListScript(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.DIR, ice.SRC, kit.Dict(nfs.DIR_DEEP, ice.TRUE, nfs.DIR_REG, kit.ExtReg(SQL))).RenameAppend(nfs.PATH, nfs.FILE).PushAction(s.CatScript, s.RunScript)
}
func (s Client) CatScript(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, m.Option(nfs.FILE))
}
func (s Client) RunScript(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), m.Option(DATABASE), func(db *Driver) {
		for _, line := range strings.Split(m.Cmdx(nfs.CAT, kit.Path(m.Option(nfs.FILE))), ";") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			res, err := db.DB.Exec(line)
			m.Push(ice.ERR, kit.Format(err)).Push(ice.RES, kit.Format(res)).Push(nfs.LINE, line)
		}
	})
}

func init() { ice.CodeModCmd(Client{}) }

func (s Client) open(m *ice.Message, h string, db string, cb func(*Driver)) {
	m.OptionFields("username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(s), "", mdb.HASH, aaa.SESS, h)
	m.Assert(h != "" && msg.Append(tcp.PORT) != "")
	Open(m, msg.Append(DRIVER), kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD), msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db)), cb)
}
