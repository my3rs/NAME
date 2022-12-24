package controller

import (
	"NAME/conf"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

type UploadController struct {
	Ctx      iris.Context
	Uploader service.Uploader
}

func NewUploadController() *UploadController {
	return &UploadController{Uploader: service.NewUploader()}
}

func (c *UploadController) Post() {
	c.Ctx.SetMaxRequestBodySize(conf.MaxBodySize)

	_, _, err := c.Ctx.UploadFormFiles("./uploads", func(ctx iris.Context, file *multipart.FileHeader) bool {
		today := time.Now().Format("2006-01")
		os.Mkdir("./uploads/"+today, 0700)

		file.Filename = strconv.FormatInt(time.Now().UnixMilli(), 10)
		return true
	})

	if err != nil {
		c.Ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
}
