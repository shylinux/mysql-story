package client

import (
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type query struct {
	Client client

	prev string `name:"prev" help:"上一页"`
	next string `name:"next" help:"下一页"`

	list string `name:"query name database table id auto page create" help:"客户端"`
}

func (q query) Modify(m *ice.Message, arg ...string) {
	p := _sql_meta(m, m.Option(kit.MDB_NAME), m.Option("database"))
	_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", m.Option(kit.MDB_TABLE), arg[0], arg[1], m.Option(kit.MDB_ID)))
}
func (q query) Prev(m *ice.Message, arg ...string) {
	mdb.PrevPage(m.Message, _query_total(m, arg...), kit.Slice(arg, 4)...)
}
func (q query) Next(m *ice.Message, arg ...string) {
	mdb.NextPage(m.Message, _query_total(m, arg...), kit.Slice(arg, 4)...)
}
func (q query) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || arg[0] == "" { // 连接列表
		q.Client.List(m)
		return
	}

	if dsn := _sql_meta(m, arg[0], ""); len(arg) == 1 || arg[1] == "" { // 数据库列表
		_sql_query(m.Spawn(), dsn, "show databases").Table(func(index int, value map[string]string, head []string) { m.Push("database", value[head[0]]) })

	} else if dsn := _sql_meta(m, arg[0], arg[1]); len(arg) == 2 || arg[2] == "" { // 关系表列表
		_sql_query(m.Spawn(), dsn, "show tables").Table(func(index int, value map[string]string, head []string) { m.Push(kit.MDB_TABLE, value[head[0]]) })
		m.Table(func(index int, value map[string]string, head []string) {
			msg := _sql_query(m.Spawn(), dsn, kit.Format("show fields from %s", value["table"]))
			m.Push("field", strings.Join(msg.Appendv("Field"), ","))
		})

	} else if len(arg) > 3 && arg[3] != "" { // 数据详情
		m.Option(mdb.FIELDS, mdb.DETAIL)
		_sql_query(m, dsn, kit.Format("select * from %s where id = %s", arg[2], arg[3]))

	} else { // 数据列表
		_sql_query(m, dsn, kit.Format("select * from %s limit %s offset %s", arg[2], kit.Select("30", arg, 4), kit.Select("0", arg, 5)))
	}
	m.StatusTimeCountTotal(_query_total(m, arg...))
}

func init() { ice.Cmd("web.code.mysql.query", query{}) }

func _query_total(m *ice.Message, arg ...string) string {
	if len(arg) > 2 {
		msg := _sql_query(m.Spawn(), _sql_meta(m, arg[0], ""), kit.Format("select count(*) as total from %s", kit.Keys(arg[1], arg[2])))
		return msg.Append("total")
	}
	return ""
}
