package core

import (
	"net/http"

	"github.com/Ranper/iam/pkg/log"
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	// Code 定义业务错误码
	Code int `json:"code"`

	// Message 包含细节. 并且是可以暴露给外界的
	Message string `json:"message"`

	// Reference 返回参考文档, 用于指引用户如何处理错误
	Reference string `json:"reference,omitempty"` // omitempty: 当字段值为该类型的“零值”或“空值”时，在序列化为JSON时忽略该字段
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)
		// coder := errors.ParserO
		// TODO: 错误码
		c.JSON(http.StatusOK, ErrResponse{
			Code:      1,
			Message:   err.Error(),
			Reference: "",
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
