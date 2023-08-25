package es

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	linux   string `data:"https://mirrors.huaweicloud.com/elasticsearch/7.6.2/elasticsearch-7.6.2-linux-x86_64.tar.gz"`
	darwin  string `data:"https://mirrors.huaweicloud.com/kibana/7.6.2/kibana-7.6.2-darwin-x86_64.tar.gz"`
	windows string `data:"https://mirrors.huaweicloud.com/elasticsearch/7.6.2/elasticsearch-7.6.2-windows-x86_64.zip"`
	start   string `name:"start port=10004" help:"启动"`
	list    string `name:"list port path auto start install" help:"搜索"`
}

func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "bin/elasticsearch", func(p string, port int) {
		nfs.Rewrite(m.Message, path.Join(p, "config/elasticsearch.yml"), func(text string) string {
			if text == "#http.port: 9200" {
				text = kit.Format("http.port: %d\ntransport.tcp.port: %d\n", port, port+10000)
			}
			return text
		})
	})
}
func (s server) List(m *ice.Message, arg ...string) { s.Code.List(m, "", arg...) }

func init() { ice.CodeCtxCmd(server{}) }
