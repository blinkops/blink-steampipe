package blinkaws

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

func DownloadFileFromS3Bucket(objectKey string, destination string, region string, bucket string, accessKeyID string, secretKey string, token string) error {
	filename := filepath.Base(objectKey)
	path := filepath.Join(destination, filename)

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "download file from s3 bucket")
	}
	defer file.Close()

	awsSession, _ := session.NewSession(
		&aws.Config{Credentials: credentials.NewStaticCredentials(accessKeyID, secretKey, token),
			Region: aws.String(region)},
	)
	downloader := s3manager.NewDownloader(awsSession)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		})

	return errors.Wrap(err, "download file from s3 bucket")
}
