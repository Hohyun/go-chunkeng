package util

import (
	"fmt"
	"os"
)

func GetMysqlDsn() string {
	mysql_host := os.Getenv("MYSQL_HOST")
	mysql_port := os.Getenv("MYSQL_PORT")
	mysql_user := os.Getenv("MYSQL_USER")
	mysql_password := os.Getenv("MYSQL_PASSWORD")
	mysql_dbname := os.Getenv("MYSQL_DBNAME")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", mysql_user, mysql_password, mysql_host, mysql_port, mysql_dbname)
}
