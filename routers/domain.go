package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/model"
)
func DomainGET(c *gin.Context) {	
	c.JSON(200, model.GetDomains())
}