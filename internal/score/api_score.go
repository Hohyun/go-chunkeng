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

type Score struct {
	ID         int64  `json:"id"`
	TestDate   string `json:"test_date"`
	TestName   string `json:"test_name"`
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
	ClassID    string `json:"class_id"`
	ClassName  string `json:"class_name"`
	Subject    string `json:"subject"`
	Teacher    string `json:"teacher"`
	ErrCnt     int64  `json:"err_cnt"`
	TtlCnt     int64  `json:"ttl_cnt"`
	Chaewoom   bool   `json:"chaewoom"`
	RegID      string `json:"reg_id"`
	RegDate    string `json:"reg_date"`
	ModID      string `json:"mod_id"`
	ModDate    string `json:"mod_date"`
	Remarks    string `json:"remarks"`
}

type ScoreWithChaewoom struct {
	Scores           []Score `json:"scores"`
	ChaewoomCriteria int64   `json:"chaewoomCriteria"`
}

func GetScores(c *fiber.Ctx) error {
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
	teacher := c.Query("teacher", "")
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
	if teacher != "" {
		cond += fmt.Sprintf(" and teacher = '%s'", teacher)
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
select id, DATE_FORMAT(test_date,"%Y-%m-%dT%T.000Z"), test_name, member_id, 
  member_name, class_id, class_name, subject, teacher, err_cnt, ttl_cnt, chaewoom, reg_id, 
  DATE_FORMAT(reg_date,"%Y-%m-%dT%T.000Z"), coalesce(mod_id, ""), 
  coalesce(DATE_FORMAT(mod_date,"%Y-%m-%dT%T.000Z"), ""), remarks
from ceng_test_score 
where true`

	if cond != "" {
		queryString += cond
	}
	queryString += "\norder by test_date desc, class_name, test_name, member_name;"
	// where test_date > DATE_SUB(NOW(), INTERVAL 30 DAY)
	// fmt.Println(queryString)

	var rr []Score
	var r Score
	rows, err := db.Query(queryString)
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	for rows.Next() {
		err := rows.Scan(&r.ID, &r.TestDate, &r.TestName, &r.MemberID, &r.MemberName,
			&r.ClassID, &r.ClassName, &r.Subject, &r.Teacher, &r.ErrCnt, &r.TtlCnt, &r.Chaewoom,
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
  workerId := c.Query("workerId", "")

	if err := c.BodyParser(p); err != nil {
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
	insert into ceng_test_score (test_date, test_name, member_id, member_name, 
		class_id, class_name, subject, teacher, err_cnt, ttl_cnt, chaewoom, 
		reg_id, reg_date, remarks)
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	queryString2 := `
	insert into ceng_test_chaewoom (score_id, homeworks, mod_date)
	values (?, ?, ?)
	`

	tx, err := db.Begin()
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	defer tx.Rollback()

	var chaewoom bool
	for _, s := range p.Scores {
		cwCriteria := p.ChaewoomCriteria
		if s.Chaewoom || s.ErrCnt >= cwCriteria {
			chaewoom = true
		} else {
			chaewoom = false
		}

		var member_id string
		if s.MemberID == "" {
			member_id = GetMember(s.ClassID, s.MemberName)
		} else {
			member_id = s.MemberID
		}

		result, err := tx.Exec(queryString1, s.TestDate, s.TestName, member_id, s.MemberName,
			s.ClassID, s.ClassName, s.Subject, s.Teacher, s.ErrCnt, s.TtlCnt, chaewoom,
			workerId, s.RegDate, s.Remarks)
	
    if err != nil {
		  fiberlog.Error(err)
		  return c.JSON(fiber.Map{
			  "result":      "FAIL",
			  "description": err.Error(),
		  })
	  }
		
		lastID, err := result.LastInsertId()
		
	  if err != nil {
		  fiberlog.Error(err)
		  return c.JSON(fiber.Map{
			  "result":      "FAIL",
			  "description": err.Error(),
		  })
	  }

    if chaewoom {
		  _, err = tx.Exec(queryString2, lastID, "", s.TestDate)
	    if err != nil {
		    fiberlog.Error(err)
		    return c.JSON(fiber.Map{
			    "result":      "FAIL",
			    "description": err.Error(),
		    })
	    }
    }
	}

	err = tx.Commit()
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}
	
	if err != nil {
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result":      "OK",
		"description": fmt.Sprintf("%v records inserted succefully", len(p.Scores)),
	})
}

func DeleteScore(c *fiber.Ctx) error {
	id := c.Params("id")
  // workerId := c.Query("workerId", "")

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
delete from ceng_test_score where id = ?
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

func UpdateScore(c *fiber.Ctx) error {
	id := c.Params("id")
	workerId := c.Query("workerId", "")

	s := new(Score)

	if err := c.BodyParser(s); err != nil {
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
  update ceng_test_score 
  set err_cnt = ?, chaewoom = ?, remarks = ?, mod_id = ?, mod_date = now()
  where id = ?
  `

	tx, err := db.Begin()
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}

	defer tx.Rollback()

	_, err = tx.Exec(queryString, s.ErrCnt, s.Chaewoom, s.Remarks, workerId, id)
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}
	
	chaewoomInsert := false
	if s.Chaewoom {
		queryString = `
    	select * from ceng_test_chaewoom
    	where score_id = ?
    	`
		rows, err := db.Query(queryString, id)
    if err != nil {
		  fiberlog.Error(err)
		  return c.JSON(fiber.Map{
			  "result":      "FAIL",
			  "description": err.Error(),
		  })
	  }

		// if rows is null, insert a new record
		if !rows.Next() {
			queryString = `
		insert into ceng_test_chaewoom (score_id, homeworks, mod_date) 
		values (?, ?, now())
		`
			_, err = tx.Exec(queryString, s.ID, "")
	    if err != nil {
		    fiberlog.Error(err)
		    return c.JSON(fiber.Map{
			    "result":      "FAIL",
			    "description": err.Error(),
		    })
	    }
			
			chaewoomInsert = true
		}
	}
	err = tx.Commit()
	if err != nil {
		fiberlog.Error(err)
		return c.JSON(fiber.Map{
			"result":      "FAIL",
			"description": err.Error(),
		})
	}
	
	description := ""
	if chaewoomInsert {
		description = " + (a 채움장 record created)"
	}

	return c.JSON(fiber.Map{
		"result":      "OK",
		"description": fmt.Sprintf("A score record (id: %v) updated successfully", s.ID) + description,
	})
}

