package leveldb

import (
	"shylinux.com/x/ice"
)

type project struct {
	ice.Code
	source string `data:"http://mirrors.aliyun.com/gnu/leveldb/leveldb-4.2.53.tar.gz"`
}

func (s project) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, "", "", func(string) {})
}
func (s project) List(m *ice.Message, arg ...string) {
	s.Code.Source(m, "", arg...)
}

func init() { ice.Cmd("web.code.leveldb.project", project{}) }
