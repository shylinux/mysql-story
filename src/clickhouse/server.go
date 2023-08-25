package clickhouse

import (
	"shylinux.com/x/ice"
)

type server struct {
	ice.Code
	linux string `data:"https://packages.clickhouse.com/tgz/stable/clickhouse-common-static-21.1.9.41.tgz"`
	list  string `name:"list port path auto start install" help:"数据库"`
}

func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "usr/bin/clickhouse", "server")
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
}

func init() { ice.CodeCtxCmd(server{}) }
