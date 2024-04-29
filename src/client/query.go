package client

import (
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	DRIVER   = "driver"
	DATABASE = "database"
	TABLE    = "table"
	WHERE    = "where"
)

type query struct {
	ice.Hash
	client Client
	short  string `data:"where"`
	field  string `data:"time,where"`
	list   string `name:"list sess database table id auto"`
}

func (s query) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		m.Cmdy(s.client, arg)
	} else if len(arg) == 1 {
		s.open(m, arg[0], kit.Select("", arg, 1), func(db *Driver) {
			db.Query(m, "show databases").ToLowerAppend()
		})
	} else if len(arg) == 2 {
		s.open(m, arg[0], arg[1], func(db *Driver) {
			db.Query(m, "show tables").RenameAppend(kit.Select("", m.Appendv(ice.MSG_APPEND), 0), TABLE).Table(func(value ice.Maps) {
				msg := db.Query(m.Spawn(), kit.Format("select count(*) as total from %s", value[TABLE])).ToLowerAppend()
				m.Push(mdb.TOTAL, msg.Append(mdb.TOTAL))
				msg = db.Query(m.Spawn(), kit.Format("show fields from %s", value[TABLE])).ToLowerAppend()
				m.Push(mdb.FIELD, strings.Join(msg.Appendv(mdb.FIELD), ice.FS))
			})
		})
	} else {
		where := kit.Select("", arg, 6)
		if where != "" {
			s.Hash.Create(m.Spawn(), WHERE, where)
			where = WHERE + ice.SP + where
		}
		mdb.OptionPage(m.Message, kit.Slice(arg, 4, 6)...)
		s.open(m, arg[0], kit.Select("", arg, 1), func(db *Driver) {
			if len(arg) < 4 || arg[3] == "" {
				db.Query(m, kit.Format("select * from %s %s limit %s offset %s", arg[2], where, kit.Select("10", m.Option(mdb.LIMIT)), kit.Select("0", m.Option(mdb.OFFEND))))
				m.Action(mdb.PAGE, "where:text=`"+kit.Select("", arg, 6)+"`@key")
				m.StatusTimeCountTotal(db.Total(m, where, arg...), TABLE, arg[2])
			} else {
				m.FieldsSetDetail()
				db.Query(m, kit.Format("select * from %s where id = %s", arg[2], arg[3]))
			}
		})
	}
}
func (s query) Prev(m *ice.Message, arg ...string) { mdb.NextPageLimit(m.Message, arg[0], arg[1:]...) }
func (s query) Next(m *ice.Message, arg ...string) { mdb.PrevPage(m.Message, arg[0], arg[1:]...) }

func init() { ice.CodeModCmd(query{}) }

func (s query) open(m *ice.Message, sess string, db string, cb func(*Driver)) {
	msg := m.Cmd(s.client, sess)
	Open(m, msg.Append(DRIVER), kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD), msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db)), cb)
}

type Query struct{ query }

func init() { ice.CodeModCmd(Query{}) }
