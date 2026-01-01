package storagepb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUploadFileRequest_GettersAndString(t *testing.T) {
	t.Parallel()

	req := &UploadFileRequest{
		UserId:   "user-1",
		FileName: "file.txt",
		FileSize: 42,
		Content:  []byte("data"),
	}

	require.Equal(t, "user-1", req.GetUserId())
	require.Equal(t, "file.txt", req.GetFileName())
	require.EqualValues(t, 42, req.GetFileSize())
	require.Equal(t, []byte("data"), req.GetContent())
	require.NotEmpty(t, req.String())
	require.NotNil(t, req.ProtoReflect())
}

func TestUploadFileResponse_GettersAndString(t *testing.T) {
	t.Parallel()

	resp := &UploadFileResponse{
		FileId:   "id-1",
		FileName: "file.txt",
	}

	require.Equal(t, "id-1", resp.GetFileId())
	require.Equal(t, "file.txt", resp.GetFileName())
	require.NotEmpty(t, resp.String())
	require.NotNil(t, resp.ProtoReflect())
}
