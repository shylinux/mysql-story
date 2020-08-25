package client

import (
	_ "github.com/go-sql-driver/mysql"
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/mysql-story/src/server"
	kit "github.com/shylinux/toolkits"

	"database/sql"
	"fmt"
	"strings"
)

const CLIENT = "client"

func _sql_meta(m *ice.Message, h string) string {
	m.Option(mdb.FIELDS, "time,hash,username,password,hostport,database")
	msg := m.Cmd(mdb.SELECT, m.Prefix("mysql"), "", mdb.HASH, h)
	p := fmt.Sprintf("%s:%s@%s/%s?charset=utf8", msg.Append("username"), msg.Append("password"), msg.Append("hostport"), msg.Append("database"))
	return p
}
func _sql_exec(m *ice.Message, p string, s string, arg ...interface{}) *ice.Message {
	m.Log_MODIFY("table", s, arg)
	if db, e := sql.Open("mysql", p); m.Assert(e) {
		if res, err := db.Exec(s, arg...); err != nil {
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
func _sql_query(m *ice.Message, p string, s string, arg ...interface{}) *ice.Message {
	m.Log_SELECT("table", s, arg)
	if db, e := sql.Open("mysql", p); m.Assert(e) {
		defer db.Close()
		if rows, err := db.Query(s, arg...); m.Assert(err) {
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
						m.Push(head[i], v)
					}
				}
			}
		}
	}
	return m
}

var Index = &ice.Context{Name: CLIENT, Help: "client",
	Configs: map[string]*ice.Config{
		CLIENT: {Name: CLIENT, Help: "client", Value: kit.Data()},
	},
	Commands: map[string]*ice.Command{
		CLIENT: {Name: "client hash=@key 执行:button 连接 cmd:textarea", Help: "client", Meta: kit.Dict(
			"连接", kit.List(
				kit.MDB_INPUT, "text", "name", "username", "value", "root",
				kit.MDB_INPUT, "text", "name", "password", "value", "root",
				kit.MDB_INPUT, "text", "name", "hostport", "value", "tcp(localhost:10035)",
				kit.MDB_INPUT, "text", "name", "database", "value", "paas",
			),
		), Action: map[string]*ice.Action{
			"connect": {Name: "connect", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				m.Option(mdb.FIELDS, "time,hash,username,hostport,database")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
			}},
			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option("hash"))
				table := _sql_query(m.Spawn(), p, "explain "+m.Option("cmd")).Append("table")
				_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", table, arg[0], arg[1], m.Option("id")))
			}},
			mdb.DELETE: {Name: "delete", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option("hash"))
				table := _sql_query(m.Spawn(), p, "explain "+m.Option("cmd")).Append("table")
				_sql_exec(m, p, kit.Format("delete from %s where id=%s", table, m.Option("id")))
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 || arg[0] == "" {
				m.Option(mdb.FIELDS, "time,hash,username,hostport,database")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				return
			}

			if p := _sql_meta(m, arg[0]); strings.Contains(arg[1], "SELECT") || strings.Contains(arg[1], "select") {
				_sql_query(m, p, arg[1])
				m.PushAction("删除")
			} else {
				_sql_exec(m, p, arg[1])
			}
		}},
	},
}

func init() { server.Index.Merge(Index, nil) }
