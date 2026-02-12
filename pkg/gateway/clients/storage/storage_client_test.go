package storage

import (
	"context"
	"net"
	"testing"
	"time"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockStorageServiceClient struct {
	lastCtx       context.Context
	lastUploadReq *storagepb.UploadFileRequest
	lastListReq   *storagepb.ListFilesByUserIdRequest

	uploadResp *storagepb.UploadFileResponse
	listResp   *storagepb.ListFilesByUserIdResponse
	err        error
}

func (m *mockStorageServiceClient) UploadFile(ctx context.Context, in *storagepb.UploadFileRequest, opts ...grpc.CallOption) (*storagepb.UploadFileResponse, error) {
	m.lastCtx = ctx
	m.lastUploadReq = in
	return m.uploadResp, m.err
}

func (m *mockStorageServiceClient) ListFilesByUserId(ctx context.Context, in *storagepb.ListFilesByUserIdRequest, opts ...grpc.CallOption) (*storagepb.ListFilesByUserIdResponse, error) {
	m.lastCtx = ctx
	m.lastListReq = in
	return m.listResp, m.err
}

func TestStorageClientListFilesByUserIdUsesTimeoutAndForwardsRequest(t *testing.T) {
	mockClient := &mockStorageServiceClient{
		listResp: &storagepb.ListFilesByUserIdResponse{
			Documents: []*storagepb.Document{
				{
					Id:       "file-id-123",
					UserId:   "user-123",
					FileName: "test.txt",
					FileSize: 123,
				},
			},
		},
	}

	client := &storageClient{
		client: mockClient,
	}

	ctx := context.Background()
	req := &storagepb.ListFilesByUserIdRequest{
		UserId: "user-123",
	}

	resp, err := client.ListFilesByUserId(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, mockClient.listResp, resp)
	assert.Equal(t, req, mockClient.lastListReq)

	deadline, ok := mockClient.lastCtx.Deadline()
	assert.True(t, ok, "expected context to have a deadline")
	remaining := time.Until(deadline)
	assert.Greater(t, remaining, time.Duration(0))
	assert.LessOrEqual(t, remaining, 30*time.Second)
}

func TestStorageClientUploadFileUsesTimeoutAndForwardsRequest(t *testing.T) {
	mockClient := &mockStorageServiceClient{
		uploadResp: &storagepb.UploadFileResponse{
			Document: &storagepb.Document{
				Id:       "file-id-123",
				UserId:   "user-123",
				FileName: "test.txt",
				FileSize: 123,
			},
		},
	}

	client := &storageClient{
		client: mockClient,
	}

	ctx := context.Background()
	req := &storagepb.UploadFileRequest{
		UserId:   "user-123",
		FileName: "test.txt",
		FileSize: 123,
		Content:  []byte("hello"),
	}

	resp, err := client.UploadFile(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, mockClient.uploadResp, resp)

	assert.Equal(t, req, mockClient.lastUploadReq)

	deadline, ok := mockClient.lastCtx.Deadline()
	assert.True(t, ok, "expected context to have a deadline")
	remaining := time.Until(deadline)
	assert.Greater(t, remaining, time.Duration(0))
	assert.LessOrEqual(t, remaining, 30*time.Second)
}

type testStorageServer struct {
	storagepb.UnimplementedStorageServiceServer
}

func (s *testStorageServer) UploadFile(ctx context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	return &storagepb.UploadFileResponse{
		Document: &storagepb.Document{
			Id:       "generated-id",
			FileName: req.FileName,
		},
	}, nil
}

func TestNewStorageClientConnectsToServerAndUploads(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(grpcServer, &testStorageServer{})

	go func() {
		_ = grpcServer.Serve(lis)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = lis.Close()
	})

	client := NewStorageClient(lis.Addr().String())

	ctx := context.Background()
	req := &storagepb.UploadFileRequest{
		UserId:   "user-123",
		FileName: "uploaded.txt",
		FileSize: 10,
	}

	resp, err := client.UploadFile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "uploaded.txt", resp.Document.FileName)
	assert.NotEmpty(t, resp.Document.Id)
}
