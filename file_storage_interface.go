package storage

import (
	"io"
	"time"
)

type FileStorageInterface interface {
	GetBucketName() string
	CreateFolder(folderName string) error
	UploadFile(folderName string, fileName string, file io.Reader, size int64) error
	DownloadFile(folderName string, fileName string, saveFolder string) error
	GetFileLink(folderName string, filename string, expires time.Duration) (string, error)
	RemoveFile(folderName string, filename string) error
}
