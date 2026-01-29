package response

type Document struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
