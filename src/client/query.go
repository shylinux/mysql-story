package client

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/mysql-story/src/server"
	kit "github.com/shylinux/toolkits"
)

const QUERY = "query"

func init() {
	server.Index.Merge(&ice.Context{Commands: map[string]*ice.Command{
		QUERY: {Name: "query hash database table id limit offset auto create", Help: "查询", Action: map[string]*ice.Action{
			mdb.CREATE: {Name: "create username=root password=root host=localhost port=10000@key database=mysql", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(CLIENT, mdb.CREATE, arg)
			}},
			mdb.REMOVE: {Name: "remove", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(CLIENT, mdb.REMOVE, arg)
			}},
			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				p := _sql_meta(m, m.Option(kit.MDB_HASH), m.Option(DATABASE))
				_sql_exec(m, p, kit.Format("update %s set %s='%s' where id=%s", m.Option(kit.MDB_TABLE), arg[0], arg[1], m.Option(kit.MDB_ID)))
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(CLIENT, mdb.INPUTS, arg)
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 || arg[0] == "" { // 连接列表
				m.Fields(true, "time,hash,username,host,port,database")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				m.PushAction(mdb.REMOVE)
				return
			}

			if dsn := _sql_meta(m, arg[0], ""); len(arg) == 1 || arg[1] == "" { // 数据库列表
				_sql_query(m.Spawn(), dsn, "show databases").Table(func(index int, value map[string]string, head []string) { m.Push(DATABASE, value[head[0]]) })

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
			m.Status(kit.MDB_TIME, m.Time(), kit.MDB_COUNT, m.FormatSize())
		}},
	}})
}
