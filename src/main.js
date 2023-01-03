Volcanos({river: {
	project: {name: "研发群", storm: {
		studio: {name: "研发 studio", list: [
			{name: "vimer", help: "编辑器", index: "web.code.vimer"},
			{name: "repos", help: "代码库", index: "web.code.git.status"},
			{name: "favor", help: "收藏夹", index: "web.chat.favor"},
			{name: "plan", help: "任务表", index: "web.team.plan"},
			{name: "ctx", help: "上下文", index: "web.wiki.word"},
		]},
		mysql: {name: "存储 mysql", list: [
			{name: "ctx", help: "数据存储", index: "web.wiki.word", args:["usr/mysql-story/src/main.shy"]},
			{name: "ctx", help: "搜索引擎", index: "web.wiki.word", args:["usr/mysql-story/src/elasticsearch/elasticsearch.shy"]},
			{name: "ctx", help: "搜索引擎", index: "web.wiki.word", args:["usr/mysql-story/src/clickhouse/clickhouse.shy"]},
		]},
	}},
	profile: {name: "测试群", storm: {
		release: {name: "发布 release", index: [
			"web.code.webpack",
			"web.code.compile",
			"web.code.publish",
			"web.code.docker.client",
			"web.space",
			"web.dream",
			"web.code.git.server",
			"web.code.git.status",
		]},
		toolkit: {name: "工具 toolkit", index: [
			"web.code.favor",
			"web.code.xterm",
			"web.code.inner",
			"web.code.vimer",
			"web.code.bench",
			"web.code.pprof",
			"web.code.oauth",
		]},
		language: {name: "语言 language", index: [
			"web.code.c",
			"web.code.sh",
			"web.code.py",
			"web.code.shy",
			"web.code.js",
			"web.code.go",
		]},
	}},
	operate: {name: "运维群", storm: {
		aaa: {name: "权限 aaa", index: ["offer", "email", "user", "totp", "sess", "role"]},
		web: {name: "应用 web", index: ["broad", "serve", "space", "dream", "share", "cache", "spide"]},
		cli: {name: "系统 cli", index: ["qrcode", "daemon", "system", "runtime", "mirrors", "forever", "host", "port"]},
		nfs: {name: "文件 nfs", index: ["dir", "cat", "pack", "tail", "trash"]},
	}},
}})

