package response

type UploadFileResponse struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
}
