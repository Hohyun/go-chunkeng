package score

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Hohyun/go-chunkeng/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
)

type Score struct {
	ID         int64     `json:"id"`
	TestDate   time.Time `json:"test_date"`
	TestName   string    `json:"test_name"`
	MemberID   string    `json:"member_id"`
	MemberName string    `json:"member_name"`
	Team       string    `json:"team"`
	Subject    string    `json:"subject"`
	Teacher    string    `json:"teacher"`
	ErrCnt     int64     `json:"err_cnt"`
	TtlCnt     int64     `json:"ttl_cnt"`
	Chaewoom   bool      `json:"chaewoom"`
	RegID      string    `json:"reg_id"`
	RegDate    time.Time `json:"reg_date"`
	ModID      string    `json:"mod_id"`
	ModDate    time.Time `json:"mod_date"`
	Remarks    string    `json:"remarks"`
}

type ScoreWithChaewoom struct {
	Scores           []Score `json:"scores"`
	ChaewoomCriteria int64   `json:"chaewoomCriteria"`
}

func GetScores(c *fiber.Ctx) error {
	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	checkError(err)

	defer db.Close()

	queryString := `
	select id, test_date, test_name, member_id, member_name, 
		team, subject, teacher, err_cnt, ttl_cnt, chaewoom, 
		reg_id, reg_date, mod_id, mod_date, remarks
	from ceng_test_score
	where test_date > DATE_SUB(NOW(), INTERVAL 30 DAY)
	order by test_date desc, team, test_name, member_name 
	`

	var rr []Score
	var r Score
	rows, err := db.Query(queryString)
	checkError(err)

	for rows.Next() {
		err := rows.Scan(&r.ID, &r.TestDate, &r.TestName, &r.MemberID, &r.MemberName,
			&r.Team, &r.Subject, &r.Teacher, &r.ErrCnt, &r.TtlCnt, &r.Chaewoom,
			&r.RegID, &r.RegDate, &r.ModID, &r.ModDate, &r.Remarks)
		if err != nil {
			fmt.Printf("log.Logger: %v\n", err.Error())
			return c.SendString(err.Error())
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

func GetScore(c *fiber.Ctx) error {
	return c.SendString("Single score")
}

func NewScore(c *fiber.Ctx) error {
	p := new(ScoreWithChaewoom)

	if err := c.BodyParser(p); err != nil {
		fiberlog.Error(err)
		return err
	}

	// log.Println(p)

	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	checkError(err)

	defer db.Close()

	queryString := `
	insert into ceng_test_score (test_date, test_name, member_id, member_name, 
		team, subject, teacher, err_cnt, ttl_cnt, chaewoom, 
		reg_id, reg_date, remarks)
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	var chaewoom bool
	for _, s := range p.Scores {
		cwCriteria := p.ChaewoomCriteria
		if s.Chaewoom == true || s.ErrCnt >= cwCriteria {
			chaewoom = true
		} else {
			chaewoom = false
		}
		_, err = tx.Exec(queryString, s.TestDate, s.TestName, s.MemberID, s.MemberName,
			s.Team, s.Subject, s.Teacher, s.ErrCnt, s.TtlCnt, chaewoom,
			s.RegID, s.RegDate, s.Remarks)
		checkError(err)

		if err != nil {
			return c.JSON(fiber.Map{
				"result":      "FAIL",
				"description": err.Error(),
			})
		}
	}
	err = tx.Commit()
	checkError(err)

	return c.JSON(fiber.Map{
		"result":      "OK",
		"description": fmt.Sprintf("%v records inserted succefully", len(p.Scores)),
	})
}

func DeleteScore(c *fiber.Ctx) error {
	return c.SendString("Delete score")
}

func checkError(err error) {
	if err != nil {
		// panic(err)
		fiberlog.Warn("Error: ", err)
	}
}
