package controller

import "github.com/kataras/iris/v12"

type PageController struct {
	Ctx iris.Context
}

// Post handles POST : http://localchost/api/v1/pages
// @func create a new page
func (c *PageController) Post() {

}

// Get handles Get : http://localchost/api/v1/pages
// @func return all pages
func (c *PageController) Get() {

}

func (c *PageController) GetBy(id int64) {

}
