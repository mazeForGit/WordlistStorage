package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func HealthGET(c *gin.Context) {
	fmt.Println("HealthGET")

	c.JSON(200, gin.H{
		"status": "UP",
	})
}
func HealthPOST(c *gin.Context) {
	fmt.Println("HealthPOST")
	
	c.JSON(200, gin.H{
		"status": "UP post",
	})
}
