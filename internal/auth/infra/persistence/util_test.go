package persistence

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestAutoMigrate_UserModel_WithSQLite(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)
}

func TestMultiString_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		src       any
		want      MultiString
		wantError bool
	}{
		{
			name: "scan from []byte",
			src:  []byte("a,b,c"),
			want: MultiString{"a", "b", "c"},
		},
		{
			name: "scan from string",
			src:  "x,y",
			want: MultiString{"x", "y"},
		},
		{
			name: "scan from nil",
			src:  nil,
			want: nil,
		},
		{
			name:      "unsupported type",
			src:       123,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var m MultiString
			err := m.Scan(tt.src)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, m)
		})
	}
}

func TestMultiString_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  MultiString
		want driver.Value
	}{
		{
			name: "nil slice returns nil",
			val:  nil,
			want: nil,
		},
		{
			name: "non-nil slice joins with comma",
			val:  MultiString{"a", "b", "c"},
			want: "a,b,c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.val.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMultiString_GormDataType(t *testing.T) {
	t.Parallel()

	var m MultiString
	assert.Equal(t, "text", m.GormDataType())
}

func TestMultiString_GormDBDataType_SQLite(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	var m MultiString
	ft := &schema.Field{}

	_ = context.Background() // make sure context is imported if needed later

	got := m.GormDBDataType(db, ft)
	assert.Equal(t, "text", got)
}
