package awscli

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.uber.org/zap"
)

type S3StorageConfig struct {
	RestEndPoint           string
	AccessKey              string
	SecretAccessKey        string
	BucketName             string
	Region                 string
	PathStyleAccessEnabled bool
}

type s3Storage struct {
	config   S3StorageConfig
	s3Client *s3.S3
	log      *zap.Logger
}

func NewS3StorageClient(
	config S3StorageConfig,
	log *zap.Logger,
) repository.StorageRepository {

	awsSession, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(config.RestEndPoint),
		Region:           aws.String(config.Region),
		Credentials:      credentials.NewStaticCredentials(config.AccessKey, config.SecretAccessKey, ""),
		S3ForcePathStyle: &config.PathStyleAccessEnabled,
	})

	if err != nil {
		panic(err)
	}

	s3Client := s3.New(awsSession)

	return &s3Storage{
		config,
		s3Client,
		log,
	}
}

func (port *s3Storage) GetStorageType() string {
	return "S3"
}

func (port *s3Storage) DeleteFile(fmd domain.FileMetaData) error {
	req, _ := port.s3Client.DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: aws.String(port.config.BucketName),
		Key:    aws.String(fmd.Path),
	})
	err := req.Send()
	if err != nil {
		return fmt.Errorf("can't delete file: %w", err)
	}
	return nil
}

func (port *s3Storage) GetFileDownloadUrl(fmd domain.FileMetaData) (string, error) {
	url := port.config.RestEndPoint + "/" + port.config.BucketName + "/" + fmd.Path
	return url, nil
}

func (port *s3Storage) GetFileUploadUrl(fmd domain.FileMetaData) (string, error) {
	req, _ := port.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(port.config.BucketName),
		Key:    aws.String(fmd.Path),
	})
	url, err := req.Presign(10 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("can't preSign url for upload: %w", err)
	}
	return url, nil
}

// privates
