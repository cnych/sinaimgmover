package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/cnych/sinaimgmover/cloud"
	"github.com/cnych/sinaimgmover/utils"
	"github.com/spf13/cobra"
)

var (
	postPath    string
	imagePrefix string
	nameLength  int
	ossCommand  *flag.FlagSet
)

const sinaImg = "sinaimg.cn"

var rootCmd = &cobra.Command{
	Use:   "mover",
	Short: "Mover 是一个用于将微博图床一键迁移到云服务的工具",
	Long: `基于 Golang 开发的一个用于将微博图床一键迁移到云服务的工具
		目前只支持阿里云 OSS，文档查看：https://github.com/cnych/sinaimgmover`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("执行 mover -h 命令查看使用方法")
	},
}

func init() {

	rootCmd.AddCommand(ossCmd)

	rootCmd.PersistentFlags().StringVarP(&imagePrefix, "prefix", "f", "images", "Bucket下面的文件夹目录，默认值： images")
	rootCmd.PersistentFlags().IntVarP(&nameLength, "length", "l", 6, "指定上传到OSS上面的图片名称长度，默认值：6")
	rootCmd.PersistentFlags().StringVarP(&postPath, "post", "p", "./", "指定markdown文章路径，默认值：当前目录")

}

func sinaImgToCloud(url string, uploader cloud.Uploader) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	objectKey := fmt.Sprintf("%s/%s.jpg", imagePrefix, utils.RandID(nameLength))
	key, err := uploader.Upload(objectKey, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.%s/%s", ossBucket, ossEndpoint, key), nil
}

func parseFile(filePath string, uploader cloud.Uploader, wg *sync.WaitGroup) error {
	bt, err := ioutil.ReadFile(filePath)
	if err != nil {
		wg.Done()
		return err
	}
	content := string(bt)
	// # Markdown中图片语法 ![](url) 或者 <img src='' />
	re := regexp.MustCompile(`!\[.*?\]\((.*?)\)|<img.*?src=[\'\"](.*?)[\'\"].*?>`)
	params := re.FindAllStringSubmatch(content, -1)
	// 获取所有的微博图床图片
	for _, param := range params {
		imgURL := param[1]
		if strings.Index(imgURL, sinaImg) != -1 {
			cloudURL, err := sinaImgToCloud(imgURL, uploader)
			if err != nil {
				fmt.Printf("图片：%s 转换到 OSS 出错了\n", err.Error())
			} else {
				newContent := strings.Replace(content, imgURL, cloudURL, -1)
				//重新写入
				ioutil.WriteFile(filePath, []byte(newContent), 0)
				content = newContent
				fmt.Printf("成功替换了图片：%s\n", imgURL)
			}
		}
	}
	wg.Done()
	return nil
}

// StartMover ... 开始执行迁移程序
func StartMover(uploader cloud.Uploader) {
	files, err := utils.GetAllFiles(postPath)
	if err != nil {
		fmt.Printf("获取markdown文件出错了: %s\n", err.Error())
	}
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(fpath string) {
			err := parseFile(fpath, uploader, &wg)
			if err != nil {
				fmt.Printf("解析文件：%s 出错了：%s\n", fpath, err.Error())
			}
		}(file)
	}
	wg.Wait()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
