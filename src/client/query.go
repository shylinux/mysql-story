package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Query struct {
	Client

	list string `name:"list session database table id auto" help:"查询"`
}

func (q Query) Modify(m *ice.Message, arg ...string) {
	if m.Option(TABLE) == "" {
		m.Cmd(q.Client, q.Modify, arg)
		return
	}
	_sql_exec(m, q.sql_meta(m, m.Option(SESSION), m.Option(DATABASE)),
		kit.Format("update %s set %s='%s' where id=%s", m.Option(TABLE), arg[0], arg[1], m.Option(kit.MDB_ID)))
	m.SetAppend()
}
func (q Query) Prev(m *ice.Message, arg ...string) {
	mdb.PrevPageLimit(m.Message, arg[0], arg[1:]...)
}
func (q Query) Next(m *ice.Message, arg ...string) {
	mdb.NextPage(m.Message, arg[0], arg[1:]...)
}
func (q Query) List(m *ice.Message, arg ...string) {
	if len(arg) < 1 || arg[0] == "" || len(arg) < 2 || arg[1] == "" || len(arg) < 3 || arg[2] == "" {
		m.Cmdy(q.Client, arg)
		return
	}

	where := kit.Select("", arg, 6)
	if where != "" {
		where = "where " + where
	}
	if dsn := q.sql_meta(m, arg[0], kit.Select("", arg, 1)); len(arg) < 4 || arg[3] == "" { // 数据列表
		_sql_query(m, dsn, kit.Format("select * from %s %s limit %s offset %s",
			arg[2], where, kit.Select("10", arg, 4), kit.Select("0", arg, 5)))
		m.Option("limit", kit.Select("", arg, 4))
		m.Option("offend", kit.Select("", arg, 5))
		m.Action("page", "where:text=`"+kit.Select("", arg, 6)+"`")

	} else { // 数据详情
		m.OptionFields(mdb.DETAIL)
		_sql_query(m, dsn, kit.Format("select * from %s where id = %s", arg[2], arg[3]))
	}
	m.StatusTimeCountTotal(_query_total(m, q.sql_meta(m, arg[0], ""), where, arg...))
}

func init() { ice.CodeModCmd(Query{}) }

func _query_total(m *ice.Message, dsn string, where string, arg ...string) string {
	if len(arg) > 2 {
		msg := _sql_query(m.Spawn(), dsn, kit.Format("select count(*) as total from %s %s", kit.Keys(arg[1], arg[2]), where))
		return msg.Append(kit.MDB_TOTAL)
	}
	return ""
}
