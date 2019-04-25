package cmd

import (
	"fmt"

	"github.com/cnych/sinaimgmover/cloud/aliyun"
	"github.com/spf13/cobra"
)

var (
	ossBucket   string
	ossKey      string
	ossSecret   string
	ossEndpoint string
)

var ossCmd = &cobra.Command{
	Use:   "oss",
	Short: "迁移到阿里云 OSS",
	Long:  `迁移到阿里云 OSS 服务，通过 oss 命令指定将微博图床迁移到阿里云 OSS 服务。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化阿里云uploader
		uploader, err := aliyun.NewAliOSS(ossBucket, ossEndpoint, ossKey, ossSecret)
		if err != nil {
			fmt.Printf("初始化OSS客户端出错了: %s\n", err.Error())
			return
		}
		// 开始迁移
		StartMover(uploader)
	},
}

func init() {
	ossCmd.Flags().StringVarP(&ossBucket, "bucket", "b", "", "指定Aliyun  OSS Bucket")
	ossCmd.Flags().StringVarP(&ossKey, "key", "k", "", "Aliyun OSS Key")
	ossCmd.Flags().StringVarP(&ossSecret, "secret", "s", "", "Aliyun OSS Secret")
	ossCmd.Flags().StringVarP(&ossEndpoint, "endpoint", "e", "oss-cn-beijing.aliyuncs.com", "OSS Endpoint（不包含http(s)），默认值：oss-cn-hangzhou.aliyuncs.com")
	ossCmd.MarkFlagRequired("bucket")
	ossCmd.MarkFlagRequired("key")
	ossCmd.MarkFlagRequired("secret")
}
