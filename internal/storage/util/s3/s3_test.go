package s3

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
)

func newTestS3Storage(t *testing.T, handler http.Handler) *S3Storage {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKID", "SECRET", "")),
	)
	require.NoError(t, err)

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Point client at our local test server and use path-style URLs.
		o.BaseEndpoint = aws.String(server.URL)
		o.UsePathStyle = true
	})

	return &S3Storage{
		s3:     client,
		bucket: "test-bucket",
	}
}

func TestS3Storage_PutObject_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Simulate successful PutObject
			w.WriteHeader(http.StatusOK)
		case http.MethodHead:
			// Waiter checks the object exists
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "unexpected method", http.StatusBadRequest)
		}
	})

	storage := newTestS3Storage(t, handler)

	ok, err := storage.PutObject(context.Background(), "path/to/object.txt", bytes.NewReader([]byte("hello")))
	require.NoError(t, err)
	require.True(t, ok)
}

func TestS3Storage_PutObject_Error(t *testing.T) {
	// Force PutObject to fail with a non-EntityTooLarge error and skip the waiter.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "unexpected method", http.StatusBadRequest)
	})

	storage := newTestS3Storage(t, handler)

	ok, err := storage.PutObject(context.Background(), "path/to/object.txt", bytes.NewReader([]byte("hello")))
	require.Error(t, err)
	require.False(t, ok)
}

func TestS3Storage_GetObject_Success(t *testing.T) {
	const body = "object-data"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Simulate successful GetObject
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, body)
		case http.MethodHead:
			// Waiter checks the object exists
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "unexpected method", http.StatusBadRequest)
		}
	})

	storage := newTestS3Storage(t, handler)

	reader, err := storage.GetObject(context.Background(), "path/to/object.txt")
	require.NoError(t, err)
	require.NotNil(t, reader)

	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, body, string(data))
}

func TestS3Storage_GetObject_NotFound(t *testing.T) {
	// Simulate S3 returning NoSuchKey on GET.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message></Error>`)
			return
		}
		http.Error(w, "unexpected method", http.StatusBadRequest)
	})

	storage := newTestS3Storage(t, handler)

	reader, err := storage.GetObject(context.Background(), "missing.txt")
	require.Error(t, err)
	require.Nil(t, reader)
}

func TestS3Storage_GetObject_WaiterError(t *testing.T) {
	const body = "object-data"

	// GET succeeds but HEAD (waiter) fails.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, body)
		case http.MethodHead:
			http.Error(w, "waiter error", http.StatusInternalServerError)
		default:
			http.Error(w, "unexpected method", http.StatusBadRequest)
		}
	})

	storage := newTestS3Storage(t, handler)

	reader, err := storage.GetObject(context.Background(), "path/to/object.txt")
	require.Error(t, err)
	require.Nil(t, reader)
}

func TestS3Storage_DeleteObject_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			// Simulate successful DeleteObject
			w.WriteHeader(http.StatusOK)
		case http.MethodHead:
			// Waiter checks the object state
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "unexpected method", http.StatusBadRequest)
		}
	})

	storage := newTestS3Storage(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ok, err := storage.DeleteObject(ctx, "path/to/object.txt")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestS3Storage_DeleteObject_Error(t *testing.T) {
	// Force DeleteObject to fail with a non-NoSuchKey error and skip the waiter.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "unexpected method", http.StatusBadRequest)
	})

	storage := newTestS3Storage(t, handler)

	ok, err := storage.DeleteObject(context.Background(), "path/to/object.txt")
	require.Error(t, err)
	require.False(t, ok)
}

func TestS3Storage_DeleteObject_WaiterError(t *testing.T) {
	// DELETE succeeds but HEAD (waiter) fails.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		case http.MethodHead:
			http.Error(w, "waiter error", http.StatusInternalServerError)
		default:
			http.Error(w, "unexpected method", http.StatusBadRequest)
		}
	})

	storage := newTestS3Storage(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ok, err := storage.DeleteObject(ctx, "path/to/object.txt")
	require.Error(t, err)
	require.False(t, ok)
}
