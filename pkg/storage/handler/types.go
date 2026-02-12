package handler

import (
	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/pkg/storage/manager/document"
	"github.com/a1y/doc-formatter/pkg/storage/manager/folder"
)

func NewHandler(documentManager *document.DocumentManager) (*Handler, error) {
	return &Handler{documentManager: documentManager}, nil
}

type Handler struct {
	storagepb.UnimplementedStorageServiceServer
	documentManager *document.DocumentManager
	folderManager   *folder.FolderManager
}
