package label

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

type Label struct {
	// all struct members must start with capital letter in order to be available for import in other packages
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

const (
	GetAllLabelsCypher = `
		MATCH (l:Label) RETURN l
		`
	CreateLabelCypher = `
		CREATE (l:Label {Name:$name, Address:$address}) RETURN l
		`
	SearchLabelsByNameCypher = `
		MATCH (l:Label) WHERE l.Name CONTAINS $name RETURN l
		`
)

type LabelRepository interface {
	GetAll(ctx context.Context) ([]*Label, error)
	CreateLabel(ctx context.Context, l Label) error
	SearchLabelsByName(ctx context.Context, name string) ([]*Label, error)
}

type Neo4jLabelRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jLabelRepository(driver neo4j.DriverWithContext) LabelRepository {
	return &Neo4jLabelRepository{
		driver: driver,
	}
}

func (labelRepo *Neo4jLabelRepository) GetAll(ctx context.Context) ([]*Label, error) {
	session := labelRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s")
	result, err := session.Run(ctx, GetAllLabelsCypher, map[string]interface{}{})

	if err != nil {
		return nil, fmt.Errorf("Error executing cypher query: %s", GetAllLabelsCypher)
	}

	labels := []*Label{}

	for result.Next(ctx) {
		record := result.Record()
		nodeId := record.AsMap()["l"].(neo4j.Node).GetElementId()
		nodeProps := record.AsMap()["l"].(neo4j.Node).Props
		labels = append(labels, getLabelNodeProps(nodeProps, nodeId))
	}

	return labels, nil
}

func (labelRepo *Neo4jLabelRepository) CreateLabel(ctx context.Context, l Label) error {
	session := labelRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s", CreateLabelCypher)
	_, err = session.Run(ctx, CreateLabelCypher,
		map[string]interface{}{"name": l.Name, "address": l.Address},
	)

	if err != nil {
		return err
	}

	return nil
}

func getLabelNodeProps(props map[string]any, nodeId string) *Label {
	name := props["Name"].(string)
	address := props["Address"].(string)
	return &Label{
		nodeId, name, address,
	}
}

func (labelRepo *Neo4jLabelRepository) SearchLabelsByName(ctx context.Context, name string) ([]*Label, error) {
	session := labelRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s", SearchLabelsByNameCypher)
	result, err := session.Run(ctx, SearchLabelsByNameCypher,
		map[string]interface{}{"name": name},
	)

	if err != nil {
		return nil, fmt.Errorf("Error executing cypher query: %s", SearchLabelsByNameCypher)
	}

	labels := []*Label{}

	for result.Next(ctx) {
		record := result.Record()
		nodeId := record.AsMap()["l"].(neo4j.Node).GetElementId()
		nodeProps := record.AsMap()["l"].(neo4j.Node).Props
		labels = append(labels, getLabelNodeProps(nodeProps, nodeId))
	}

	return labels, nil
}
