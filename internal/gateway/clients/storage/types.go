package storage

import (
	"context"
	"log"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StorageClient interface {
	UploadFile(ctx context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error)
}

var _ StorageClient = &storageClient{}

type storageClient struct {
	conn   *grpc.ClientConn
	client storagepb.StorageServiceClient
}

func NewStorageClient(addr string) StorageClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect to StorageService: %v", err)
		return nil
	}
	c := storagepb.NewStorageServiceClient(conn)
	return &storageClient{conn: conn, client: c}
}
