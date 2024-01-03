package artist

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Artist struct {
	// all struct members must start with capital letter in order to be available for import in other packages
	Id              string `json:"id"`
	Name            string `json:"name"`
	CountryOfOrigin string `json:"country"`
	Age             int32  `json:"age"`
}

func CreateArtist(driver neo4j.Driver, newArtist Artist) (*Artist, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"CREATE (a:Artist {Name:$name, CountryOfOrigin:$country, Age:$age}) RETURN a",
		map[string]interface{}{"name": newArtist.Name, "country": newArtist.CountryOfOrigin, "age": newArtist.Age},
	)
	if err != nil {
		return nil, err
	}

	//	TODO: extract id from just created node
	for result.Next() {

	}

	return &Artist{Id: "0", Name: newArtist.Name, CountryOfOrigin: newArtist.CountryOfOrigin, Age: newArtist.Age}, nil
}
