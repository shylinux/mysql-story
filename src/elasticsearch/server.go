package es

import (
	"path"

	"shylinux.com/x/ice"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	linux string `data:"https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.3.2-linux-x86_64.tar.gz"`
	start string `name:"start port=10004" help:"启动"`
	list  string `name:"list port path auto start install" help:"搜索"`
}

func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "", func(p string, port int) {
		kit.Rewrite(path.Join(p, "config/elasticsearch.yml"), func(text string) string {
			if text == "#http.port: 9200" {
				text = kit.Format("http.port: %d\ntransport.tcp.port: %d\n", port, port+10000)
			}
			return text
		})
	})
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
}

func init() { ice.CodeCtxCmd(server{}) }
