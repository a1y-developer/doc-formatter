package storage

import (
	"github.com/a1y/doc-formatter/internal/gateway/clients/storage"
)

type StorageManager struct {
	client storage.StorageClient
}

func NewStorageManager(client storage.StorageClient) *StorageManager {
	return &StorageManager{client: client}
}
