package controller

import (
	"NAME/model"
	"github.com/kataras/iris/v12"
)

type SettingController struct {
	Ctx iris.Context
}

func (c *SettingController) GetBy(no string) {
	setting, found := model.GetSettingsItem(no)
	if !found {
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
