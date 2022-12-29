package controller

import (
	"NAME/conf"
	"NAME/model"
	"NAME/service"
	"github.com/kataras/iris/v12"
	"os"
	"path"
	"strconv"
	"time"
)

type AttachmentController struct {
	Ctx               iris.Context
	AttachmentService service.AttachmentService
}

func NewAttachmentController() *AttachmentController {
	return &AttachmentController{AttachmentService: service.NewAttachmentService()}
}

func (c *AttachmentController) Post() {
	c.Ctx.SetMaxRequestBodySize(conf.MaxBodySize)

	// Parse request
	maxSize := c.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

	err := c.Ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		c.Ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	form := c.Ctx.Request().MultipartForm
	files := form.File["files"]
	//c.Ctx.FormFiles("files")
	c.Ctx.Application().Logger().Debugf("form.Value %+v", form.Value)
	c.Ctx.Application().Logger().Debugf("form.Value[\"contentID\"] %+v", form.Value["contentID"])
	c.Ctx.Application().Logger().Debugf("form.Value[\"files\"] %+v", form.Value["files"])
	c.Ctx.Application().Logger().Debug("form.File: ", form.File)
	c.Ctx.Application().Logger().Debugf("Receive %d files", len(files))
	if len(files) == 0 {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		c.Ctx.JSON(iris.Map{
			"Success": false,
			"Message": "received 0 file",
		})
		return
	}

	// Check whether these files belong to a content
	var contentID int
	id := c.Ctx.FormValue("contentID")
	if len(id) != 0 {
		c.Ctx.Application().Logger().Info("content ID ", id)
		contentID, err = strconv.Atoi(id)
	} else {
		c.Ctx.Application().Logger().Info("no content ID specified")
		contentID = 0
	}

	// Create folder to store files
	today := time.Now().Format("2006-01")
	destDirectory := conf.Config().DataPath + "/uploads/" + today + "/"
	err = os.Mkdir(destDirectory, 0700)
	if err != nil {
		c.Ctx.Application().Logger().Error("Failed to create uploads folder: ", err)
	}

	// Store files
	failures := 0
	message := ""
	for _, file := range files {
		oldFileName := file.Filename
		suffix := path.Ext(file.Filename)
		file.Filename = strconv.FormatInt(time.Now().UnixMilli(), 10) + suffix

		// Insert `Attachment` to database
		attachment := model.Attachment{
			ContentID: uint(contentID),
			Name:      file.Filename,
			Path:      "/uploads/" + today,
		}
		c.Ctx.Application().Logger().Debugf("Inserting attachment: %+v", attachment)

		if e := c.AttachmentService.InsertAttachment(attachment); e != nil {
			message += " failed to insert attachment: " + oldFileName
		}

		_, err := c.Ctx.SaveFormFile(file, destDirectory+file.Filename)
		if err != nil {
			failures++
			message += " failed to upload: " + oldFileName
		}
	}

	c.Ctx.JSON(iris.Map{
		"Success": failures == 0,
		"Message": message,
	})

}
