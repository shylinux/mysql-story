package es

import (
	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"

	"path"
	"runtime"
	"strings"
)

const ES = "es"

var Index = &ice.Context{Name: ES, Help: "搜索",
	Configs: ice.Configs{
		ES: {Name: ES, Help: "搜索", Value: kit.Data(
			"address", "http://localhost:9200",
			"windows", "https://elasticsearch.thans.cn/downloads/elasticsearch/elasticsearch-7.3.2-windows-x86_64.zip",
			"darwin", "https://elasticsearch.thans.cn/downloads/elasticsearch/elasticsearch-7.3.2-darwin-x86_64.tar.gz",
			"linux", "https://elasticsearch.thans.cn/downloads/elasticsearch/elasticsearch-7.3.2-linux-x86_64.tar.gz",
		)},
	},
	Commands: ice.Commands{
		ice.CTX_INIT: {Hand: func(m *ice.Message, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, arg ...string) {}},

		ES: {Name: "es port=auto path=auto auto 启动:button 下载", Help: "搜索", Actions: ice.Actions{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(ES, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},

			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				m.Option("install", ".")
				name := path.Base(m.Conf(ES, kit.Keys(kit.MDB_META, runtime.GOOS)))
				name = strings.Join(strings.Split(name, "-")[:2], "-")
				m.Cmdy(code.INSTALL, "start", name, "bin/elasticsearch")
			}},
		}, Hand: func(m *ice.Message, arg ...string) {
			name := path.Base(m.Conf(ES, kit.Keys(kit.MDB_META, runtime.GOOS)))
			name = strings.Join(strings.Split(name, "-")[:2], "-")
			m.Cmdy(code.INSTALL, name, arg)
		}},

		"GET": {Name: "GET 查看:button cmd:text=/", Help: "命令", Hand: func(m *ice.Message, arg ...string) {
			if pod := m.Option("_pod"); pod != "" {
				m.Option("_pod", "")
				m.Cmdy(web.SPACE, pod, m.PrefixKey(), arg)

				if m.Result(0) != ice.ErrWarn || m.Result(1) != ice.ErrNotFound {
					return
				}
				m.Set(ice.MSG_RESULT)
			}

			m.Option(web.SPIDE_HEADER, web.ContentType, web.ContentJSON)
			m.Echo(kit.Formats(kit.UnMarshal(m.Cmdx(web.SPIDE, ice.DEV, web.SPIDE_RAW,
				web.SPIDE_GET, kit.MergeURL2(m.Conf(ES, "meta.address"), kit.Select("/", arg, 0))))))
		}},
		"CMD": {Name: "CMD 执行:button method:select=GET|PUT|POST|DELETE cmd:text=/ data:textarea", Help: "命令", Hand: func(m *ice.Message, arg ...string) {
			if pod := m.Option("_pod"); pod != "" {
				m.Option("_pod", "")
				m.Cmdy(web.SPACE, pod, m.PrefixKey(), arg)

				if m.Result(0) != ice.ErrWarn || m.Result(1) != ice.ErrNotFound {
					return
				}
				m.Set(ice.MSG_RESULT)
			}

			m.Option(web.SPIDE_HEADER, web.ContentType, web.ContentJSON)
			prefix := []string{web.SPIDE, ice.DEV, web.SPIDE_RAW, arg[0], kit.MergeURL2(m.Conf(ES, "meta.address"), arg[1])}

			if len(arg) > 2 {
				prefix = append(prefix, web.SPIDE_DATA, arg[2])
			}
			m.Echo(kit.Formats(kit.UnMarshal(m.Cmdx(prefix))))
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
