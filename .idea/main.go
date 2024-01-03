package main

import (
	"example/my-dependencies/record"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"net/http"
	//"strconv"
)

var DRIVER, ERR = neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "parola123", ""))

var records = []record.Record{}

func main() {
	router := gin.Default()
	router.GET("/records", getRecords)
	router.POST("/records", addRecord)
	router.GET("/records/:title", getRecordByTitle)

	if ERR != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", ERR)
	}
	defer DRIVER.Close()

	router.Run("localhost:8080")
}

// ENDPOINT IMPLEMENTATIONS

func getRecords(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, record.GetAllRecords(DRIVER))
}

func addRecord(c *gin.Context) {
	var newRecord record.Record

	// Call BindJSON to bind the received JSON to newRecord
	if err := c.BindJSON(&newRecord); err != nil {
		return
	}

	addedRecord, err := record.CreateRecord(DRIVER, newRecord)
	if err != nil {
		log.Fatalf("Failed to create record: %v", err)
	}
	log.Printf("Created record: %+v\n", addedRecord)

	// Add the new album to the slice.
	records = append(records, *addedRecord)
	c.IndentedJSON(http.StatusCreated, addedRecord)
}

func getRecordByTitle(c *gin.Context) {
	title := c.Param("title")

	for _, rec := range records {
		if rec.Title == title {
			c.IndentedJSON(http.StatusOK, rec)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "record not found"})
}
