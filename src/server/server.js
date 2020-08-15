Volcanos("onimport", {help: "导入数据", list: [],
    _init: function(can, msg, list, cb, target) {
        can.onappend.table(can, target, "table", msg)
        can.onappend.board(can, target, "board", msg)
        return typeof cb == "function" && cb(msg)
    },
})
Volcanos("onaction", {help: "交互操作", list: [],
    _init: function(can, msg, list, cb, target) {},
})
Volcanos("onexport", {help: "导出数据", list: [],
    _init: function(can, msg, list, cb, target) {},
})
