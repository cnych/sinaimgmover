package aliyun

import (
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliOSS ... 阿里云OSS服务
type AliOSS struct {
	Bucket *oss.Bucket
}

// NewAliOSS ... 新建阿里云OSS对象
func NewAliOSS(bucket *oss.Bucket) *AliOSS {
	return &AliOSS{
		Bucket: bucket,
	}
}

// Upload ... 阿里云OSS上传内容
func (alioss *AliOSS) Upload(objectKey string, r io.Reader) (string, error) {
	err := alioss.Bucket.PutObject(objectKey, r)
	if err != nil {
		return "", err
	}
	return objectKey, nil
}
