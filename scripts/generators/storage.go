package generators

import (
	"os"
	"os/exec"
	"path/filepath"

	blinkaws "github.com/blinkops/blink-steampipe/internal/aws"
	"github.com/pkg/errors"
)

const (
	SteampipeStorageFileIdentifier = "STORAGE_FILE_IDENTIFIER"
	SteampipeStorageAWSKeyID       = "STORAGE_AWS_KEY_ID"
	SteampipeStorageAWSSecretKey   = "STORAGE_AWS_SECRET_KEY"
	SteampipeStorageAWSBucket      = "STORAGE_AWS_BUCKET"
	SteampipeStorageAWSRootDir     = "STORAGE_AWS_ROOT_DIR"
	SteampipeStorageAWSToken       = "STORAGE_AWS_TOKEN"
	SteampipeStorageAWSRegion      = "STORAGE_AWS_REGION"
	storageDestination             = "/home/steampipe/storage"
	tempDesitnation                = "/home/steampipe/tmp"
)

type StorageCredentialsGenerator struct{}

func (gen StorageCredentialsGenerator) Generate() error {
	fileID, ok := os.LookupEnv(SteampipeStorageFileIdentifier)
	if !ok {
		return nil
	}
	awsAccessKeyID, ok := os.LookupEnv(SteampipeStorageAWSKeyID)
	if !ok {
		return errors.New("missing STORAGE_AWS_KEY_ID env. key")
	}
	awsSecretKey, ok := os.LookupEnv(SteampipeStorageAWSSecretKey)
	if !ok {
		return errors.New("missing STORAGE_AWS_SECRET_KEY env. key")
	}
	bucket, ok := os.LookupEnv(SteampipeStorageAWSBucket)
	if !ok {
		return errors.New("missing SteampipeStorageAWSBucket env. key")
	}
	region, ok := os.LookupEnv(SteampipeStorageAWSRegion)
	if !ok {
		return errors.New("missing SteampipeStorageAWSBucket env. key")
	}
	rootDir, ok := os.LookupEnv(SteampipeStorageAWSRootDir)
	if !ok {
		return errors.New("missing STORAGE_AWS_ROOT_DIR env. key")
	}
	token, ok := os.LookupEnv(SteampipeStorageAWSToken)
	if !ok {
		return errors.New("missing STORAGE_AWS_ROOT_DIR env. key")
	}

	if err := blinkaws.DownloadFileFromS3Bucket(filepath.Join(rootDir, fileID),
		tempDesitnation,
		region,
		bucket,
		awsAccessKeyID,
		awsSecretKey,
		token); err != nil {
		return errors.New("download content for steampipe")
	}

	sourcePath := filepath.Join(tempDesitnation, fileID)

	cmd := exec.Command("tar", "-xzf", sourcePath, "-C", storageDestination)
	if _, err := cmd.Output(); err != nil {
		return errors.Wrap(err, "extract user content")
	}

	return nil
}
