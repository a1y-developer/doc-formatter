package s3

import (
	"context"
	"errors"

	"github.com/a1y/doc-formatter/internal/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	s3     *s3.Client
	bucket string
}

func NewS3Storage(ctx context.Context, config *storage.Config) (*S3Storage, error) {
	cred := credentials.NewStaticCredentialsProvider(
		config.AccessKeyID,
		config.AccessKeySecret,
		"",
	)
	awsConfig, err := awsconfig.LoadDefaultConfig(
		ctx, awsconfig.WithRegion(config.Region),
		awsconfig.WithCredentialsProvider(cred),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsConfig)
	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(config.Bucket),
	})
	if err != nil {
		return nil, errors.New("bucket does not exist: " + config.Bucket)
	}

	return &S3Storage{
		s3:     s3Client,
		bucket: config.Bucket,
	}, nil
}
