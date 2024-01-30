module shylinux.com/x/mysql-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	shylinux.com/x/ice v1.4.5
	shylinux.com/x/icebergs v1.8.5
	shylinux.com/x/toolkits v1.0.0
)

require (
	github.com/mattn/go-sqlite3 v1.14.16
	shylinux.com/x/go-sql-mysql v0.0.2
)
