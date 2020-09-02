package storage

import (
	"fmt"
	"github.com/minio/minio-go/v6"
	"io"
	"net/url"
	"time"
	"mime"
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

func (m *MinioFileStorage) GetFileLink(folderName string, fileName string, filePath string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-Content-Disposition", "attachment; filename=\""+fileName+"\"")
	mime, err := mimetype.DetectFile(fileName)
	if err != nil {
		return "", err
	}
	reqParams.Set("response-Content-Type", mime)
	url, err := m.client.PresignedGetObject(folderName,
		filePath,
		expires,
		reqParams,
	)
	if err != nil {
		return "", err
	}
	return url.Path + "?" + url.RawQuery, nil
}

func (m *MinioFileStorage) RemoveFile(folderName string, fileName string) error {
	err := m.client.RemoveObject(folderName, fileName)
	return err
}

func (m *MinioFileStorage) RemoveFolder(bucketName string, folderName string) error {
	objectsCh := make(chan string)
	go func() {
		defer close(objectsCh)
		objectCh := m.client.ListObjects(bucketName, folderName, true, nil)
		for object := range objectCh {
			if object.Err != nil {
				fmt.Println(object.Err)
				return
			}
			objectsCh <- object.Key
		}
	}()
	for rErr := range m.client.RemoveObjects(bucketName, objectsCh) {
		return rErr.Err
	}
	return nil
}

func (m *MinioFileStorage) GetFilesIntoFolder(bucketName string, folderName string) ([]string, error) {
	var filePaths []string
	objectCh := m.client.ListObjects(bucketName, folderName, true, nil)
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		filePaths = append(filePaths, object.Key)
	}
	return filePaths, nil
}
