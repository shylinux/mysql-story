package client

import (
	"shylinux.com/x/ice"
)

type client struct {
	ice.Hash
	short  string `data:"sess"`
	field  string `data:"time,sess,username,password,host,port,driver,database"`
	create string `name:"create sess*=biz username*=root password*=root host*=127.0.0.1 port*=10001 database*=mysql driver*=mysql"`
	list   string `name:"list sess auto"`
}

func init() { ice.CodeModCmd(client{}) }

type Client struct{ client }

func init() { ice.CodeModCmd(Client{}) }
