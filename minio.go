package storage

import (
	"github.com/minio/minio-go/v6"
	"io"
	"net/url"
	"time"
)

type MinioFileStorage struct {
	bucket string
	client *minio.Client
}

func NewMinioFileStorage(c *minio.Client, defaultBucket string) (*MinioFileStorage, error) {
	fileStorage := MinioFileStorage{client: c, bucket: defaultBucket}
	err := fileStorage.CreateFolder(defaultBucket)
	if err != nil {
		return nil, err
	}
	return &fileStorage, nil
}

func (m *MinioFileStorage) GetBucketName() string {
	return m.bucket
}

func (m *MinioFileStorage) CreateFolder(folderName string) error {
	err := m.client.MakeBucket(folderName, "")
	if err != nil {
		exists, errBucketExists := m.client.BucketExists(folderName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (m *MinioFileStorage) UploadFile(folderName string, fileName string, file io.Reader, size int64) error {
	_, err := m.client.PutObject(folderName,
		fileName,
		file,
		size,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioFileStorage) DownloadFile(folderName string, fileName string, saveFolder string) error {
	err := m.client.FGetObject(folderName, fileName, saveFolder, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioFileStorage) GetFileLink(folderName string, filename string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	url, err := m.client.PresignedGetObject(folderName,
		filename,
		expires,
		reqParams,
	)
	if err != nil {
		return "", err
	}
	return url.Host + url.Path + "?" + url.RawQuery, nil
}

func (m *MinioFileStorage) RemoveFile(folderName string, filename string) error {
	err := m.client.RemoveObject(folderName, filename)
	return err
}
