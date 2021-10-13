package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	*minio.Client
	params *Params
}

type Params struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Secure    bool
}

func NewClient(params *Params) (*Client, error) {
	minioClient, err := minio.New(params.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(params.AccessKey, params.SecretKey, ""),
		Secure: params.Secure,
	})
	if err != nil {
		return nil, err
	}

	c := Client{minioClient, params}

	return &c, nil
}

