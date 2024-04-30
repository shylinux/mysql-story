package client

import (
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	DRIVER   = "driver"
	MYSQL    = "mysql"
	DATABASE = "database"
	TABLE    = "table"
)

type client struct {
	ice.Hash
	short  string `data:"sess"`
	field  string `data:"time,sess,username,password,host,port,database,driver"`
	create string `name:"create sess*=biz username*=root password*=root host*=127.0.0.1 port*=10001 database*=mysql driver*=mysql"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case aaa.SESS:
		m.Cmdy(s).Cut(arg[0])
	case aaa.USERNAME:
		m.Cmdy(aaa.USER).Cut("username,usernick")
	case tcp.PORT:
		m.Cmdy(tcp.PORT, mdb.INPUTS, arg).Push(arg[0], "3306")
	case DRIVER:
		m.Push(arg[0], MYSQL)
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s client) Create(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), "", func(db *Driver) {
		db.Exec(m, kit.Format("create database %s charset=utf8mb4", m.Option(DATABASE)))
	}).ProcessRefresh()
}
func (s client) Drop(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), "", func(db *Driver) {
		db.Exec(m, kit.Format("drop database %s", m.Option(DATABASE)))
	}).ProcessRefresh()
}
func (s client) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) < 1 || arg[0] == "" {
		s.Hash.List(m, arg...).Sort(aaa.SESS).PushAction(s.Create, s.Disconn)
		return m
	}
	s.open(m, arg[0], kit.Select("", arg, 1), func(db *Driver) {
		if len(arg) < 2 {
			db.Query(m, "show databases").ToLowerAppend().PushAction(s.Script, s.Xterm, s.Grant, s.Drop)
		} else if len(arg) < 3 || arg[2] == "" {
			db.Query(m, "show tables").RenameAppend(kit.Select("", m.Appendv(ice.MSG_APPEND), 0), TABLE).Table(func(value ice.Maps) {
				msg := db.Query(m.Spawn(), kit.Format("select count(*) as total from %s", value[TABLE])).ToLowerAppend()
				m.Push(mdb.TOTAL, msg.Append(mdb.TOTAL))
				msg = db.Query(m.Spawn(), kit.Format("show fields from %s", value[TABLE])).ToLowerAppend()
				m.Push(mdb.FIELD, strings.Join(msg.Appendv(mdb.FIELD), ice.FS))
			})
		} else if kit.HasPrefix(strings.ToLower(strings.TrimSpace(arg[2])), mdb.SELECT, "show") {
			db.Query(m, arg[2])
		} else {
			db.Exec(m, arg[2])
		}
	})
	return m
}
func (s client) Script(m *ice.Message, arg ...string) {
	ctx.ProcessField(m.Message, ice.GetTypeKey(sql{}), []string{m.Option(aaa.SESS), m.Option(DATABASE)}, arg...)
}
func (s client) Xterm(m *ice.Message, arg ...string) {
	m.ProcessXterm("", func() []string {
		msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(s), "", mdb.HASH, m.OptionSimple(aaa.SESS), kit.Dict(ice.MSG_FIELDS, "username,password,host,port"))
		return []string{mdb.TYPE, kit.Format("%s -h%s -P%s -u%s -p%s %s",
			MYSQL, msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD), m.Option(DATABASE)),
			mdb.NAME, m.Option(DATABASE),
		}
	}, arg...)
}
func (s client) Grant(m *ice.Message, arg ...string) {
	ctx.ProcessField(m.Message, ice.GetTypeKey(Grant{}), []string{m.Option(aaa.SESS)}, arg...)
}

func init() { ice.CodeModCmd(client{}) }

func (s client) open(m *ice.Message, h string, db string, cb func(*Driver)) *ice.Message {
	msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(s), "", mdb.HASH, aaa.SESS, h, kit.Dict(ice.MSG_FIELDS, "username,password,host,port,database,driver"))
	Open(m, msg.Append(DRIVER), kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD), msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db)), cb)
	return m
}

type Client struct{ client }
