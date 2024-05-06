package score

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Hohyun/go-chunkeng/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type Framework struct {
	Tests     []string `json:"tests"`
	Teachers  []string `json:"teachers"`
	Subjects  []string `json:"subjects"`
	Homeworks []string `json:"homeworks"`
}

func GetFramework(c *fiber.Ctx) error {

	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	queryString := `
select keystr, value 
from ceng_test_framework
where keystr in ('tests', 'teachers', 'subjects', 'chaewoom:homeworks')
	`

	var r Framework
	var keystr string
	var value string
	rows, err := db.Query(queryString)
	if err != nil {
		fmt.Printf("log.Logger: %v\n", err.Error())
		return c.SendString(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&keystr, &value)
		if err != nil {
			fmt.Printf("log.Logger: %v\n", err.Error())
			return c.SendString(err.Error())
		}
		if keystr == "tests" {
			r.Tests = strings.Split(value, ",")
		} else if keystr == "teachers" {
			r.Teachers = strings.Split(value, ",")
		} else if keystr == "subjects" {
			r.Subjects = strings.Split(value, ",")
		} else if keystr == "chaewoom:homeworks" {
			r.Homeworks = strings.Split(value, ",")
		}
	}

	b, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("log.Logger: %v\n", err.Error())
	}

	c.Set("Content-type", "application/json")
	return c.Send(b)
}
