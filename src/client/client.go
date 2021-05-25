package client

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/aaa"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/mysql-story/src/server"
	kit "github.com/shylinux/toolkits"
)

func _sql_meta(m *ice.Message, h string, db string) string {
	m.Option(mdb.FIELDS, "time,hash,username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH, h)
	m.Assert(msg.Append(tcp.PORT) != "")

	return kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db))
}
func _sql_exec(m *ice.Message, dsn string, stm string, arg ...interface{}) *ice.Message {
	m.Log_MODIFY("dsn", dsn, "stm", stm, "arg", arg)
	if db, e := sql.Open(MYSQL, dsn); m.Assert(e) {
		defer db.Close()

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
	}
	return m
}
func _sql_query(m *ice.Message, dsn string, stm string, arg ...interface{}) *ice.Message {
	m.Log_SELECT("dsn", dsn, "stm", stm, "arg", arg)
	if db, e := sql.Open(MYSQL, dsn); m.Assert(e) {
		defer db.Close()

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
	}
	return m
}

const (
	MYSQL    = "mysql"
	DATABASE = "database"
)

const CLIENT = "client"

func init() {
	server.Index.Merge(&ice.Context{Name: CLIENT, Help: "客户端",
		Configs: map[string]*ice.Config{
			CLIENT: {Name: CLIENT, Help: "客户端", Value: kit.Data()},
		},
		Commands: map[string]*ice.Command{
			ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				m.Watch(server.MYSQL_SERVER_START, m.Prefix(CLIENT))
			}},
			ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			}},

			CLIENT: {Name: "client hash 执行:button create cmd:textarea", Help: "客户端", Action: map[string]*ice.Action{
				server.MYSQL_SERVER_START: {Name: "mysql.server.start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
					m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg, DATABASE, MYSQL)
				}},
				mdb.CREATE: {Name: "create username=root password=root host=localhost port=10000 database=mysql", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
					m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
				}},
				mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
					m.Cmdy(mdb.MODIFY, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH), arg)
				}},
				mdb.REMOVE: {Name: "remove", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
					m.Cmdy(mdb.DELETE, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH))
				}},
				mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
					switch arg[0] {
					case tcp.PORT:
						m.Cmdy(server.SERVER).Appendv(ice.MSG_APPEND, kit.Split("port,time"))
					default:
						m.Option(mdb.FIELDS, "time,hash,username,host,port,database")
						m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
					}
				}},
			}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				if len(arg) < 2 || arg[0] == "" { // 连接列表
					m.Fields(!(len(arg) > 0 && arg[0] != ""), "time,hash,username,host,port,database")
					m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, arg)
					m.PushAction(mdb.REMOVE)
					return
				}

				if dsn := _sql_meta(m, kit.Select(arg[0], kit.MDB_RANDOMS, arg[0] == "random"), ""); strings.Contains(strings.ToLower(arg[1]), "show") {
					_sql_query(m, dsn, arg[1]) // 查询定义
				} else if strings.Contains(strings.ToLower(arg[1]), "select") {
					_sql_query(m, dsn, arg[1]) // 查询数据
				} else {
					_sql_exec(m, dsn, arg[1]) // 操作数据
				}
			}},
		},
	})
}
