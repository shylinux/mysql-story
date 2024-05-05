Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) { can.ui = can.onappend.layout(can), can.onimport._project(can, msg)
		can.sup.onimport._field = function(sup, msg, cb) { msg.Table(function(value) {
			can.onimport._content(can, [msg.Option(aaa.SESS), msg.Option("database"), msg.Option("table"), value.index], value)
		}); return true }
	},
	_project: function(can, msg) { var _select
		msg.Table(function(value) { value.icon = icon.sess, value.nick = `${value.sess}(${value.host}:${value.port})`
			var target = can.onimport.item(can, value, function(event, value, show) {
				show == undefined && can.run(event, [value.sess], function(msg) {
					can.onimport._database(can, msg, value.sess, target)
				})
			}); _select = _select||target, value.sess == can.db.hash[0] && (_select = target)
		}), _select && _select.click()
	},
	_database: function(can, msg, sess, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) { value.icon = icon.database, value.nick = value.database
			value._select = sess == can.db.hash[0] && value.database == can.db.hash[1]
			return value
		}), function(event, value, show, target) {
			show == undefined && can.run(event, [sess, value.database], function(msg) {
				can.onimport._table(can, msg, sess, value.database, target)
			})
		}, function() {}, target)
		can.onimport.itemlist(can, [{icon: "bi bi-person-lock", nick: "_grant", _select: sess == can.db.hash[0] && "_grant" == can.db.hash[1]}], function(event, value, show, target) {
			can.onimport._content(can, [sess, "_grant"], {index: "web.code.mysql.grant", args: [sess]}, target, value)
		}, function() {}, target)
		can.onimport.itemlist(can, [{icon: icon.xterm, nick: "_shell", _select: sess == can.db.hash[0] && "_shell" == can.db.hash[1]}], function(event, value, show, target) {
			can.onimport._content(can, [sess, "_shell"], {index: "web.code.mysql.shell", args: [sess], style: html.OUTPUT}, target, value)
		}, function() {}, target)
		can.onimport.itemlist(can, [{icon: icon.path, nick: "_script", _select: sess == can.db.hash[0] && "_script" == can.db.hash[1]}], function(event, value, show, target) {
			show == undefined && can.run(event, [nfs.SCRIPT, sess, nfs.SRC], function(msg) {
				can.onimport._script(can, msg, sess, target)
			})
		}, function() {}, target)
	},
	_script: function(can, msg, sess, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) { value.icon = icon.file, value.nick = value.file
			value._select = sess == can.db.hash[0] && "_script" == can.db.hash[1] && value.file == can.db.hash[2]
			return value
		}), function(event, value, show, target) {
			can.onimport._content(can, [sess, "_script", value.file], {index: "web.code.mysql.script", args: [sess, nfs.SRC, value.file]}, target, value)
		}, function() {

		}, target)
	},
	_table: function(can, msg, sess, database, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) { value.icon = icon.table, value.nick = `${value.table}(${value.total})`
			value._select = sess == can.db.hash[0] && database == can.db.hash[1] && value.table == can.db.hash[2]
			return value
		}), function(event, value, show, target) {
			can.onimport._content(can, [sess, database, value.table, "query"], {index: "web.code.mysql.query", args: [sess, database, value.table]}, target, value)
		}, function() {

		}, target)
	},
	_content: function(can, keys, meta, target, value) { if (target._tabs) { return target._tabs.click() } var key = keys.join(".")
		return target._tabs = can.onimport.tabs(can, [{icon: value.icon, nick: can.core.Keys(keys.slice(1, 3)), title: key}], function() { can.onexport.hash(can, keys)

			target && can.page.Select(can, can.ui.project, html.DIV_ITEM, function(target) { can.page.ClassList.del(can, target, html.SELECT) })
			for (var p = target; p; p = p.parentNode.previousElementSibling) { can.page.ClassList.add(can, p, html.SELECT) }

			if (can.onmotion.cache(can, function(save, load) { save({_content_plugin: can.ui._content_plugin})
				load(key, function(bak) { can.ui._content_plugin = bak._content_plugin }); return key
			}, can.ui.content)) { return can.onimport.layout(can) }
			can.onappend.plugin(can, meta, function(sub) { can.ui._content_plugin = sub, can.onimport.layout(can)
				sub.onexport.output = function(_sub, msg) {
					can.page.Select(sub, sub._option, "div.item.text.id", function(target) {
						can.onmotion.toggle(can, target, msg.append && msg.append.indexOf(mdb.ID) > -1)
					})
				}
			}, can.ui.content)
		}, function() {
			delete(target._tabs), can.onmotion.cacheClear(can, key, can.ui.content)
		})
	},
	layout: function(can) { can.ui.layout(can.ConfHeight(), can.ConfWidth(), 0, function(height, width) {
		can.ui._content_plugin && can.ui._content_plugin.onimport.size(can.ui._content_plugin, height, width, false)
	})},
})
