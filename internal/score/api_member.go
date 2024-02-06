package score

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Hohyun/go-chunkeng/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type ClassMember struct {
	CmID    string `json:"id"`
	CmName  string `json:"name"`
	CmClass string `json:"class"`
}

func GetMembers(c *fiber.Ctx) error {
	class := c.Params("class")

	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	queryString := `
	select m_id, m_name, m_class
	from ceng_member cm 
	where c_idx = '1' and m_out = 'N' and m_level = '1' and m_class = ?
	order by m_name 
	`

	var rr []ClassMember
	var r ClassMember
	rows, err := db.Query(queryString, class)
	for rows.Next() {
		err := rows.Scan(&r.CmID, &r.CmName, &r.CmClass)
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
