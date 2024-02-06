package score

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Hohyun/go-chunkeng/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type Score struct {
	TestDate   string `json:"test_date"`
	TestName   string `json:"test_name"`
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
	Team       string `json:"team"`
	Subject    string `json:"subject"`
	Teacher    string `json:"teacher"`
	ErrorCount string `json:"err_cnt"`
	TotalCount string `json:"ttl_cnt"`
	Chaewoom   bool   `json:"chaewoom"`
	Remarks    string `json:"remarks"`
}

func GetScores(c *fiber.Ctx) error {
	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	queryString := `
	select t_date, t_name, m_id, m_name, team, subject, teacher, err_cnt, ttl_cnt, chaewoom, remarks
	from ceng_test_score
	where t_date > DATE_SUB(NOW(), INTERVAL 30 DAY)
	order by t_date desc, team, subject, test_name, member_name 
	`

	var rr []Score
	var r Score
	rows, err := db.Query(queryString)
	for rows.Next() {
		err := rows.Scan(&r.TestDate, &r.MemberID, &r.MemberName, &r.Team, &r.Subject, &r.Teacher, &r.ErrorCount, &r.TotalCount, &r.Chaewoom, &r.Remarks)
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
	return c.SendString("New score")
}

func DeleteScore(c *fiber.Ctx) error {
	return c.SendString("Delete score")
}
