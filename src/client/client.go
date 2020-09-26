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

const (
	MYSQL  = "mysql"
	CLIENT = "client"
	SELECT = "select"
)

const (
	USERNAME = "username"
	PASSWORD = "password"
	HOSTPORT = "hostport"
	DATABASE = "database"
)

func _sql_meta(m *ice.Message, h string, db string) string {
	if h == "random" {
		h = kit.MDB_RANDOMS
	}
	m.Option(mdb.FIELDS, "time,hash,username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH, h)
	if msg.Append(kit.SSH_PORT) == "" {
		return ""
	}
	p := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", msg.Append(USERNAME), msg.Append(PASSWORD),
		msg.Append(kit.SSH_HOST), msg.Append(kit.SSH_PORT), kit.Select(msg.Append(DATABASE), db))
	return p
}
func _sql_exec(m *ice.Message, p string, s string, arg ...interface{}) *ice.Message {
	m.Log_MODIFY("table", s, "p", p)
	if p == "" {
		return m
	}
	if db, e := sql.Open(MYSQL, p); m.Assert(e) {
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
	m.Log_SELECT("table", s, "p", p)
	if p == "" {
		return m
	}

	if db, e := sql.Open(MYSQL, p); m.Assert(e) {
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
		CLIENT: {Name: "client hash=@key 执行:button 连接 cmd:textarea", Help: "client", Action: map[string]*ice.Action{
			mdb.CREATE: {Name: "create username=root password=root host=localhost port=10000 database=mysql", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option(kit.MDB_HASH), "")
				table := _sql_query(m.Spawn(), p, "explain "+m.Option("cmd")).Append("table")
				_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", table, arg[0], arg[1], m.Option("id")))
			}},
			mdb.REMOVE: {Name: "remove", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option(kit.MDB_HASH), "")
				table := _sql_query(m.Spawn(), p, "explain "+m.Option("cmd")).Append("table")
				_sql_exec(m, p, kit.Format("delete from %s where id=%s", table, m.Option("id")))
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				m.Option(mdb.FIELDS, "time,hash,username,host,port,database")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 || arg[0] == "" {
				m.Option(mdb.FIELDS, "time,hash,username,host,port,database")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
			} else if p := _sql_meta(m, arg[0], ""); strings.Contains(arg[1], "SHOW") || strings.Contains(arg[1], "show") {
				_sql_query(m, p, arg[1])
			} else if strings.Contains(arg[1], "SELECT") || strings.Contains(arg[1], "select") {
				_sql_query(m, p, arg[1])
			} else {
				_sql_exec(m, p, arg[1])
			}
			m.PushAction("删除")
		}},

		SELECT: {Name: "select hash@key database table limit offset auto 连接", Help: "查询", Action: map[string]*ice.Action{
			mdb.CREATE: {Name: "create username=root password=root host=localhost port=10000@key database=mysql", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},

			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				if m.Option(kit.SSH_PORT) != "" {
					m.Cmdy(mdb.MODIFY, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH), arg)
					return
				}

				p := _sql_meta(m, m.Option(kit.MDB_HASH), m.Option(DATABASE))
				_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", m.Option("table"), arg[0], arg[1], m.Option("id")))
			}},
			mdb.REMOVE: {Name: "remove", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				if m.Option(kit.SSH_PORT) != "" {
					m.Cmdy(mdb.DELETE, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH))
					return
				}

				p := _sql_meta(m, m.Option(kit.MDB_HASH), m.Option(DATABASE))
				_sql_exec(m, p, kit.Format("delete from %s where id=%s", m.Option("table"), m.Option("id")))
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				switch arg[0] {
				case kit.SSH_PORT:
					m.Cmdy(server.SERVER)
				default:
					m.Option(mdb.FIELDS, "time,hash,username,host,port,database")
					m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				}
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 || arg[0] == "" {
				m.Option(mdb.FIELDS, "time,hash,username,host,port,database")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
			} else if p := _sql_meta(m, arg[0], ""); len(arg) == 1 || arg[1] == "" {
				_sql_query(m.Spawn(), p, "show databases").Table(func(index int, value map[string]string, head []string) { m.Push(DATABASE, value[head[0]]) })
			} else if p := _sql_meta(m, arg[0], arg[1]); len(arg) == 2 || arg[2] == "" {
				_sql_query(m.Spawn(), p, "show tables").Table(func(index int, value map[string]string, head []string) { m.Push("table", value[head[0]]) })
			} else {
				_sql_query(m, p, fmt.Sprintf("select * from %s limit %s offset %s", arg[2], kit.Select("30", arg, 3), kit.Select("0", arg, 4)))
			}
			m.PushAction("删除")
		}},
	},
}

func init() { server.Index.Merge(Index, nil) }
