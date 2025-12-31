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
	lastCtx context.Context
	lastReq *storagepb.UploadFileRequest

	resp *storagepb.UploadFileResponse
	err  error
}

func (m *mockStorageServiceClient) UploadFile(ctx context.Context, in *storagepb.UploadFileRequest, opts ...grpc.CallOption) (*storagepb.UploadFileResponse, error) {
	m.lastCtx = ctx
	m.lastReq = in
	return m.resp, m.err
}

func TestStorageClientUploadFileUsesTimeoutAndForwardsRequest(t *testing.T) {
	mockClient := &mockStorageServiceClient{
		resp: &storagepb.UploadFileResponse{
			FileId:   "file-id",
			FileName: "test.txt",
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
	assert.Equal(t, mockClient.resp, resp)

	assert.Equal(t, req, mockClient.lastReq)

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
		FileId:   "generated-id",
		FileName: req.FileName,
	}, nil
}

func TestNewStorageClientConnectsToServerAndUploads(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(grpcServer, &testStorageServer{})

	go grpcServer.Serve(lis)
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
	assert.Equal(t, "uploaded.txt", resp.FileName)
	assert.NotEmpty(t, resp.FileId)
}
