package client

import (
	_sql "database/sql"
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
	ice.Code
	ice.Hash
	short string `data:"session"`
	field string `data:"time,session,username,host,port,database"`

	connect string `name:"connect session=biz username=root password=root host=127.0.0.1 port=10002@key database=mysql" help:"连接"`
	list    string `name:"list session@key database@key run cmd:textarea" help:"会话"`
}

func (s Client) meta(m *ice.Message, h string, db string) string {
	m.OptionFields("username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(s), "", mdb.HASH, SESSION, h)
	m.Assert(h != "" && msg.Append(tcp.PORT) != "")
	return kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db))
}

func (s Client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case SESSION:
		s.List(m).Cut(SESSION)
	case DATABASE:
		s.List(m, m.Option(SESSION)).Cut(DATABASE)
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	}
}
func (s Client) Connect(m *ice.Message, arg ...string) {
	m.Cmd(mdb.INSERT, ice.GetTypeKey(s), "", mdb.HASH, arg)
}
func (s Client) Xterm(m *ice.Message, arg ...string) {
	m.OptionFields("username,password,host,port,database")
	msg := m.Cmd(mdb.SELECT, ice.GetTypeKey(s), "", mdb.HASH, m.OptionSimple(SESSION))
	s.Code.Xterm(m, kit.Format("%s -h%s -P%s -u%s -p%s", kit.Path(ice.USR_LOCAL_DAEMON, msg.Append(tcp.PORT), ice.BIN, MYSQL),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD)), arg...)
}
func (s Client) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) < 1 || arg[0] == "" { // 会话列表
		s.Hash.List(m, arg...).Sort(SESSION).Action(s.Connect).PushAction(s.Xterm, s.Remove)

	} else if dsn := s.meta(m, arg[0], kit.Select("", arg, 1)); len(arg) < 2 { // 数据库列表
		_sql_query(m, dsn, "show databases").ToLowerAppend()

	} else if len(arg) < 3 || arg[2] == "" { // 关系表列表
		_sql_query(m, dsn, "show tables").RenameAppend(kit.Select("", m.Appendv(ice.MSG_APPEND), 0), TABLE).Tables(func(value ice.Maps) {
			msg := _sql_query(m.Spawn(), dsn, kit.Format("show fields from %s", value[TABLE])).ToLowerAppend()
			m.Push(FIELD, strings.Join(msg.Appendv(FIELD), ice.FS))
		}).Action(s.ListScript)

	} else if cmd := strings.ToLower(strings.TrimSpace(arg[2])); strings.HasPrefix(cmd, ice.SHOW) { // 查询定义
		_sql_query(m, dsn, arg[2])
	} else if strings.HasPrefix(cmd, mdb.SELECT) { // 查询数据
		_sql_query(m, dsn, arg[2])
	} else { // 操作数据
		_sql_exec(m, dsn, arg[2])
	}
	return m
}
func (s Client) ListScript(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.DIR, ice.SRC, kit.Dict(nfs.DIR_DEEP, ice.TRUE, nfs.DIR_REG, ".*.sql")).RenameAppend(nfs.PATH, nfs.FILE)
	m.PushAction(s.CatScript, s.RunScript)
}
func (s Client) CatScript(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, m.Option(nfs.FILE))
}
func (s Client) RunScript(m *ice.Message, arg ...string) {
	if db, e := _sql.Open(MYSQL, s.meta(m, m.Option(SESSION), m.Option(DATABASE))); m.Assert(e) {
		defer db.Close()

		for _, line := range strings.Split(m.Cmdx(nfs.CAT, kit.Path(m.Option(nfs.FILE))), ";") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			res, err := db.Exec(line)
			m.Push(ice.ERR, kit.Format(err))
			m.Push(ice.RES, kit.Format(res))
			m.Push(nfs.LINE, line)
		}
	}
}

func init() { ice.CodeModCmd(Client{}) }

func _sql_open(m *ice.Message, dsn, stm string, cb func(*_sql.DB)) *ice.Message {
	if db, e := _sql.Open(MYSQL, dsn); m.Assert(e) {
		defer db.Close()
		cb(db)
	}
	return m
}
func _sql_exec(m *ice.Message, dsn string, stm string, arg ...ice.Any) *ice.Message {
	return _sql_open(m, dsn, stm, func(db *_sql.DB) {
		m.Logs(mdb.MODIFY, "dsn", dsn, "stm", stm, "arg", arg)
		m.Push(mdb.TIME, m.Time())
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
func _sql_query(m *ice.Message, dsn string, stm string, arg ...ice.Any) *ice.Message {
	return _sql_open(m, dsn, stm, func(db *_sql.DB) {
		m.Logs(mdb.SELECT, "dsn", dsn, "stm", stm, "arg", arg)
		if rows, err := db.Query(stm, arg...); m.Assert(err) {
			head, err := rows.Columns()
			m.Assert(err)

			var data ice.List
			for _, _ = range head {
				var item ice.Any
				data = append(data, &item)
			}

			defer m.StatusTimeCount()
			for rows.Next() {
				rows.Scan(data...)
				for i, v := range data {
					v = *(v.(*ice.Any))
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
