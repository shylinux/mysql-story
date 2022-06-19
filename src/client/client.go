package client

import (
	sqls "database/sql"
	"strings"

	_ "shylinux.com/x/go-sql-mysql"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	MYSQL    = "mysql"
	SESSION  = "session"
	DATABASE = "database"
	TABLE    = "table"
	FIELD    = "field"
)

type Client struct {
	ice.Hash

	short  string `data:"session"`
	field  string `data:"time,session,username,host,port,database"`
	script string `data:"src/sql/"`

	create     string `name:"create session=biz username=root password=root host=localhost port=10000@key database=mysql" help:"连接"`
	list       string `name:"list session database run listScript cmd:textarea" help:"客户端"`
	listScript string `name:"listScript" help:"脚本"`
	catScript  string `name:"catScript" help:"查看"`
	runScript  string `name:"runScript session database file@key" help:"执行"`
}

func (c Client) sql_meta(m *ice.Message, h string, db string) string {
	m.OptionFields("time,session,username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(c), "", mdb.HASH, SESSION, h)
	m.Assert(msg.Append(tcp.PORT) != "")

	return kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db))
}

func (c Client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	case DATABASE:
		m.Cmdy(c, m.Option(SESSION)).Cut(DATABASE)
	case nfs.FILE:
		m.Cmdy(nfs.DIR, arg[1:]).ProcessAgain()
	}
}
func (c Client) Create(m *ice.Message, arg ...string) {
	m.Cmdy(mdb.INSERT, ice.GetTypeKey(c), "", mdb.HASH, arg)
}
func (c Client) List(m *ice.Message, arg ...string) {
	if len(arg) < 1 || arg[0] == "" { // 连接列表
		m.Fields(len(kit.Slice(arg, 0, 1)), m.Config(mdb.FIELD))
		m.Cmdy(mdb.SELECT, ice.GetTypeKey(c), "", mdb.HASH, kit.Slice(arg, 0, 1))
		m.PushAction(c.Hash.Remove)
		m.Action(c.Create)
		m.Sort(SESSION)
	} else if dsn := c.sql_meta(m, kit.Select(arg[0], mdb.RANDOMS, arg[0] == mdb.RANDOM), kit.Select("", arg, 1)); len(arg) < 2 {
		_sql_query(m, dsn, "show databases").ToLowerAppend()
	} else if len(arg) < 3 {
		_sql_query(m, dsn, "show tables")
		m.RenameAppend(m.Appendv(ice.MSG_APPEND)[0], TABLE).Table(func(index int, value map[string]string, head []string) {
			msg := _sql_query(m.Spawn(), dsn, kit.Format("show fields from %s", value[TABLE])).ToLowerAppend()
			m.Push(FIELD, strings.Join(msg.Appendv(FIELD), ice.FS))
		})
	} else if strings.Contains(strings.ToLower(arg[2]), ice.SHOW) {
		_sql_query(m, dsn, arg[2])
	} else if strings.Contains(strings.ToLower(arg[2]), mdb.SELECT) {
		_sql_query(m, dsn, arg[2])
	} else { // 操作数据
		_sql_exec(m, dsn, arg[2])
	}
	m.StatusTimeCount()
}
func (c Client) ListScript(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.DIR, m.Conf(c, kit.Keym(nfs.SCRIPT)), kit.Dict(nfs.DIR_DEEP, ice.TRUE, nfs.DIR_TYPE, nfs.CAT)).RenameAppend(nfs.PATH, nfs.FILE)
	m.PushAction(c.CatScript, c.RunScript)
}
func (c Client) CatScript(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, m.Option(nfs.FILE))
}
func (c Client) RunScript(m *ice.Message, arg ...string) {
	if db, e := sqls.Open(MYSQL, c.sql_meta(m, m.Option(SESSION), m.Option(DATABASE))); m.Assert(e) {
		defer db.Close()

		for _, line := range strings.Split(m.Cmdx(nfs.CAT, kit.Path(m.Option(nfs.FILE))), ";") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			res, err := db.Exec(line)
			m.Push(ice.RES, kit.Format(res))
			m.Push(ice.ERR, kit.Format(err))
			m.Push(nfs.LINE, line)
		}
	}
}

func init() { ice.CodeModCmd(Client{}) }

func _sql_open(m *ice.Message, dsn, stm string, cb func(*sqls.DB)) *ice.Message {
	if db, e := sqls.Open(MYSQL, dsn); m.Assert(e) {
		defer db.Close()
		cb(db)
	}
	return m
}
func _sql_exec(m *ice.Message, dsn string, stm string, arg ...interface{}) *ice.Message {
	return _sql_open(m, dsn, stm, func(db *sqls.DB) {
		m.Log_MODIFY("dsn", dsn, "stm", stm, "arg", arg)
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
		m.Info("dsn: %v stm: %v", dsn, stm)
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
