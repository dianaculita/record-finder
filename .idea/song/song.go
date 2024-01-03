package song

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Song struct {
	// all struct members must start with capital letter in order to be available for import in other packages
	Id          string `json:"id"`
	Name        string `json:"name"`
	DurationSec int64  `json:"duration"`
	Theme       string `json:"theme"`
}

func CreateSong(driver neo4j.Driver, newSong Song) (*Song, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"CREATE (s:Song {Name:$name, DurationSec:$duration, Theme:$theme}) RETURN s",
		map[string]interface{}{"name": newSong.Name, "duration": newSong.DurationSec, "theme": newSong.Theme},
	)
	if err != nil {
		return nil, err
	}

	//	TODO: extract id from just created node
	for result.Next() {

	}

	return &Song{Id: "0", Name: newSong.Name, DurationSec: newSong.DurationSec, Theme: newSong.Theme}, nil
}
