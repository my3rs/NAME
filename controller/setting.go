package controller

import (
	"NAME/database"
	"NAME/model"
	"github.com/kataras/iris/v12"
)

type SettingController struct {
	Ctx iris.Context
}

func (c *SettingController) GetBy(no string) {
	var setting model.Setting
	err := database.GetDB().Where("key = ?", no).First(&setting).Error
	if err != nil {
		c.Ctx.StatusCode(iris.StatusNotFound)
		c.Ctx.JSON(iris.Map{
			"success": false,
			"message": "没有找到对应的配置项",
		})
		return
	}

	c.Ctx.StatusCode(iris.StatusOK)
	c.Ctx.JSON(iris.Map{
		"success": true,
		"data":    setting,
	})

	return
}
