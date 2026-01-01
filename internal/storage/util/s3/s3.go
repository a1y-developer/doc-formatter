package s3

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func (s *S3Storage) PutObject(ctx context.Context, objectKey string, file io.Reader) (bool, error) {
	_, err := s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			return false, errors.New("file is too large: " + objectKey + " in bucket: " + s.bucket)
		}
		return false, errors.New("failed to put object: " + objectKey + " in bucket: " + s.bucket + " with error: " + err.Error())
	} else {
		err = s3.NewObjectExistsWaiter(s.s3).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			return false, errors.New("failed to wait for object: " + objectKey + " in bucket: " + s.bucket + " with error: " + err.Error())
		}
	}
	return true, nil
}

func (s *S3Storage) GetObject(ctx context.Context, objectKey string) (io.Reader, error) {
	resp, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchKey" {
			return nil, errors.New("object not found: " + objectKey + " in bucket: " + s.bucket)
		}
		return nil, errors.New("failed to get object: " + objectKey + " in bucket: " + s.bucket + " with error: " + err.Error())
	} else {
		err = s3.NewObjectExistsWaiter(s.s3).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			return nil, errors.New("failed to wait for object: " + objectKey + " in bucket: " + s.bucket + " with error: " + err.Error())
		}
	}
	return resp.Body, nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, objectKey string) (bool, error) {
	_, err := s.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchKey" {
			return false, errors.New("object not found: " + objectKey + " in bucket: " + s.bucket)
		}
		return false, errors.New("failed to delete object: " + objectKey + " in bucket: " + s.bucket + " with error: " + err.Error())
	} else {
		err = s3.NewObjectExistsWaiter(s.s3).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			return false, errors.New("failed to wait for object: " + objectKey + " in bucket: " + s.bucket + " with error: " + err.Error())
		}
	}
	return true, nil
}
