Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) {
		can.ui = can.onappend.layout(can), can.onimport._project(can, msg)
	},
	_project: function(can, msg) {
		can.onimport.itemlist(can, msg.Table(function(value) {
			value.icon = icon.sess, value.nick = `${value.sess}(${value.host}:${value.port})`
			return value
		}), function(event, value, show, target) {
			show == undefined && can.run(event, [value.sess], function(msg) { can.onimport._database(can, msg, value.sess, target) })
		})
	},
	_database: function(can, msg, sess, target) {
		can.onimport.itemlist(can, [{icon: icon.path, nick: "_script", _hash: [sess, "_script"]}], function(event, value, show, target) {
			show == undefined && can.run(event, [nfs.SCRIPT, sess, nfs.SRC], function(msg) { can.onimport._script(can, msg, sess, target) })
		}, function() {}, target)
		can.onimport.itemlist(can, [{icon: icon.xterm, nick: "_shell", _hash: [sess, "_shell"]}], function(event, value, show, target) {
			can.onimport._content(can, {index: "web.code.mysql.shell", args: [sess], style: html.OUTPUT}, value, target)
		}, function() {}, target)
		can.onimport.itemlist(can, [{icon: "bi bi-person-lock", nick: "_grant", _hash: [sess, "_grant"]}], function(event, value, show, target) {
			can.onimport._content(can, {index: "web.code.mysql.grant", args: [sess]}, value, target)
		}, function() {}, target)
		can.onimport.itemlist(can, msg.Table(function(value, index) {
			value.icon = icon.database, value.nick = value.database
			value._hash = [sess, value.database]
			return value
		}), function(event, value, show, target) {
			show == undefined && can.run(event, [sess, value.database], function(msg) { can.onimport._table(can, msg, sess, value.database, target) })
		}, function() {}, target)
	},
	_script: function(can, msg, sess, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) {
			value.icon = icon.file, value.nick = value.file
			value._hash = [sess, "_script", value.file]
			return value
		}), function(event, value, show, target) {
			can.onimport._content(can, {index: "web.code.mysql.script", args: [sess, nfs.SRC, value.file]}, value, target)
		}, function() {}, target)
	},
	_table: function(can, msg, sess, database, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) {
			value.icon = icon.table, value.nick = `${value.table}(${value.total})`
			value._hash = [sess, database, value.table, "query"]
			return value
		}), function(event, value, show, target) {
			can.onimport._content(can, {index: "web.code.mysql.query", args: [sess, database, value.table]}, value, target)
		}, function() {}, target)
	},
	_content: function(can, meta, value, target) {
		if (!value._msg) { var msg = can.request({}); msg.Push(meta), value._msg = msg }
		return can.onimport.tabsCache(can, value, target, value._msg)
	},
	layout: function(can) { can.ui.layout(can.ConfHeight(), can.ConfWidth(), 0, function(height, width) {
		var sub = can.db.value._content_plugin; if (sub) {
			sub.onexport.output = function(_, msg) {
				can.page.Select(sub, sub._option, "div.item.text.id", function(target) {
					can.onmotion.toggle(can, target, msg.append && msg.append.indexOf(mdb.ID) > -1)
				})
			}
		}
	}) },
})