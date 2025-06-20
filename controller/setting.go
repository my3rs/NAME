package controller

import (
	"NAME/database"
	"NAME/model"
	"github.com/kataras/iris/v12"
)

type SettingController struct {
	Ctx iris.Context
}

func (c *SettingController) GetBy(no string) model.DetailResponse {
	var setting model.Setting
	err := database.GetDB().Where("key = ?", no).First(&setting).Error
	if err != nil {
		c.Ctx.StatusCode(iris.StatusNotFound)
		return model.DetailResponse{
			Success: false,
			Message: "没有找到对应的配置项",
		}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return model.DetailResponse{
		Success: true,
		Message: "success",
		Data:    setting,
	}
}
