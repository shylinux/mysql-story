package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

const (
	WHERE = "where"
)

type Query struct {
	Client
	short string `data:"where"`
	field string `data:"hash,time,where"`
	list  string `name:"list sess@key database@key table@key id auto" help:"查询"`
}

func (s Query) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case TABLE:
		s.List(m, m.Option(aaa.SESS), m.Option(DATABASE)).Cut(arg[0])
	case WHERE:
		s.Hash.Inputs(m, arg...).Sort(arg[0])
	default:
		s.Client.Inputs(m, arg...)
	}
}
func (s Query) Modify(m *ice.Message, arg ...string) {
	if m.Option(TABLE) == "" {
		m.Cmd(s.Client, s.Modify, arg)
		return
	}
	defer m.ProcessRefresh()
	s.open(m, m.Option(aaa.SESS), m.Option(DATABASE), func(db *Driver) {
		db.Exec(m, kit.Format("update %s set %s='%s' where id=%s", m.Option(TABLE), arg[0], arg[1], m.Option(mdb.ID)))
	})
}
func (s Query) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) < 3 || arg[0] == "" || arg[1] == "" || arg[2] == "" {
		m.Cmdy(s.Client, arg)
		return m
	}
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
			m.OptionFields(ice.FIELDS_DETAIL)
			db.Query(m, kit.Format("select * from %s where id = %s", arg[2], arg[3]))
		}
	})
	return m
}
func (s Query) Prev(m *ice.Message, arg ...string) { mdb.NextPageLimit(m.Message, arg[0], arg[1:]...) }
func (s Query) Next(m *ice.Message, arg ...string) { mdb.PrevPage(m.Message, arg[0], arg[1:]...) }

func init() { ice.CodeModCmd(Query{}) }
