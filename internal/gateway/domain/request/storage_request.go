package request

type UploadFileRequest struct {
	UserID string `form:"user_id" binding:"required,uuid4"`
}
