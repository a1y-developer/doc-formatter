package storage

import (
	"io"
	"net/http"

	"github.com/a1y/doc-formatter/internal/gateway/domain/request"
	"github.com/gin-gonic/gin"
)

// UploadFile godoc
//
//	@Summary		Upload file
//	@Description	Upload a file for a user
//	@Tags			Storage
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user_id	formData	string	true	"User ID (UUID)"
//	@Param			file	formData	file	true	"File to upload"
//	@Success		201		{object}	response.UploadFileResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/v1/storage/upload [post]
func (h *StorageHandler) UploadFile(c *gin.Context) {
	var req request.UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.storageManager.UploadFile(c.Request.Context(), req.UserID, header.Filename, int64(len(fileBytes)), fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file_id":   resp.FileID,
		"file_name": resp.FileName,
	})
}
