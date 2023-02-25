package service

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/vmihailenco/msgpack/v5"
)

func extractBodyWithMsgpack(body io.ReadCloser, o interface{}) error {
	d := msgpack.NewDecoder(body)
	d.SetCustomStructTag("json")
	err := d.Decode(o)
	if err != nil {
		return err
	}
	return nil
}

// @Summary upload functions
// @Accept  json
// @Produce json
// @Success 200
// @Param   payload body object.FunctionUploadUnit true "Payload description"
// @Router  /api/v1/func [post]
// @Tags    Upload
func HandleFunctionUpload(c *gin.Context) {
	result := &object.FunctionUploadUnit{}
	if c.GetHeader("Content-Type") == object.BodyTypeMsgpack {
		err := extractBodyWithMsgpack(c.Request.Body, result)
		if err != nil {
			core.Log.Errorf("error when parse msgpack: %v\n", err)
			c.JSON(http.StatusBadRequest, fmt.Sprintf("parse msgpack error: %v", err))
			return
		}
	} else {
		err := c.BindJSON(result)
		if err != nil {
			core.Log.Errorf("error when parse json: %v\n", err)
			c.JSON(http.StatusBadRequest, fmt.Sprintf("parse json error: %v", err))
			return
		}
	}

	if err := result.WorkspaceConfig.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	go sharedQueue.SubmitFunc(result)
	c.JSON(http.StatusOK, "received")
}

// @Summary upload functions ctx
// @Accept  json
// @Produce json
// @Success 200
// @Param   payload body object.FunctionContextUploadUnit true "Payload description"
// @Router  /api/v1/funcctx [post]
// @Tags    Upload
func HandleFunctionContextUpload(c *gin.Context) {
	result := &object.FunctionContextUploadUnit{}
	if c.GetHeader("Content-Type") == object.BodyTypeMsgpack {
		err := extractBodyWithMsgpack(c.Request.Body, result)
		if err != nil {
			core.Log.Errorf("error when parse msgpack: %v\n", err)
			c.JSON(http.StatusBadRequest, fmt.Sprintf("parse msgpack error: %v", err))
			return
		}
	} else {
		err := c.BindJSON(result)
		if err != nil {
			core.Log.Errorf("error when parse json: %v\n", err)
			c.JSON(http.StatusBadRequest, fmt.Sprintf("parse json error: %v", err))
			return
		}
	}

	if err := result.WorkspaceConfig.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	go sharedQueue.SubmitFuncCtx(result)
	c.JSON(http.StatusOK, "received")
}

// @Summary upload class
// @Accept  json
// @Produce json
// @Success 200
// @Param   payload body object.ClazzUploadUnit true "Payload description"
// @Router  /api/v1/clazz [post]
// @Tags    Upload
func HandleClazzUpload(c *gin.Context) {
	result := &object.ClazzUploadUnit{}
	if c.GetHeader("Content-Type") == object.BodyTypeMsgpack {
		err := extractBodyWithMsgpack(c.Request.Body, result)
		if err != nil {
			core.Log.Errorf("error when parse msgpack: %v\n", err)
			c.JSON(http.StatusBadRequest, fmt.Sprintf("parse msgpack error: %v", err))
			return
		}
	} else {
		err := c.BindJSON(result)
		if err != nil {
			core.Log.Errorf("error when parse json: %v\n", err)
			c.JSON(http.StatusBadRequest, fmt.Sprintf("parse json error: %v", err))
			return
		}
	}

	if err := result.WorkspaceConfig.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	go sharedQueue.SubmitClazz(result)
	c.JSON(http.StatusOK, "received")
}
