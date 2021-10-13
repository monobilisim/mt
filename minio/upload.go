package minio

import (
	"context"
	"errors"
	"github.com/fatih/structs"
	"github.com/minio/minio-go/v7"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type UploadParams struct {
	Source            string
	Destination       string
	Recursive         bool
	RemoveSourceFiles bool
	Md5sum            bool
	StopOnError       bool
	NotifyErrors      bool
}

func Upload(logger Logger, serverParams *Params, uploadParams *UploadParams) {
	uploadParams.Source = strings.TrimSuffix(uploadParams.Source, "/")
	sourceFile, err := os.Open(uploadParams.Source)
	if err != nil {
		logger.FatalWithFields(map[string]interface{}{
			"source": uploadParams.Source,
		},
			"Unable to open source path",
		)
	}
	sourceAbs, err := filepath.Abs(sourceFile.Name())
	sourceBase := filepath.Base(sourceAbs)
	if err != nil {
		logger.FatalWithFields(map[string]interface{}{
			"source": uploadParams.Source,
		},
			"Unable to get absolute path of the source",
		)
	}
	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		logger.FatalWithFields(map[string]interface{}{
			"source": uploadParams.Source,
		},
			"Unable to stat source path",
		)
	}
	sourceIsDir := sourceFileInfo.IsDir()
	logger.DebugWithFields(map[string]interface{}{
		"sourceIsDir": sourceIsDir,
	})

	uploadParams.Destination = strings.TrimSuffix(uploadParams.Destination, "/") + "/"
	bucket := strings.Split(uploadParams.Destination, "/")[0]
	objectNamePrefix := strings.TrimPrefix(uploadParams.Destination, bucket+"/")
	if objectNamePrefix != "" {
		objectNamePrefix = strings.TrimSuffix(objectNamePrefix, "/") + "/"
	}

	sourcePrefix := strings.TrimSuffix(sourceAbs, sourceBase)

	logger.DebugWithFields(map[string]interface{}{
		"source":           uploadParams.Source,
		"sourceAbs":        sourceAbs,
		"sourceBase":       sourceBase,
		"sourcePrefix":     sourcePrefix,
		"sourceIsDir":      sourceIsDir,
		"destination":      uploadParams.Destination,
		"bucket":           bucket,
		"objectNamePrefix": objectNamePrefix,
	})

	c, err := NewClient(serverParams)
	if err != nil {
		logger.Fatal(err)
	}
	logger.DebugWithFields(map[string]interface{}{
		"client": c.Client,
	},
		"MinIO client initialized",
	)

	sourceFiles := make([]string, 0)
	if sourceIsDir {
		if !uploadParams.Recursive {
			logger.Fatal(errors.New("recursive flag must be used to upload directories"))
		}
		sourceFiles, err = getFiles(sourceAbs)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		sourceFiles = append(sourceFiles, sourceAbs)
	}
	logger.Info(strconv.Itoa(len(sourceFiles)) + " files will be uploaded")

	for _, file := range sourceFiles {
		objectName := objectNamePrefix + strings.TrimPrefix(file, sourcePrefix)

		logger.DebugWithFields(map[string]interface{}{
			"file": file,
			"objectName": objectName,
		},
		"File will be uploaded",
		)

		uploadInfo, err := c.FPutObject(context.Background(), bucket, objectName, file, minio.PutObjectOptions{})
		if err != nil {
			if uploadParams.StopOnError {
				logger.Fatal(err)
			} else {
				logger.Error(err)
			}
			continue
		}
		logger.InfoWithFields(map[string]interface{}{
			"file": file,
		},
			"File uploaded",
		)
		logger.DebugWithFields(structs.Map(uploadInfo), "File uploaded")
	}

}

func getFiles(path string) ([]string, error) {
	Files := make([]string, 0)
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			Files = append(Files, filePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return Files, nil
}
