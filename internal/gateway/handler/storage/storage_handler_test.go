package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	clientstorage "github.com/a1y/doc-formatter/internal/gateway/clients/storage"
	storagemgr "github.com/a1y/doc-formatter/internal/gateway/manager/storage"
	"github.com/a1y/doc-formatter/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockStorageClient struct {
	clientstorage.StorageClient

	resp *storagepb.UploadFileResponse
	err  error

	lastReq *storagepb.UploadFileRequest
}

func (m *mockStorageClient) UploadFile(_ context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	m.lastReq = req
	return m.resp, m.err
}

func newTestHandler(t *testing.T, mockClient *mockStorageClient) *StorageHandler {
	t.Helper()

	manager := storagemgr.NewStorageManager(mockClient)
	h, err := NewStorageHandler(manager)
	assert.NoError(t, err)
	return h
}

func setupRouter(h *StorageHandler) *gin.Engine {
	r := testutil.NewGinEngine()
	r.POST("/api/v1/storage/upload", h.UploadFile)
	return r
}

func createMultipartRequest(t *testing.T, userID string, includeFile bool) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if userID != "" {
		err := writer.WriteField("user_id", userID)
		assert.NoError(t, err)
	}

	if includeFile {
		fileWriter, err := writer.CreateFormFile("file", "test.txt")
		assert.NoError(t, err)
		content := []byte("hello world")
		_, err = io.Copy(fileWriter, bytes.NewReader(content))
		assert.NoError(t, err)
	}

	assert.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func TestStorageHandler_UploadFileSuccess(t *testing.T) {
	mockClient := &mockStorageClient{
		resp: &storagepb.UploadFileResponse{
			FileId:   "file-id-123",
			FileName: "test.txt",
		},
	}

	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	req := createMultipartRequest(t, userID, true)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var respBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "file-id-123", respBody["file_id"])
	assert.Equal(t, "test.txt", respBody["file_name"])

	if assert.NotNil(t, mockClient.lastReq) {
		assert.Equal(t, userID, mockClient.lastReq.GetUserId())
		assert.Equal(t, "test.txt", mockClient.lastReq.GetFileName())
		assert.Equal(t, int64(len("hello world")), mockClient.lastReq.GetFileSize())
		assert.Equal(t, []byte("hello world"), mockClient.lastReq.GetContent())
	}
}

func TestStorageHandler_UploadFileBindError(t *testing.T) {
	mockClient := &mockStorageClient{}
	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)
	_, err = io.Copy(fileWriter, bytes.NewReader([]byte("hello world")))
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_UploadFileFileMissing(t *testing.T) {
	mockClient := &mockStorageClient{}
	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	req := createMultipartRequest(t, "550e8400-e29b-41d4-a716-446655440000", false)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_UploadFileManagerError(t *testing.T) {
	mockClient := &mockStorageClient{
		err: errors.New("upload failed"),
	}

	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	req := createMultipartRequest(t, "550e8400-e29b-41d4-a716-446655440000", true)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
