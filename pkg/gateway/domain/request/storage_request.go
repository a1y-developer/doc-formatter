package request

import (
	"mime/multipart"

	"github.com/a1y/doc-formatter/pkg/gateway/domain/constant"
	"github.com/gin-gonic/gin"
)

type UploadFileRequest struct {
	UserID string                `form:"user_id" binding:"required,uuid4"`
	File   *multipart.FileHeader `form:"file"    binding:"required"`
}

func (r *UploadFileRequest) Validate() error {
	if r.UserID == "" {
		return constant.ErrEmptyUserID
	}
	return nil
}

func (r *UploadFileRequest) Decode(c *gin.Context) error {
	return c.ShouldBind(r)
}

type ListFilesByUserIdRequest struct {
	UserID string `form:"user_id" binding:"required,uuid4"`
}

func (r *ListFilesByUserIdRequest) Validate() error {
	if r.UserID == "" {
		return constant.ErrEmptyUserID
	}
	return nil
}

func (r *ListFilesByUserIdRequest) Decode(c *gin.Context) error {
	return c.ShouldBind(r)
}
