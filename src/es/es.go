package es

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/web"
)

type server struct {
	ice.Code
	linux   string `data:"https://elasticsearch.thans.cn/downloads/elasticsearch/elasticsearch-7.3.2-linux-x86_64.tar.gz"`
	darwin  string `data:"https://elasticsearch.thans.cn/downloads/elasticsearch/elasticsearch-7.3.2-linux-x86_64.tar.gz"`
	windows string `data:"https://elasticsearch.thans.cn/downloads/elasticsearch/elasticsearch-7.3.2-linux-x86_64.tar.gz"`
}

func (s server) Get(m *ice.Message, arg ...string) {
	m.Option(web.SPIDE_HEADER, web.ContentType, web.ContentJSON)
}
func (s server) Cmd(m *ice.Message, arg ...string) {
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "bin/elasticsearch")
}

func init() { ice.CodeCtxCmd(server{}) }
