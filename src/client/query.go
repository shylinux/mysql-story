package client

import (
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Query struct {
	Client

	prev string `name:"prev" help:"上一页"`
	next string `name:"next" help:"下一页"`
	list string `name:"list session database table id auto page script create" help:"查询"`
}

func (q Query) Modify(m *ice.Message, arg ...string) {
	p := _sql_meta(m, m.Option(kit.MDB_NAME), m.Option(DATABASE))
	_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", m.Option(kit.MDB_TABLE), arg[0], arg[1], m.Option(kit.MDB_ID)))
}
func (q Query) Prev(m *ice.Message, arg ...string) {
	mdb.PrevPageLimit(m.Message, arg[0], arg[1:]...)
}
func (q Query) Next(m *ice.Message, arg ...string) {
	mdb.NextPage(m.Message, arg[0], arg[1:]...)
}
func (q Query) List(m *ice.Message, arg ...string) {
	if len(arg) < 1 || arg[0] == "" { // 连接列表
		q.Client.List(m)
		return
	}

	if dsn := _sql_meta(m, arg[0], ""); len(arg) < 2 || arg[1] == "" { // 数据库列表
		_sql_query(m.Spawn(), dsn, "show databases").Table(func(index int, value map[string]string, head []string) { m.Push(DATABASE, value[head[0]]) })

	} else if dsn := _sql_meta(m, arg[0], arg[1]); len(arg) < 3 || arg[2] == "" { // 关系表列表
		_sql_query(m.Spawn(), dsn, "show tables").Table(func(index int, value map[string]string, head []string) { m.Push(kit.MDB_TABLE, value[head[0]]) })
		m.Table(func(index int, value map[string]string, head []string) {
			msg := _sql_query(m.Spawn(), dsn, kit.Format("show fields from %s", value["table"]))
			m.Push("field", strings.Join(msg.Appendv("Field"), ","))
		})

	} else if len(arg) < 4 || arg[3] == "" { // 数据列表
		_sql_query(m, dsn, kit.Format("select * from %s limit %s offset %s", arg[2], kit.Select("10", arg, 4), kit.Select("0", arg, 5)))

	} else { // 数据详情
		m.Option(mdb.FIELDS, mdb.DETAIL)
		_sql_query(m, dsn, kit.Format("select * from %s where id = %s", arg[2], arg[3]))
	}
	m.StatusTimeCountTotal(_query_total(m, arg...))
}

func init() { ice.CodeModCmd(Query{}) }

func _query_total(m *ice.Message, arg ...string) string {
	if len(arg) > 2 {
		msg := _sql_query(m.Spawn(), _sql_meta(m, arg[0], ""), kit.Format("select count(*) as total from %s", kit.Keys(arg[1], arg[2])))
		return msg.Append(kit.MDB_TOTAL)
	}
	return ""
}
