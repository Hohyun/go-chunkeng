package score

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hohyun/go-chunkeng/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type Score struct {
	ID         int    `json:"id"`
	TestDate   time.Time `json:"test_date"`
	TestName   string `json:"test_name"`
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
	Team       string `json:"team"`
	Subject    string `json:"subject"`
	Teacher    string `json:"teacher"`
	ErrCnt     int    `json:"err_cnt"`
	TtlCnt     int    `json:"ttl_cnt"`
	Chaewoom   bool   `json:"chaewoom"`
	RegID      string `json:"reg_id"`
	RegDate    time.Time  `json:"reg_date"`
	ModDate    time.Time  `json:"mod_date"`
	Remarks    string `json:"remarks"`
}

func GetScores(c *fiber.Ctx) error {
	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	queryString := `
	select id, test_date, test_name, member_id, member_name, team, subject, teacher, err_cnt, ttl_cnt, chaewoom, reg_id, reg_date, mod_date, remarks
	from ceng_test_score
	where t_date > DATE_SUB(NOW(), INTERVAL 30 DAY)
	order by t_date desc, team, test_name, member_name 
	`

	var rr []Score
	var r Score
	rows, err := db.Query(queryString)
	for rows.Next() {
		err := rows.Scan(&r.ID, &r.TestDate, &r.TestName, &r.MemberID, &r.MemberName, &r.Team, &r.Subject, &r.Teacher, 
			&r.ErrCnt, &r.TtlCnt, &r.Chaewoom, &r.RegID, &r.RegDate, &r.ModDate, &r.Remarks)
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
