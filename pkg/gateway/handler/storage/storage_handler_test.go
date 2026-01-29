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
	clientstorage "github.com/a1y/doc-formatter/pkg/gateway/clients/storage"
	storagemgr "github.com/a1y/doc-formatter/pkg/gateway/manager/storage"
	"github.com/a1y/doc-formatter/pkg/util/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockStorageClient struct {
	clientstorage.StorageClient

	resp     *storagepb.UploadFileResponse
	listResp *storagepb.ListFilesByUserIdResponse
	err      error

	lastReq *storagepb.UploadFileRequest
}

func (m *mockStorageClient) UploadFile(_ context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	m.lastReq = req
	return m.resp, m.err
}

func (m *mockStorageClient) ListFilesByUserId(_ context.Context, req *storagepb.ListFilesByUserIdRequest) (*storagepb.ListFilesByUserIdResponse, error) {
	return m.listResp, m.err
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
	r.GET("/api/v1/storage/files", h.ListFilesByUserId)
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
			Document: &storagepb.Document{
				Id:       "file-id-123",
				UserId:   "550e8400-e29b-41d4-a716-446655440000",
				FileName: "test.txt",
				FileSize: int64(len("hello world")),
			},
		},
	}

	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	req := createMultipartRequest(t, userID, true)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var respBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)

	data, ok := respBody["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "file-id-123", data["id"])
	assert.Equal(t, "test.txt", data["file_name"])

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

func TestStorageHandler_ListFilesByUserIdSuccess(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	mockClient := &mockStorageClient{
		listResp: &storagepb.ListFilesByUserIdResponse{
			Documents: []*storagepb.Document{
				{
					Id:       "doc-1",
					UserId:   userID,
					FileName: "file1.pdf",
					FileSize: 1024,
				},
				{
					Id:       "doc-2",
					UserId:   userID,
					FileName: "file2.txt",
					FileSize: 512,
				},
			},
		},
	}

	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/storage/files?user_id="+userID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var respBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	assert.NoError(t, err)

	data, ok := respBody["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)
}

func TestStorageHandler_ListFilesByUserIdInvalidUserID(t *testing.T) {
	mockClient := &mockStorageClient{}
	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/storage/files?user_id=invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_ListFilesByUserIdMissingUserID(t *testing.T) {
	mockClient := &mockStorageClient{}
	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/storage/files", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_ListFilesByUserIdManagerError(t *testing.T) {
	mockClient := &mockStorageClient{
		err: errors.New("list failed"),
	}

	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/storage/files?user_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestStorageHandler_UploadFileValidationError(t *testing.T) {
	mockClient := &mockStorageClient{}
	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	// Create multipart request without user_id (should fail validation)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file but no user_id
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

func TestStorageHandler_ListFilesByUserIdValidationError(t *testing.T) {
	mockClient := &mockStorageClient{}
	h := newTestHandler(t, mockClient)
	router := setupRouter(h)

	// Request without user_id should fail validation
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/storage/files?user_id=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
