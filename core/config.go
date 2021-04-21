package core

import "os"

const (
	driver       = "mysql"
	hostname     = "127.0.0.1:3306"
	databaseName = "queue"
	table        = "Queue"
	idField      = "id"
	dataField    = "data"
)

var (
	username = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASS")
)
