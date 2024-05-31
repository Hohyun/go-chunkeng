package score

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Hohyun/go-chunkeng/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type ClassInfo struct {
	LevelCode string `json:"levelcode"`
	LevelName string `json:"levelname"`
	ClassCode string `json:"classcode"`
	ClassName string `json:"classname"`
}

// type ClassTreeInfo struct {
// 	ID     string `json:"id"`
// 	Parent string `json:"parent"`
// 	Text   string `json:"text"`
// }

type ClassTreeP struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Children []ClassTreeC `json:"children"`
}

type ClassTreeC struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetClasses(c *fiber.Ctx) error {
	rr := getClassesInfo()
	b, err := json.Marshal(rr)
	if err != nil {
		fmt.Printf("log.Logger: %v\n", err.Error())
	}

	c.Set("Content-type", "application/json")
	return c.Send(b)
}

func GetClassGroups(c *fiber.Ctx) error {
	data := GetClassesTreeData()

	// return c.SendString("Hello World 👋! This is great. Isn't it?")
	return c.Render("class_groups", fiber.Map{
		"Title": "Select Group",
		"Groups": data,
	})
}


func GetClassTeams(c *fiber.Ctx) error {
	group_id := c.Params("group_id")

	data := GetClassesTreeData()

	var teams []ClassTreeC
	for i := range data {
		if data[i].ID == group_id {
			teams = data[i].Children 
		}
	}

	// return c.SendString("Hello World 👋! This is great. Isn't it?")
	return c.Render("class_teams", fiber.Map{
		"Title": "Select Team",
		"Teams": teams,
	})
}

func GetClassesTree(c *fiber.Ctx) error {
	pp := GetClassesTreeData()

	b, err := json.Marshal(pp)
	if err != nil {
		fmt.Printf("log.Logger: %v\n", err.Error())
	}

	c.Set("Content-type", "application/json")
	return c.Send(b)
}

func GetClassesTreeData() []ClassTreeP {
	classes := getClassesInfo()

	var pp []ClassTreeP
	for i := range classes {
		if !parentExist(classes[i], pp) {
			p := ClassTreeP{
				ID:   classes[i].LevelCode,
				Name: classes[i].LevelName,
			}
			pp = append(pp, p)
		}
	}

	for i := range pp {
		cc := getChildren(pp[i], classes)
		pp[i].Children = cc
	}

	return pp
}

func parentExist(item ClassInfo, array []ClassTreeP) bool {
	for i := range array {
		if item.LevelCode == array[i].ID {
			return true
		}
	}
	return false
}

func getChildren(item ClassTreeP, array []ClassInfo) []ClassTreeC {
	var cc []ClassTreeC
	for i := range array {
		if array[i].LevelCode == item.ID {
			c := ClassTreeC{
				ID:   array[i].ClassCode,
				Name: array[i].ClassName,
			}
			cc = append(cc, c)
		}
	}
	return cc
}

func getClassesInfo() []ClassInfo {
	dsn := util.GetMysqlDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	queryString := `
	select aa.level_code, aa.level_name, bb.class_code, bb.class_name
	from (
		select d_category as level_code, d_name as level_name
		from ceng_directory
		where c_idx='1' and d_deleted='N' and d_group='C' and length(d_category)=3
	) aa join (
		select left(d_category, 3) as level_code, d_category as class_code,d_name as class_name
		from ceng_directory
		where c_idx='1' and d_deleted='N' and d_group='C' and length(d_category)=6
	) bb on aa.level_code = bb.level_code
	order by aa.level_name, bb.class_name
	`

	var rr []ClassInfo
	var r ClassInfo
	rows, err := db.Query(queryString)
	for rows.Next() {
		err := rows.Scan(&r.LevelCode, &r.LevelName, &r.ClassCode, &r.ClassName)
		if err != nil {
			fmt.Printf("log.Logger: %v\n", err.Error())
			return rr
		}
		rr = append(rr, r)
	}
	return rr
}

/*
func GetClassesTreeData(c *fiber.Ctx) error {
	classes := getClassesInfo()

	var rr []ClassTreeInfo
	var r ClassTreeInfo
	var parent string
	for i := range classes {
		parent = parentID(classes[i], rr)
		if parent == "0" {
			r = ClassTreeInfo{
				ID:     classes[i].LevelCode,
				Parent: parent,
				Text:   classes[i].LevelName,
			}
			parent = r.ID
			rr = append(rr, r)
		}
		r = ClassTreeInfo{
			ID:     classes[i].ClassCode,
			Parent: parent,
			Text:   classes[i].ClassName,
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

func parentID(item ClassInfo, array []ClassTreeInfo) string {
	for i := range array {
		if item.LevelCode == array[i].ID {
			return item.LevelCode
		}
	}
	return "0"
}
*/
