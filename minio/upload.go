package minio

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/fatih/structs"
	"github.com/minio/minio-go/v7"
	"io"
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

func Upload(logger Logger, notifier Notifier, serverParams *Params, uploadParams *UploadParams) {
	uploadParams.Source = strings.TrimSuffix(uploadParams.Source, "/")
	sourceFile, err := os.Open(uploadParams.Source)
	if err != nil {
		logger.FatalWithFields(map[string]interface{}{
			"source": uploadParams.Source,
		},
			"Unable to open source path",
		)
		err = notifier.Notify("Kaynak dosya/dizin açılamadı: " + uploadParams.Source)
		if err != nil {
			logger.WarnWithFields(map[string]interface{}{
				"error": err,
			},
			"Unable to send notification",
			)
		}
	}
	sourceAbs, err := filepath.Abs(sourceFile.Name())
	sourceBase := filepath.Base(sourceAbs)
	if err != nil {
		logger.FatalWithFields(map[string]interface{}{
			"source": uploadParams.Source,
		},
			"Unable to get absolute path of the source",
		)
		err = notifier.Notify("Kaynak dosya/dizinin tam yolu belirlenemedi: " + uploadParams.Source)
		if err != nil {
			logger.WarnWithFields(map[string]interface{}{
				"error": err,
			},
				"Unable to send notification",
			)
		}
	}
	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		logger.FatalWithFields(map[string]interface{}{
			"source": uploadParams.Source,
		},
			"Unable to stat source path",
		)
		err = notifier.Notify("Kaynak dosya/dizin bilgileri alınamadı: " + uploadParams.Source)
		if err != nil {
			logger.WarnWithFields(map[string]interface{}{
				"error": err,
			},
				"Unable to send notification",
			)
		}
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
		err = notifier.Notify("MinIO client oluşturulamadı: " + err.Error())
		if err != nil {
			logger.WarnWithFields(map[string]interface{}{
				"error": err,
			},
				"Unable to send notification",
			)
		}
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
			err = notifier.Notify("\"recursive\" parametresi kullanılmadığı için dizin yüklenemedi.")
			if err != nil {
				logger.WarnWithFields(map[string]interface{}{
					"error": err,
				},
					"Unable to send notification",
				)
			}
		}
		sourceFiles, err = getFiles(sourceAbs)
		if err != nil {
			logger.Fatal(err)
			err = notifier.Notify(sourceAbs + " dizinindeki dosyalar alınamadı.")
			if err != nil {
				logger.WarnWithFields(map[string]interface{}{
					"error": err,
				},
					"Unable to send notification",
				)
			}
		}
	} else {
		sourceFiles = append(sourceFiles, sourceAbs)
	}
	logger.Info(strconv.Itoa(len(sourceFiles)) + " files will be uploaded")

	uploaded := make([]string, 0)

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


		if uploadParams.Md5sum {
			md5sum, err := md5sum(file)
			if err != nil {
				if uploadParams.StopOnError {
					logger.Fatal(err)
				} else {
					logger.Error(err)
				}
				continue
			}

			if md5sum != uploadInfo.ETag {
				if uploadParams.StopOnError {
					logger.FatalWithFields(map[string]interface{}{
						"md5sum": md5sum,
						"ETag": uploadInfo.ETag,
					},
						"md5sums doesn't match",
					)
				} else {
					logger.ErrorWithFields(map[string]interface{}{
						"md5sum": md5sum,
						"ETag": uploadInfo.ETag,
					},
						"md5sums doesn't match",
					)
				}
				continue
			}

			logger.DebugWithFields(map[string]interface{}{
				"md5sum": md5sum,
				"ETag": uploadInfo.ETag,
			},
				"md5sums match",
			)

			if uploadParams.RemoveSourceFiles {
				err = os.Remove(file)
				if err != nil {
					logger.Error(err)
				}
			}
		}
		uploaded = append(uploaded, file)
	}

	if len(sourceFiles) > len(uploaded) {
		err = notifier.Notify("Bazı dosyalar MinIO'ya yüklenemedi. Lütfen logu kontrol edin.")
		if err != nil {
			logger.WarnWithFields(map[string]interface{}{
				"error": err,
			},
				"Unable to send notification",
			)
		}
	}

}

func md5sum(file string) (string, error) {
	var md5Str string

	f, err := os.Open(file)
	if err != nil {
		return md5Str, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	hash := md5.New()

	if _, err := io.Copy(hash, f); err != nil {
		return md5Str, err
	}

	hashInBytes := hash.Sum(nil)[:16]

	md5Str = hex.EncodeToString(hashInBytes)

	return md5Str, nil
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
