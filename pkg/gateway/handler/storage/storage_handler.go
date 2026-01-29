package storage

import (
	"io"
	"net/http"

	"github.com/a1y/doc-formatter/pkg/gateway/domain/request"
	"github.com/a1y/doc-formatter/pkg/gateway/handler"
	"github.com/gin-gonic/gin"
)

//	@Id				UploadFile
//	@Summary		Upload file
//	@Description	Upload a file for a user
//	@Tags			Storage
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user_id	formData	string										true	"User ID (UUID)"
//	@Param			file	formData	file										true	"File to upload"
//	@Success		201		{object}	handler.Response{data=response.Document}	"Success"
//	@Failure		400		{object}	error										"Bad Request"
//	@Failure		401		{object}	error										"Unauthorized"
//	@Failure		429		{object}	error										"Too Many Requests"
//	@Failure		404		{object}	error										"Not Found"
//	@Failure		500		{object}	error										"Internal Server Error"
//	@Router			/api/v1/storage/upload [post]
func (h *StorageHandler) UploadFile(c *gin.Context) {
	var requestPayload request.UploadFileRequest
	if err := requestPayload.Decode(c); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	// Validate request payload
	if err := requestPayload.Validate(); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	file, err := requestPayload.File.Open()
	if err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		handler.HandleResult(c, err, nil)
		return
	}

	response, err := h.storageManager.UploadFile(c.Request.Context(), requestPayload.UserID, requestPayload.File.Filename, int64(len(fileBytes)), fileBytes)
	if err != nil {
		handler.HandleResult(c, err, response)
		return
	}
	handler.HandleResultWithStatus(c, nil, response, http.StatusCreated)
}

//	@Id				ListFilesByUserId
//	@Summary		List files by user ID
//	@Description	Retrieve a list of files uploaded by a specific user
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Param			user_id	query		string										true	"User ID (UUID)"
//	@Success		200		{object}	handler.Response{data=[]response.Document}	"Success"
//	@Failure		400		{object}	error										"Bad Request"
//	@Failure		401		{object}	error										"Unauthorized"
//	@Failure		429		{object}	error										"Too Many Requests"
//	@Failure		404		{object}	error										"Not Found"
//	@Failure		500		{object}	error										"Internal Server Error"
//	@Router			/api/v1/storage/files [get]
func (h *StorageHandler) ListFilesByUserId(c *gin.Context) {
	var requestPayload request.ListFilesByUserIdRequest
	if err := requestPayload.Decode(c); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	// Validate request payload
	if err := requestPayload.Validate(); err != nil {
		handler.HandleResultWithStatus(c, err, nil, http.StatusBadRequest)
		return
	}

	response, err := h.storageManager.ListFilesByUserId(c.Request.Context(), requestPayload.UserID)
	if err != nil {
		handler.HandleResult(c, err, response)
		return
	}
	handler.HandleResult(c, nil, response)
}
