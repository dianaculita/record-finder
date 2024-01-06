package record

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

type Record struct {
	Id    string  `json:"id"`
	Title string  `json:"title"`
	Year  int64   `json:"year"`
	Price float64 `json:"price"`
}

const (
	GetAllRecordsCypher = `
		MATCH (r:Record) RETURN r
		`
	CreateRecordCypher = `
		CREATE (r:Record {Title: $title, Year:$year, Price:$price}) RETURN r
		`
	SearchRecordsByTitleCypher = `
		MATCH (r:Record) WHERE r.Title CONTAINS $title RETURN r
		`
)

type RecordRepository interface {
	GetAll(ctx context.Context) ([]*Record, error) //unfiltered
	CreateRecord(ctx context.Context, r Record) error
	SearchRecordsByTitle(ctx context.Context, title string) ([]*Record, error)
}

type Neo4jRecordRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jRecordRepository(driver neo4j.DriverWithContext) RecordRepository {
	return &Neo4jRecordRepository{
		driver: driver,
	}
}

func (recordRepo *Neo4jRecordRepository) GetAll(ctx context.Context) ([]*Record, error) {
	session := recordRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s", GetAllRecordsCypher)
	result, err := session.Run(ctx, GetAllRecordsCypher, map[string]interface{}{})

	if err != nil {
		return nil, fmt.Errorf("Error executing cypher query: %s", GetAllRecordsCypher)
	}

	records := []*Record{}

	for result.Next(ctx) {
		record := result.Record()
		nodeId := record.AsMap()["r"].(neo4j.Node).GetElementId()
		nodeProps := record.AsMap()["r"].(neo4j.Node).Props
		records = append(records, getRecordNodeProps(nodeProps, nodeId))
	}

	return records, nil
}

func (recordRepo *Neo4jRecordRepository) CreateRecord(ctx context.Context, r Record) error {
	session := recordRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s", CreateRecordCypher)
	_, err = session.Run(ctx, CreateRecordCypher,
		map[string]interface{}{"title": r.Title, "year": r.Year, "price": r.Price},
	)

	if err != nil {
		return err
	}

	return nil
}

func getRecordNodeProps(props map[string]any, nodeId string) *Record {
	price := props["Price"].(float64)
	title := props["Title"].(string)
	year := props["Year"].(int64)
	return &Record{
		nodeId, title, year, price,
	}
}

func (recordRepo *Neo4jRecordRepository) SearchRecordsByTitle(ctx context.Context, title string) ([]*Record, error) {
	session := recordRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s", SearchRecordsByTitleCypher)
	result, err := session.Run(ctx, SearchRecordsByTitleCypher,
		map[string]interface{}{"title": title},
	)

	if err != nil {
		return nil, fmt.Errorf("Error executing cypher query: %s", SearchRecordsByTitleCypher)
	}

	records := []*Record{}

	for result.Next(ctx) {
		record := result.Record()
		nodeId := record.AsMap()["r"].(neo4j.Node).GetElementId()
		nodeProps := record.AsMap()["r"].(neo4j.Node).Props
		records = append(records, getRecordNodeProps(nodeProps, nodeId))
	}

	return records, nil
}
