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
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	defer db.Close()

	team := c.Query("team", "")
	subject := c.Query("subject", "")
	fromDate := c.Query("fromDate", "")
	toDate := c.Query("toDate", "")

	cond := ""
	if team != "" {
		cond = fmt.Sprintf(" and class_id = '%s'", team)
	}
	if subject != "" {
		cond += fmt.Sprintf(" and subject = '%s'", subject)
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
from ceng_test_chaewoom
where true`

	if cond != "" {
		queryString += cond
	}

	queryString += " order by test_date desc, member_name"

	var rr []Chaewoom
	var r Chaewoom
	rows, err := db.Query(queryString)
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

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
	// workerId := c.Query("workerId", "")

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

func NewChaewoom(c *fiber.Ctx) error {
	cw := new(Chaewoom)

	if err := c.BodyParser(cw); err != nil {
		fiberlog.Error(err)
		return err
	}

	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	defer db.Close()

	queryString1 := `
	insert into ceng_test_chaewoom (
		test_date, class_id, class_name, subject, test_name, 
    	member_id, member_name, err_cnt, due_date, homeworks, 
		done, remarks, mod_id, mod_date)
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	result, err := db.Exec(queryString1,
		cw.TestDate, cw.ClassID, cw.ClassName, cw.Subject, cw.TestName,
		cw.MemberID, cw.MemberName, cw.ErrCnt, cw.DueDate, cw.Homeworks,
		cw.Done, cw.Remarks, cw.ModID, cw.TestDate)

	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result":      "OK",
		"description": fmt.Sprintf("Record (id: %d) inserted succefully", lastId),
	})
}

func UpdateChaewoom(c *fiber.Ctx) error {
	cw := new(Chaewoom)
	workerId := c.Query("workerId", "")

	if err := c.BodyParser(cw); err != nil {
		fiberlog.Error(err)
		return err
	}

	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	defer db.Close()

	queryString := `
	update ceng_test_chaewoom 
  set due_date = ?, homeworks = ?, done = ?, remarks = ?, mod_id = ?, mod_date = NOW() 
  where id = ?
	`
	_, err = db.Exec(queryString, cw.DueDate, cw.Homeworks, cw.Done, cw.Remarks, workerId, cw.ID)

	if err != nil {
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"result":      "OK",
		"description": fmt.Sprintf("Record (id: %d) updated succefully", cw.ID),
	})
}
