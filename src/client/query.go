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
		defer m.Sort(WHERE)
		s.Hash.Inputs(m, arg...)
	default:
		s.Client.Inputs(m, arg...)
	}
}
func (s Query) Modify(m *ice.Message, arg ...string) {
	if m.Option(TABLE) == "" {
		m.Cmd(s.Client, s.Modify, arg)
		return
	}
	_sql_exec(m.Spawn(), s.meta(m, m.Option(aaa.SESS), m.Option(DATABASE)), kit.Format("update %s set %s='%s' where id=%s",
		m.Option(TABLE), arg[0], arg[1], m.Option(mdb.ID)))
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
	if dsn := s.meta(m, arg[0], kit.Select("", arg, 1)); len(arg) < 4 || arg[3] == "" {
		_sql_query(m, dsn, kit.Format("select * from %s %s limit %s offset %s",
			arg[2], where, kit.Select("10", m.Option(mdb.LIMIT)), kit.Select("0", m.Option(mdb.OFFEND))))
		m.Action(mdb.PAGE, "where:text=`"+kit.Select("", arg, 6)+"`@key")
		m.StatusTimeCountTotal(_query_total(m, s.meta(m, arg[0], ""), where, arg...), TABLE, arg[2])
	} else {
		m.OptionFields(ice.FIELDS_DETAIL)
		_sql_query(m, dsn, kit.Format("select * from %s where id = %s", arg[2], arg[3]))
	}
	return m
}
func (s Query) Prev(m *ice.Message, arg ...string) {
	mdb.NextPageLimit(m.Message, arg[0], arg[1:]...)
}
func (s Query) Next(m *ice.Message, arg ...string) {
	mdb.PrevPage(m.Message, arg[0], arg[1:]...)
}

func init() { ice.CodeModCmd(Query{}) }

func _query_total(m *ice.Message, dsn string, where string, arg ...string) string {
	if len(arg) > 2 {
		msg := _sql_query(m.Spawn(), dsn, kit.Format("select count(*) as total from %s %s", kit.Keys(arg[1], arg[2]), where))
		return msg.Append(mdb.TOTAL)
	}
	return ""
}
