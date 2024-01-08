package artist

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

type Artist struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	CountryOfOrigin string `json:"country"`
	Age             int64  `json:"age"`
}

const (
	GetAllArtistsCypher = `
		MATCH (a:Artist) RETURN a
		`
	CreateArtistCypher = `
		CREATE (a:Artist {Name:$name, CountryOfOrigin:$country, Age:$age}) RETURN a
		`
	SearchArtistsByNameCypher = `
		MATCH (a:Artist) WHERE a.Name CONTAINS $name RETURN a
		`
)

type ArtistRepository interface {
	GetAll(ctx context.Context) ([]*Artist, error)
	CreateArtist(ctx context.Context, a Artist) error
	SearchArtistsByName(ctx context.Context, title string) ([]*Artist, error)
}

type Neo4jArtistRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jArtistRepository(driver neo4j.DriverWithContext) ArtistRepository {
	return &Neo4jArtistRepository{
		driver: driver,
	}
}

func (artistRepo *Neo4jArtistRepository) GetAll(ctx context.Context) ([]*Artist, error) {
	session := artistRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s")
	result, err := session.Run(ctx, GetAllArtistsCypher, map[string]interface{}{})

	if err != nil {
		return nil, fmt.Errorf("Error executing cypher query: %s", GetAllArtistsCypher)
	}

	records := []*Artist{}

	for result.Next(ctx) {
		record := result.Record()
		nodeId := record.AsMap()["a"].(neo4j.Node).GetElementId()
		nodeProps := record.AsMap()["a"].(neo4j.Node).Props
		records = append(records, getArtistNodeProps(nodeProps, nodeId))
	}

	return records, nil
}

func (artistRepo *Neo4jArtistRepository) CreateArtist(ctx context.Context, a Artist) error {
	session := artistRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s", CreateArtistCypher)
	_, err = session.Run(ctx, CreateArtistCypher,
		map[string]interface{}{"name": a.Name, "country": a.CountryOfOrigin, "age": a.Age},
	)

	if err != nil {
		return err
	}

	return nil
}

func getArtistNodeProps(props map[string]any, nodeId string) *Artist {
	name := props["Name"].(string)
	country := props["CountryOfOrigin"].(string)
	age := props["Age"].(int64)
	return &Artist{
		nodeId, name, country, age,
	}
}

func (artistRepo *Neo4jArtistRepository) SearchArtistsByName(ctx context.Context, name string) ([]*Artist, error) {
	session := artistRepo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	err := session.Close(ctx)
	if err != nil {
		fmt.Println("session.Close failed")
	}

	log.Printf("Attempt to run cypher query: %s", SearchArtistsByNameCypher)
	result, err := session.Run(ctx, SearchArtistsByNameCypher,
		map[string]interface{}{"name": name},
	)

	if err != nil {
		return nil, fmt.Errorf("Error executing cypher query: %s", SearchArtistsByNameCypher)
	}

	artists := []*Artist{}

	for result.Next(ctx) {
		record := result.Record()
		nodeId := record.AsMap()["a"].(neo4j.Node).GetElementId()
		nodeProps := record.AsMap()["a"].(neo4j.Node).Props
		artists = append(artists, getArtistNodeProps(nodeProps, nodeId))
	}

	return artists, nil
}
