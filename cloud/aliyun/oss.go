package aliyun

import (
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliOSS ... 阿里云OSS服务
type AliOSS struct {
	Endpoint, AccessKey, AccessSecret string
	Bucket                            *oss.Bucket
}

// NewAliOSS ... 新建阿里云OSS对象
func NewAliOSS(bucket, endpoint, key, secret string) (*AliOSS, error) {
	client, err := oss.New(endpoint, key, secret)
	if err != nil {
		return nil, err
	}
	bucketObject, err := client.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	return &AliOSS{
		Endpoint:     endpoint,
		AccessKey:    key,
		AccessSecret: secret,
		Bucket:       bucketObject,
	}, nil
}

// Upload ... 阿里云OSS上传内容
func (alioss *AliOSS) Upload(objectKey string, r io.Reader) (string, error) {
	err := alioss.Bucket.PutObject(objectKey, r)
	if err != nil {
		return "", err
	}
	return objectKey, nil
}
