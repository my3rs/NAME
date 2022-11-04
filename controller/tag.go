package controller

import (
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
)

type TagController struct {
	Ctx     iris.Context
	Service service.TagService
}

func NewTagController() *TagController {
	return &TagController{Service: service.NewTagService()}
}

// Get handles GET /api/v1/tags
// @Summary Get tags list
// @Description Get tags list
// @Accept  json
// @Produce  json
// @Param 	Authorization	header	string	true	"Access token with the prefix `Bearer `"
// @Param   pageSize    query  	int      true        "page size"	 default(1)
// @Param 	pageIndex	query	int 		true		"page index"	default(1)
// @Param	orderBy		query	string  	false 		"order"		default("created_at desc")
// @Success 200 		{object} model.Response	"success"
// @Failure 400 		{object} model.Response "Bad request"
// @Router /api/v1/posts [get]
func (c *TagController) Get(request model.QueryRequest) iris.Map {
	tags, err := c.Service.GetAllTags()
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return iris.Map{"success": true}
	}

	c.Ctx.StatusCode(iris.StatusOK)
	return iris.Map{"success": true, "data": tags}
}
