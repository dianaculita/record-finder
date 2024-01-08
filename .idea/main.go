package main

import (
	"example/my-dependencies/artist"
	//"context"
	"example/my-dependencies/record"
	"fmt"

	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	//"strconv"
)

var DRIVER, ERR = neo4j.NewDriverWithContext("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "parola123", ""))

var records = []record.Record{}
var neo4jRecordRepo = record.NewNeo4jRecordRepository(DRIVER)
var neo4jArtistRepo = artist.NewNeo4jArtistRepository(DRIVER)

func main() {

	router := gin.Default()
	router.GET("/rf-app/records", getRecords)
	router.POST("/rf-app/records", addRecord)
	router.GET("/rf-app/records/search/:title", searchRecordsByTitle)

	router.GET("/rf-app/artists", getArtists)
	router.POST("/rf-app/artists", addArtist)
	router.GET("/rf-app/artists/search/:name", searchArtistsByName)

	if ERR != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", ERR)
	}
	//defer DRIVER.Close(router)

	router.Run("localhost:8080")
}

func searchArtistsByName(c *gin.Context) {
	name := c.Param("name")

	//todo: add support for pagination

	res, err := neo4jArtistRepo.SearchArtistsByName(c, name)

	if err != nil {
		log.Fatalf("Failed to search all artists by name: %v", err)
	}

	if res == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no artists found"})
	}

	c.IndentedJSON(http.StatusOK, res)
}

func addArtist(c *gin.Context) {
	var newArtist artist.Artist

	// Call BindJSON to bind the received JSON to newRecord
	if err := c.BindJSON(&newArtist); err != nil {
		return
	}

	fmt.Println(newArtist)
	err := neo4jArtistRepo.CreateArtist(c, newArtist)
	if err != nil {
		log.Fatalf("Failed to create artist: %v", err)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "artist successfully created"})
}

func getArtists(c *gin.Context) {
	res, err := neo4jArtistRepo.GetAll(c)

	if err != nil {
		log.Fatalf("Failed to get all artists: %v", err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

// ENDPOINT IMPLEMENTATIONS

func getRecords(c *gin.Context) {
	res, err := neo4jRecordRepo.GetAll(c)

	if err != nil {
		log.Fatalf("Failed to get all records: %v", err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func addRecord(c *gin.Context) {
	var newRecord record.Record

	// Call BindJSON to bind the received JSON to newRecord
	if err := c.BindJSON(&newRecord); err != nil {
		return
	}

	fmt.Println(newRecord)
	err := neo4jRecordRepo.CreateRecord(c, newRecord)
	if err != nil {
		log.Fatalf("Failed to create record: %v", err)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "record successfully created"})
}

func searchRecordsByTitle(c *gin.Context) {
	title := c.Param("title")

	//todo: add support for pagination

	res, err := neo4jRecordRepo.SearchRecordsByTitle(c, title)

	if err != nil {
		log.Fatalf("Failed to search all records by title: %v", err)
	}

	if res == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no record found"})
	}

	c.IndentedJSON(http.StatusOK, res)
}
