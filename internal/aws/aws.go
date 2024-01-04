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

func DownloadFileFromS3Bucket(objectKey, destination, region, bucket, accessKeyID, secretKey, endpoint, token string) error {
	filename := filepath.Base(objectKey)
	path := filepath.Join(destination, filename)

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "create file path '%s'", path)
	}

	defer func() {
		_ = file.Close()
	}()

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretKey, token),
		Region:      aws.String(region),
	}

	if endpoint != "" {
		config.Endpoint = aws.String(endpoint)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(sess)
	if _, err = downloader.Download(file, &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(objectKey)}); err != nil {
		return errors.Wrap(err, "download file from s3 bucket")
	}

	return nil
}
