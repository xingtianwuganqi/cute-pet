package admin

import (
	"github.com/gin-gonic/gin"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
)

func CreatePetType(c *gin.Context) {
	var model models.PetActionType
	if err := c.ShouldBind(&model); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	var result = db.DB.Create(&model)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, map[string]interface{}{})

}

func UpdatePetType(c *gin.Context) {

}
