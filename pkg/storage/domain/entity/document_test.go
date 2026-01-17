package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDocument_Validate_Success(t *testing.T) {
	t.Parallel()

	d := &Document{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		FileName:  "file.txt",
		FileSize:  123,
		ObjectKey: "user/file.txt",
	}

	require.NoError(t, d.Validate())
}

func TestDocument_Validate_MissingUserID(t *testing.T) {
	t.Parallel()

	d := &Document{
		FileName:  "file.txt",
		ObjectKey: "file.txt",
	}

	err := d.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "user id is required")
}

func TestDocument_Validate_MissingFileName(t *testing.T) {
	t.Parallel()

	d := &Document{
		UserID:    uuid.New(),
		ObjectKey: "file.txt",
	}

	err := d.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "file name is required")
}

func TestDocument_Validate_MissingObjectKey(t *testing.T) {
	t.Parallel()

	d := &Document{
		UserID:   uuid.New(),
		FileName: "file.txt",
	}

	err := d.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "object key is required")
}
