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
	GetFileLink(folderName string, fileName string, filePath string, expires time.Duration) (string, error)
	GetFile(bucketName string, fileName string) (io.Reader, error)
	RemoveFile(folderName string, fileName string) error
	RemoveFolder(bucketName string, folderName string) error
	GetFilesIntoFolder(bucketName string, folderName string) ([]string, error)
	CheckIfFileExists(bucketName string, filePath string) bool
}
