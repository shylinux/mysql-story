package client

import (
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

const (
	DRIVER   = "driver"
	DATABASE = "database"
	TABLE    = "table"
	WHERE    = "where"
	ID       = "id"
)

type query struct {
	client
	short string `data:"where"`
	field string `data:"time,where"`
	list  string `name:"list sess database table id auto" help:"查询"`
}

func (s query) Drop(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), kit.Select(MYSQL, m.Option(DATABASE)), func(db *Driver) {
		if m.Option(TABLE) != "" {
			db.Query(m, "drop table "+m.Option(TABLE)).ToLowerAppend()
		} else if m.Option(DATABASE) != "" {
			if !kit.IsIn(m.Option(DATABASE), "mysql", "information_schema", "performance_schema") {
				db.Query(m, "drop database "+m.Option(DATABASE)).ToLowerAppend()
			}
		}
	})
}
func (s query) Modify(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), m.Option(DATABASE), func(db *Driver) {
		db.Exec(m, kit.Format("UPDATE %s SET %s = '%s' WHERE id = %s", m.Option(TABLE), arg[0], arg[1], m.Option(ID)))
	})
}
func (s query) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || arg[0] == "" {
		m.Cmdy(s.client, arg).PushAction()
	} else if len(arg) == 1 || arg[1] == "" {
		s.open(m, arg[0], kit.Select("", arg, 1), func(db *Driver) {
			db.Query(m, "show databases").ToLowerAppend()
		})
		m.Table(func(value ice.Maps) {
			if !kit.IsIn(value[DATABASE], "mysql", "information_schema", "performance_schema") {
				m.PushButton(s.Drop)
			} else {
				m.PushButton()
			}
		})
	} else if len(arg) == 2 || arg[2] == "" {
		s.open(m, arg[0], arg[1], func(db *Driver) {
			db.Query(m, "show tables").RenameAppend(kit.Select("", m.Appendv(ice.MSG_APPEND), 0), TABLE).Table(func(value ice.Maps) {
				msg := db.Query(m.Spawn(), kit.Format("select count(*) as total from %s", value[TABLE])).ToLowerAppend()
				m.Push(mdb.TOTAL, msg.Append(mdb.TOTAL))
				msg = db.Query(m.Spawn(), kit.Format("show fields from %s", value[TABLE])).ToLowerAppend()
				m.Push(mdb.FIELD, strings.Join(msg.Appendv(mdb.FIELD), ice.FS))
			})
		})
		if !kit.IsIn(arg[1], "mysql", "information_schema", "performance_schema") {
			m.PushAction(s.Drop)
		}
	} else {
		where := kit.Select("", arg, 6)
		kit.If(where, func() { s.Hash.Create(m.Spawn(), WHERE, where); where = kit.JoinWord(WHERE, where) })
		mdb.OptionPage(m.Message, kit.Slice(arg, 4, 6)...)
		m.OptionDefault(mdb.OFFEND, "0", mdb.LIMIT, "30")
		s.open(m, arg[0], arg[1], func(db *Driver) {
			if len(arg) == 3 || arg[3] == "" {
				db.Query(m, kit.Format("select * from %s %s limit %s offset %s", arg[2], where, m.Option(mdb.LIMIT), m.Option(mdb.OFFEND)))
				total := db.Total(m, where, arg...)
				if where != "" || kit.Int(total) > kit.Int(m.Option(mdb.LIMIT)) {
					m.Action(s.Describe, mdb.PAGE, "where:text=`"+kit.Select("", arg, 6)+"`@key")
				} else {
					m.Action(s.Describe)
				}
				m.StatusTimeCountTotal(db.Total(m, where, arg...), mdb.LIMIT, m.Option(mdb.LIMIT), mdb.OFFEND, m.Option(mdb.OFFEND), TABLE, arg[2])
			} else {
				db.Query(m.FieldsSetDetail(), kit.Format("select * from %s where id = %s", arg[2], arg[3]))
			}
		})
	}
}
func (s query) Prev(m *ice.Message, arg ...string) { mdb.NextPageLimit(m.Message, arg[0], arg[1:]...) }
func (s query) Next(m *ice.Message, arg ...string) { mdb.PrevPage(m.Message, arg[0], arg[1:]...) }
func (s query) Describe(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), m.Option(DATABASE), func(db *Driver) {
		db.Query(m, kit.Format("desc %s", m.Option(TABLE))).ToLowerAppend()
	})
}

func init() { ice.CodeModCmd(query{}) }

type Query struct{ query }

func init() { ice.CodeModCmd(Query{}) }
