module shylinux.com/x/mysql-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

replace shylinux.com/x/go-sql-mysql => ./usr/go-sql-mysql

require (
	shylinux.com/x/ice v1.4.7
	shylinux.com/x/icebergs v1.8.7
	shylinux.com/x/toolkits v1.0.2
)

require (
	github.com/mattn/go-sqlite3 v1.14.16
	shylinux.com/x/go-sql-mysql v0.0.2
)
