Volcanos(chat.ONSYNTAX, {
	sql: {
		prefix: {"--": code.COMMENT},
		regexp: {
			"^[A-Z0-9_]+$": code.KEYWORD,
		},
		prepare: kit.Dict(
			code.KEYWORD, [
				"CREATE", "DATABASE", "TABLE", "IF",
				"DROP",
				"SHOW", "DATABASES", "TABLES", "COLUMNS",
				"USE",
				"ALTER", "ADD",
				"DESCRIBE",
				"RENAME", "TO",
				
				"INSERT", "INTO",
				"DELETE", "FROM",
				"UPDATE", "SET",
				"SELECT", "FROM", "WHERE",
				"GROUP", "BY",
				"ORDER", "BY",
				"HAVING",
				"OFFSET", "LIMIT",
				"INNER", "JOIN", "ON",
				"LEFT", "JOIN", "ON",
				"RIGHT", "JOIN", "ON",
				"CROSS", "JOIN",
			],
			code.CONSTANT, [
				"InnoDB", "utf8mb4",
				"EXISTS", "NULL",
			],
			code.DATATYPE, [
				"TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT",
				"FLOAT", "DOUBLE", "DECIMAL",
				"DATE", "TIME", "DATETIME", "TIMESTAMP",
				"CHAR", "VARCHAR", "TEXT", "BLOB",
			],
			code.FUNCTION, [
				"ENGINE", "CHARSET",
				
				"AUTO_INCREMENT",
				"PRIMARY", "KEY",
				"UNIQUE", "INDEX",
				"DEFAULT", "COMMENT",
				"VALUES",
				
				"DISTINCT",
				"SUM", "AVG", "MAX", "MIN",
				"COUNT",
				"CONCAT",
				"UPPER",
				"DATE",
				
				"IS", "NOT", "AND", "OR",
				"LIKE", "GLOB",
				"IN", "BETWEEN",
				
				"ASC", "DESC",
			],
		),
	},
})
