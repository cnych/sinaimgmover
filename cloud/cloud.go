package cloud

import "io"

// Uploader ... 文件上传统一接口
type Uploader interface {
	Upload(objectKey string, r io.Reader) (string, error)
}
