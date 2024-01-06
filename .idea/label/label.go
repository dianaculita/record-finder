package label

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Label struct {
	// all struct members must start with capital letter in order to be available for import in other packages
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func CreateLabel(driver neo4j.Driver, newLabel Label) (*Label, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"CREATE (l:Label {Name:$name, Address:$address}) RETURN l",
		map[string]interface{}{"name": newLabel.Name, "address": newLabel.Address},
	)
	if err != nil {
		return nil, err
	}

	//	TODO: extract id from just created node
	for result.Next() {

	}

	return &Label{Id: "0", Name: newLabel.Name, Address: newLabel.Address}, nil
}
