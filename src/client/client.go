package client

import (
	sqls "database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type client struct {
	ice.Hash

	short string `data:"name"`
	field string `data:"time,name,username,host,port,database"`

	create string `name:"create name=biz username=root password=root host=localhost port=10000@key database=mysql" help:"连接"`
	list   string `name:"list name run:button create cmd:textarea" help:"客户端"`
}

func (c client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER)
	}
}
func (c client) List(m *ice.Message, arg ...string) {
	if len(arg) < 2 || arg[0] == "" { // 连接列表
		defer m.PushAction(mdb.REMOVE)
		c.Hash.List(m, kit.Slice(arg, 0, 1)...)
		return
	}

	if dsn := _sql_meta(m, kit.Select(arg[0], kit.MDB_RANDOMS, arg[0] == "random"), ""); strings.Contains(strings.ToLower(arg[1]), "show") {
		_sql_query(m, dsn, arg[1]) // 查询定义
	} else if strings.Contains(strings.ToLower(arg[1]), "select") {
		_sql_query(m, dsn, arg[1]) // 查询数据
	} else {
		_sql_exec(m, dsn, arg[1]) // 操作数据
	}
}

func init() { ice.Cmd("web.code.mysql.client", client{}) }

func _sql_meta(m *ice.Message, h string, db string) string {
	m.Option(mdb.FIELDS, "time,name,username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, m.PrefixKey(), "", mdb.HASH, kit.MDB_NAME, h)
	m.Assert(msg.Append(tcp.PORT) != "")

	return kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append("database"), db))
}
func _sql_open(m *ice.Message, dsn, stm string, cb func(*sqls.DB)) *ice.Message {
	m.Log_MODIFY("dsn", dsn, "stm", stm)
	if db, e := sqls.Open("mysql", dsn); m.Assert(e) {
		defer db.Close()
		cb(db)
	}
	return m
}
func _sql_exec(m *ice.Message, dsn string, stm string, arg ...interface{}) *ice.Message {
	return _sql_open(m, dsn, stm, func(db *sqls.DB) {
		m.Push(kit.MDB_TIME, m.Time())
		if res, err := db.Exec(stm, arg...); err != nil {
			m.Push("", kit.UnMarshal(kit.Format(err)))
		} else {
			if i, e := res.LastInsertId(); e == nil {
				m.Push("lastInsertId", i)
			}
			if i, e := res.RowsAffected(); e == nil {
				m.Push("rowsAffected", i)
			}
		}
	})
}
func _sql_query(m *ice.Message, dsn string, stm string, arg ...interface{}) *ice.Message {
	return _sql_open(m, dsn, stm, func(db *sqls.DB) {
		if rows, err := db.Query(stm, arg...); m.Assert(err) {
			head, err := rows.Columns()
			m.Assert(err)

			var data []interface{}
			for _, _ = range head {
				var item interface{}
				data = append(data, &item)
			}

			for rows.Next() {
				rows.Scan(data...)
				for i, v := range data {
					v = *(v.(*interface{}))
					switch v := v.(type) {
					case []byte:
						m.Push(head[i], string(v))
					default:
						m.Push(head[i], kit.Format("%v", v))
					}
				}
			}
		}
	})
}
