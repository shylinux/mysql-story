package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
)

const (
	SQL = "sql"

	SHOW   = "SHOW"
	CREATE = "CREATE"
	ALTER  = "ALTER"
	DROP   = "DROP"

	INSERT = "INSERT"
	DELETE = "DELETE"
	SELECT = "SELECT"
	UPDATE = "UPDATE"
)

type sql struct {
	ice.Lang
}

func (s sql) Init(m *ice.Message, arg ...string) {
	s.Lang.Init(m, nfs.SCRIPT, m.Resource(""))
}

func init() { ice.CodeModCmd(sql{}) }
