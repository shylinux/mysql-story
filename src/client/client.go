package client

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/aaa"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/mysql-story/src/server"
	kit "github.com/shylinux/toolkits"
)

func _sql_meta(m *ice.Message, h string, db string) string {
	if h == "random" {
		h = kit.MDB_RANDOMS
	}

	m.Option(mdb.FIELDS, "time,hash,username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH, h)
	m.Assert(msg.Append(tcp.PORT) != "")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db))
}
func _sql_exec(m *ice.Message, p string, s string, arg ...interface{}) *ice.Message {
	if p == "" {
		return m
	}

	m.Log_MODIFY("table", s, "p", p)
	if db, e := sql.Open(MYSQL, p); m.Assert(e) {
		defer db.Close()

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
	if p == "" {
		return m
	}

	m.Log_SELECT("table", s, "p", p)
	if db, e := sql.Open(MYSQL, p); m.Assert(e) {
		defer db.Close()

		m.Debug("what %v %v", s, arg)
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
						m.Push(head[i], fmt.Sprintf("%v", v))
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
	STATMENT = "statment"
	QUERY    = "query"
)

const CLIENT = "client"

var Index = &ice.Context{Name: CLIENT, Help: "客户端",
	Configs: map[string]*ice.Config{
		CLIENT: {Name: CLIENT, Help: "客户端", Value: kit.Data()},
	},
	Commands: map[string]*ice.Command{
		CLIENT: {Name: "client hash 执行:button create cmd:textarea", Help: "客户端", Action: map[string]*ice.Action{
			mdb.CREATE: {Name: "create username=root password=root host=localhost port=10000 database=mysql", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
			mdb.SELECT: {Name: "select hash database statment:textarea", Help: "查询", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option(kit.MDB_HASH), m.Option(DATABASE))
				_sql_query(m, p, m.Option(STATMENT))
			}},
			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option(kit.MDB_HASH), "")
				table := _sql_query(m.Spawn(), p, "explain "+m.Option(cli.CMD)).Append(kit.MDB_TABLE)
				_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", table, arg[0], arg[1], m.Option(kit.MDB_ID)))
			}},
			mdb.REMOVE: {Name: "remove", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option(kit.MDB_HASH), "")
				table := _sql_query(m.Spawn(), p, "explain "+m.Option(cli.CMD)).Append(kit.MDB_TABLE)
				_sql_exec(m, p, kit.Format("delete from %s where id=%s", table, m.Option(kit.MDB_ID)))
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				m.Option(mdb.FIELDS, "time,hash,username,host,port,database")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) < 2 || arg[0] == "" {
				m.Option(mdb.FIELDS, kit.Select("time,hash,username,host,port,database", mdb.DETAIL, len(arg) > 0 && arg[0] != ""))
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				return
			}

			if p := _sql_meta(m, arg[0], ""); strings.Contains(strings.ToLower(arg[1]), "show") {
				_sql_query(m, p, arg[1])
			} else if strings.Contains(strings.ToLower(arg[1]), "select") {
				_sql_query(m, p, arg[1])
			} else {
				_sql_exec(m, p, arg[1])
			}
		}},

		QUERY: {Name: "query hash database table id limit offset auto create", Help: "查询", Action: map[string]*ice.Action{
			mdb.CREATE: {Name: "create username=root password=root host=localhost port=10000@key database=mysql", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				if m.Option(tcp.PORT) != "" {
					m.Cmdy(mdb.MODIFY, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH), arg)
					return
				}

				p := _sql_meta(m, m.Option(kit.MDB_HASH), m.Option(DATABASE))
				_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", m.Option(kit.MDB_TABLE), arg[0], arg[1], m.Option(kit.MDB_ID)))
			}},
			mdb.REMOVE: {Name: "remove", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				if m.Option(tcp.PORT) != "" {
					m.Cmdy(mdb.DELETE, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH))
					return
				}

				p := _sql_meta(m, m.Option(kit.MDB_HASH), m.Option(DATABASE))
				_sql_exec(m, p, kit.Format("delete from %s where id=%s", m.Option(kit.MDB_TABLE), m.Option(kit.MDB_ID)))
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				switch arg[0] {
				case tcp.PORT:
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
				m.PushAction(mdb.REMOVE)
				return
			}

			if p := _sql_meta(m, arg[0], ""); len(arg) == 1 || arg[1] == "" {
				_sql_query(m.Spawn(), p, "show databases").Table(func(index int, value map[string]string, head []string) { m.Push(DATABASE, value[head[0]]) })

			} else if p := _sql_meta(m, arg[0], arg[1]); len(arg) == 2 || arg[2] == "" {
				_sql_query(m.Spawn(), p, "show tables").Table(func(index int, value map[string]string, head []string) { m.Push(kit.MDB_TABLE, value[head[0]]) })

			} else if len(arg) > 3 && arg[3] != "" {
				m.Option(mdb.FIELDS, mdb.DETAIL)
				_sql_query(m, p, fmt.Sprintf("select * from %s where id = %s", arg[2], arg[3]))

			} else {
				_sql_query(m, p, fmt.Sprintf("select * from %s limit %s offset %s", arg[2], kit.Select("30", arg, 4), kit.Select("0", arg, 5)))
			}
		}},
	},
}

func init() { server.Index.Merge(Index) }
