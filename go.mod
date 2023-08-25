module shylinux.com/x/mysql-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	shylinux.com/x/ice v1.3.11
	shylinux.com/x/icebergs v1.5.19
	shylinux.com/x/toolkits v0.7.10
)

require (
	github.com/elastic/go-elasticsearch v0.0.0 // indirect
	github.com/glebarez/sqlite v1.9.0 // indirect
	github.com/go-sqlite/sqlite3 v0.0.0-20180313105335-53dd8e640ee7 // indirect
	github.com/gonuts/binary v0.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.16
	gorm.io/gorm v1.25.4 // indirect
	shylinux.com/x/go-sql-mysql v0.0.2
)
