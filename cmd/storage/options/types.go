package options

import (
	"os"
)

const (
	DefaultDBPort = 5432
	DefaultPort   = 8082
)

var (
	DBHostEnv      = os.Getenv("STORAGE_DB_HOST")
	DBPortEnv      = os.Getenv("STORAGE_DB_PORT")
	DBUserEnv      = os.Getenv("STORAGE_DB_USER")
	DBPassEnv      = os.Getenv("STORAGE_DB_PASS")
	DBNameEnv      = os.Getenv("STORAGE_DB_NAME")
	PortEnv        = os.Getenv("STORAGE_PORT")
	AutoMigrateEnv = os.Getenv("STORAGE_AUTO_MIGRATE")
	S3EndpointEnv  = os.Getenv("STORAGE_S3_ENDPOINT")
	S3RegionEnv    = os.Getenv("STORAGE_S3_REGION")
	S3AccessIDEnv  = os.Getenv("STORAGE_S3_ACCESS_KEY_ID")
	S3AccessKeyEnv = os.Getenv("STORAGE_S3_ACCESS_KEY_SECRET")
	S3BucketEnv    = os.Getenv("STORAGE_S3_BUCKET")
	S3ForcePathEnv = os.Getenv("STORAGE_S3_FORCE_PATH_STYLE")
)
