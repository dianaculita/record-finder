package main

import (
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

func main() {

	router := gin.Default()
	router.GET("/records", getRecords)
	router.POST("/records", addRecord)
	router.GET("/records/search/:title", searchRecordsByTitle)

	if ERR != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", ERR)
	}
	//defer DRIVER.Close(router)

	router.Run("localhost:8080")
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
