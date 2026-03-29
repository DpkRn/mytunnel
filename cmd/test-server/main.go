package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Server starting : 8080")
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "welcome to the server"})
	})
	router.GET("/request1", func(c *gin.Context) {
		time.Sleep(10 * time.Second)
		c.JSON(200, gin.H{"message": "served request1"})
	})
	router.GET("/request2", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "served request2"})
	})

	router.POST("/fullname", func(c *gin.Context) {

		payload := struct {
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
		}{}
		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(500, gin.H{"message": "error binding json"})
			return
		}
		fmt.Println("payload:", payload)
		fullname := payload.Firstname + " " + payload.Lastname
		fmt.Println("fullname:", fullname)
		c.JSON(http.StatusOK, gin.H{"message": fullname})

	})
	router.Run(":8080")
}
