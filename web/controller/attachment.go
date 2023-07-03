package controller

import (
	"NAME/service"
	"github.com/kataras/iris/v12"
)

type AttachmentController struct {
	Ctx               iris.Context
	AttachmentService service.AttachmentService
}

func NewAttachmentController() *AttachmentController {
	return &AttachmentController{AttachmentService: service.NewAttachmentService()}
}
