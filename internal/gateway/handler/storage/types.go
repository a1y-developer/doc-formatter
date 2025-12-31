package storage

import (
	"github.com/a1y/doc-formatter/internal/gateway/manager/storage"
)

type StorageHandler struct {
	storageManager *storage.StorageManager
}

func NewStorageHandler(storageManager *storage.StorageManager) (*StorageHandler, error) {
	return &StorageHandler{storageManager: storageManager}, nil
}
