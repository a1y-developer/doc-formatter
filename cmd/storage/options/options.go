package options

import (
	"context"
	"net"
	"strconv"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/internal/storage"
	"github.com/a1y/doc-formatter/internal/storage/handler"
	storagepersistence "github.com/a1y/doc-formatter/internal/storage/infra/persistence"
	"github.com/a1y/doc-formatter/internal/storage/manager/document"
	storages3 "github.com/a1y/doc-formatter/internal/storage/util/s3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"k8s.io/kubectl/pkg/util/i18n"
)

type StorageOptions struct {
	Port int

	Database DatabaseOptions

	S3Endpoint        string
	S3Region          string
	S3AccessKeyID     string
	S3AccessKeySecret string
	S3Bucket          string
	S3ForcePathStyle  bool
}

func NewStorageOptions() *StorageOptions {
	return &StorageOptions{
		Port:     DefaultPort,
		Database: DatabaseOptions{},
	}
}

func (o *StorageOptions) Complete(args []string) {}

func (o *StorageOptions) Validate() error {
	return o.Database.Validate()
}

func (o *StorageOptions) Config() (*storage.Config, error) {
	cfg := storage.NewConfig()
	if err := o.Database.ApplyTo(&cfg.DB); err != nil {
		return nil, err
	}

	cfg.Port = o.Port
	cfg.EndPoint = o.S3Endpoint
	cfg.Region = o.S3Region
	cfg.AccessKeyID = o.S3AccessKeyID
	cfg.AccessKeySecret = o.S3AccessKeySecret
	cfg.Bucket = o.S3Bucket
	cfg.ForcePathStyle = o.S3ForcePathStyle

	return cfg, nil
}

func (o *StorageOptions) AddFlags(cmd *cobra.Command) {
	port, err := strconv.Atoi(PortEnv)
	if err != nil {
		port = DefaultPort
	}
	cmd.Flags().IntVarP(&o.Port, "port", "p", port,
		i18n.T("specify the port for the storage service to listen on"))

	cmd.Flags().StringVar(&o.S3Endpoint, "s3-endpoint", S3EndpointEnv,
		i18n.T("specify the S3 endpoint for the storage service"))
	cmd.Flags().StringVar(&o.S3Region, "s3-region", S3RegionEnv,
		i18n.T("specify the S3 region for the storage service"))
	cmd.Flags().StringVar(&o.S3AccessKeyID, "s3-access-key-id", S3AccessIDEnv,
		i18n.T("specify the S3 access key ID"))
	cmd.Flags().StringVar(&o.S3AccessKeySecret, "s3-access-key-secret", S3AccessKeyEnv,
		i18n.T("specify the S3 access key secret"))
	cmd.Flags().StringVar(&o.S3Bucket, "s3-bucket", S3BucketEnv,
		i18n.T("specify the S3 bucket name"))
	cmd.Flags().BoolVar(&o.S3ForcePathStyle, "s3-force-path-style", false,
		i18n.T("whether to enable path-style access for S3"))

	o.Database.AddFlags(cmd.Flags())
}

func (o *StorageOptions) Run() error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	documentRepository := storagepersistence.NewDocumentRepository(config.DB)

	ctx := context.Background()
	s3Storage, err := storages3.NewS3Storage(ctx, config)
	if err != nil {
		return err
	}

	documentManager := document.NewDocumentManager(documentRepository, s3Storage)
	storageHandler, err := handler.NewHandler(documentManager)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(server, storageHandler)

	logrus.Infof("Storage service running at :%d", config.Port)

	if err = server.Serve(lis); err != nil {
		return err
	}

	return nil
}
