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
	MemberID    string `json:"member_id"`
	MemberName  string `json:"member_name"`
	ClassID     string `json:"class_id"`
  ClassName   string `json:"class_name"`
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
select m_id member_id, m_name member_name, m_class class_id, class_name
from ceng_member cm join ceng_class_info cci on cm.m_class = cci.class_code 
where c_idx = '1' and m_out = 'N' and m_level = '1' and m_class = ?
order by m_name
	`

	var rr []ClassMember
	var r ClassMember
	rows, err := db.Query(queryString, class)
	for rows.Next() {
		err := rows.Scan(&r.MemberID, &r.MemberName, &r.ClassID, &r.ClassName)
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

func GetMember(class_id string , member_name string) string {
  dsn := util.GetMysqlDsn()
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    panic(err)
  }
  defer db.Close()

	queryString := `
select m_id member_id, m_name member_name, m_class class_id, class_name
from ceng_member cm join ceng_class_info cci on cm.m_class = cci.class_code 
where c_idx = '1' and m_out = 'N' and m_level = '1' and m_class = ? and m_name = ?
	`

	var rr []ClassMember
	var r ClassMember
	rows, err := db.Query(queryString, class_id, member_name)
	for rows.Next() {
		err := rows.Scan(&r.MemberID, &r.MemberName, &r.ClassID, &r.ClassName)
		if err != nil {
			fmt.Printf("log.Logger: %v\n", err.Error())
      return ""
		}
		rr = append(rr, r)
	}

  if len(rr) == 0 {
    return ""
  }

	return rr[0].MemberID
}
