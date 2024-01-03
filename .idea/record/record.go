package record

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Record struct {
	// all struct members must start with capital letter in order to be available for import in other packages
	Id    string  `json:"id"`
	Title string  `json:"title"`
	Year  int32   `json:"year"`
	Price float64 `json:"price"`
}

func CreateRecord(driver neo4j.Driver, newRecord Record) (*Record, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"CREATE (r:Record {Title: $title, Year:$year, Price:$price}) RETURN r",
		map[string]interface{}{"title": newRecord.Title, "year": newRecord.Year, "price": newRecord.Price},
	)
	if err != nil {
		return nil, err
	}

	//	TODO: extract id from just created node
	for result.Next() {

	}

	return &Record{Id: "0", Title: newRecord.Title, Year: newRecord.Year, Price: newRecord.Price}, nil
}

func GetAllRecords(driver neo4j.Driver) []Record {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"MATCH (r:Record) RETURN r",
		map[string]interface{}{},
	)
	if err != nil {
		return nil
	}

	var records = []Record{}

	for result.Next() {
		record := result.Record().GetByIndex(0).(neo4j.Node)
		fmt.Println(record)
		fmt.Println()
		name, _ := record.Props()["Artist"].(string)
		//properties := record.Props()
		fmt.Println(name)
		fmt.Println()
		//recDb := Record{
		//	Title: record.Get("title").(string),
		//	//Year:  int(record.Get("age").(int64)),
		//}
		//title, _ := res.Get("title")
		//year, _ := res.Get("year")
		//price, _ := res.Get("price")
		//var record = Record{Title: title, Year: year, Price: price}
	}

	return records
}
