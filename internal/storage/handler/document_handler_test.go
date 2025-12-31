package handler

import (
	"context"
	"testing"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/internal/storage/manager/document"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewHandler_ReturnsHandlerWithDocumentManager(t *testing.T) {
	dm := &document.DocumentManager{}

	h, err := NewHandler(dm)
	require.NoError(t, err)
	require.NotNil(t, h)
	require.Equal(t, dm, h.documentManager)
}

func TestHandler_UploadFile_InvalidUserID_Panics(t *testing.T) {
	h := &Handler{}

	req := &storagepb.UploadFileRequest{
		UserId:   "not-a-uuid",
		FileName: "file.txt",
		FileSize: 10,
		Content:  []byte("data"),
	}

	require.Panics(t, func() {
		_, _ = h.UploadFile(context.Background(), req)
	})
}

func TestHandler_UploadFile_NilDocumentManager_Panics(t *testing.T) {
	// documentManager is nil; a valid request will cause a panic when attempting
	// to call UploadDocument on the nil manager.
	h := &Handler{}

	req := &storagepb.UploadFileRequest{
		UserId:   uuid.New().String(),
		FileName: "ok.txt",
		FileSize: 3,
		Content:  []byte("ok!"),
	}

	require.Panics(t, func() {
		_, _ = h.UploadFile(context.Background(), req)
	})
}
