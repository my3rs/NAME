package controller

import (
	"NAME/service"
	"github.com/kataras/iris/v12"
)

type MetaController struct {
	Ctx        iris.Context
	Service    service.ContentService
	TagService service.TagService
}

func (c *MetaController) Get() iris.Map {
	return c.Service.GetMeta()
}
