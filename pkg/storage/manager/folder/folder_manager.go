package folder

import (
	"context"

	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/google/uuid"
)

func (f *FolderManager) CreateFolder(ctx context.Context, folder *entity.Folder) (*entity.Folder, error) {
	if folder.ParentID == uuid.Nil {
		folder.Name = folder.OwnerID.String()
		folder.PathTokens = []string{folder.OwnerID.String()}
	}
	createdFolder, err := f.folderRepo.Create(ctx, folder)
	if err != nil {
		return nil, err
	}
	return createdFolder, nil
}

func (f *FolderManager) GetFolder(ctx context.Context, id uuid.UUID) (*entity.Folder, error) {
	folder, err := f.folderRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return folder, nil
}
