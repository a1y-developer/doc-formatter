package folder

import "github.com/a1y/doc-formatter/pkg/storage/domain/repository"

type FolderManager struct {
	folderRepo repository.FolderRepository
}

func NewFolderManager(
	folderRepo repository.FolderRepository,
) *FolderManager {
	return &FolderManager{
		folderRepo: folderRepo,
	}
}
