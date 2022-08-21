module shylinux.com/x/mysql-story

go 1.11

require shylinux.com/x/go-sql-mysql v0.0.1

require (
	shylinux.com/x/ice v1.0.8
	shylinux.com/x/icebergs v1.3.6
	shylinux.com/x/toolkits v0.6.8
)

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)
