package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Record struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Year   int32   `json:"year"`
	Price  float64 `json:"price"`
}

var records = []Record{
	{Id: "1", Title: "STMPD RCRDS", Artist: "Martin Garrix", Year: 2023, Price: 56.99},
	{Id: "2", Title: "Ether", Artist: "Martin Garrix", Year: 2019, Price: 90.99},
	{Id: "3", Title: "Anima", Artist: "Martin Garrix", Year: 2018, Price: 70.99},
}

func main() {
	router := gin.Default()
	router.GET("/records", getRecords)
	router.POST("/records", addRecord)
	router.GET("/records/:id", getRecordById)

	router.Run("localhost:8080")
}

func getRecords(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, records)
}

func addRecord(c *gin.Context) {
	var newRecord Record

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newRecord); err != nil {
		return
	}

	// Add the new album to the slice.
	records = append(records, newRecord)
	c.IndentedJSON(http.StatusCreated, newRecord)
}

func getRecordById(c *gin.Context) {
	id := c.Param("id")

	for _, rec := range records {
		if rec.Id == id {
			c.IndentedJSON(http.StatusOK, rec)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "record not found"})
}
