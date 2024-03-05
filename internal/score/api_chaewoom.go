
package score

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Hohyun/go-chunkeng/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
)

type Chaewoom struct {
	ID         int64  `json:"id"`
	TestDate   string `json:"test_date"`
	ClassID    string `json:"class_id"`
	ClassName  string `json:"class_name"`
	Subject    string `json:"subject"`
	TestName   string `json:"test_name"`
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
	ErrCnt     int64  `json:"err_cnt"`
	DueDate    string `json:"due_date"`
  Homeworks  string `json:"homeworks"`
	Done       bool   `json:"done"`
	Remarks    string `json:"remarks"`
	ModID      string `json:"mod_id"`
	ModDate    string `json:"mod_date"`
}

func GetChaewooms(c *fiber.Ctx) error {
	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	checkError(err)

	defer db.Close()

	team := c.Query("team", "")
	subject := c.Query("subject", "")
	testName := c.Query("testName", "")
	fromDate := c.Query("fromDate", "")
	toDate := c.Query("toDate", "")

	cond := ""
	if team != "" {
		cond = fmt.Sprintf(" and class_id = '%s'", team)
	}
	if subject != "" {
		cond += fmt.Sprintf(" and subject = '%s'", subject)
	}
	if testName != "" {
		cond += fmt.Sprintf(" and test_name = '%s'", testName)
	}
	if fromDate != "" {
		cond += fmt.Sprintf(" and test_date >= '%s'", fromDate)
	}
	if toDate != "" {
		cond += fmt.Sprintf(" and test_date <= '%s'", toDate)
	}
	if fromDate == "" && toDate == "" {
		cond += " and test_date > DATE_SUB(NOW(), INTERVAL 30 DAY)"
	}

	queryString := `
select id, DATE_FORMAT(test_date,"%Y-%m-%dT%T.000Z"), class_id, class_name, 
  subject, test_name, member_id, member_name, err_cnt, 
  coalesce(DATE_FORMAT(due_date,"%Y-%m-%dT%T.000Z"), ""), homeworks, done, 
  coalesce(remarks, ""), coalesce(mod_id, ""), 
  coalesce(DATE_FORMAT(mod_date,"%Y-%m-%dT%T.000Z"), "")
from ceng_chaewoom_info
where true`

	if cond != "" {
		queryString += cond
	}

	var rr []Chaewoom
	var r Chaewoom
	rows, err := db.Query(queryString)
	checkError(err)

	for rows.Next() {
		err := rows.Scan(&r.ID, &r.TestDate, &r.ClassID, &r.ClassName, 
      &r.Subject, &r.TestName, &r.MemberID, &r.MemberName, &r.ErrCnt, 
      &r.DueDate, &r.Homeworks, &r.Done, &r.Remarks, &r.ModID, &r.ModDate)

		if err != nil {
      fiberlog.Error(err.Error())
	    return c.JSON(fiber.Map{
		    "result":      "FAIL",
		    "description": err.Error(),
	    })
		}
		rr = append(rr, r)
	}

	b, err := json.Marshal(rr)
	if err != nil {
		fmt.Printf("log.Logger: %v\n", err.Error())
	}

	c.Set("Content-type", "application/json")
	return c.Send(b)
}


func DeleteChaewoom(c *fiber.Ctx) error {
	id := c.Params("id")

	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	queryString := `
delete from ceng_test_chaewoom where id = ?
	`
	_, err = db.Exec(queryString, id)
	if err != nil {
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"result":      "OK",
		"description": fmt.Sprintf("Record %s deleted succefully", id),
	})
}

