package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cnych/sinaimgmover/cloud"
	"github.com/cnych/sinaimgmover/cloud/aliyun"
	"github.com/cnych/sinaimgmover/utils"
)

var (
	postPath    string
	ossBucket   string
	ossKey      string
	ossSecret   string
	ossEndpoint string
	ossFolder   string
	nameLength  int
	bucket      *oss.Bucket
	uploader    cloud.Uploader
)

const sinaImg = "sinaimg.cn"

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&nameLength, "length", 6, "指定上传到OSS上面的图片名称长度，默认为6")
	flag.StringVar(&postPath, "post", "./", "指定markdown文章路径，默认当前目录")
	flag.StringVar(&ossBucket, "bucket", "", "指定Aliyun  OSS Bucket")
	flag.StringVar(&ossKey, "key", "", "Aliyun OSS Key")
	flag.StringVar(&ossSecret, "secret", "", "Aliyun OSS Secret")
	flag.StringVar(&ossFolder, "folder", "images", "Bucket下面的文件夹目录，默认为 images")
	flag.StringVar(&ossEndpoint, "endpoint", "oss-cn-beijing.aliyuncs.com", "OSS Endpoint（不包含http(s)），如：oss-cn-hangzhou.aliyuncs.com")
}

func initBucket() (*oss.Bucket, error) {
	client, err := oss.New(ossEndpoint, ossKey, ossSecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(ossBucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func exit(msg string) {
	flag.Usage()
	fmt.Fprintln(os.Stderr, "\n[Error] "+msg)
	os.Exit(1)
}

func sinaImgToCloud(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	objectKey := fmt.Sprintf("%s/%s.jpg", ossFolder, utils.RandID(nameLength))
	key, err := uploader.Upload(objectKey, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.%s/%s", ossBucket, ossEndpoint, key), nil
}

func parseFile(filePath string, wg *sync.WaitGroup) error {
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
			cloudURL, err := sinaImgToCloud(imgURL)
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

func main() {
	flag.Parse()
	if ossBucket == "" {
		exit("没有指定 bucket 参数")
	}
	if ossEndpoint == "" {
		exit("没有指定 endpoint 参数")
	}
	if ossKey == "" {
		exit("没有指定 key 参数")
	}
	if ossSecret == "" {
		exit("没有指定 secret 参数")
	}
	files, err := utils.GetAllFiles(postPath)
	if err != nil {
		fmt.Printf("获取markdown文件出错了: %s\n", err.Error())
	} else {
		bucket, err = initBucket()
		// todo，根据参数传入决定使用哪个uploader
		uploader = aliyun.NewAliOSS(bucket)

		if err != nil {
			fmt.Printf("初始化OSS客户端出错了: %s\n", err.Error())
		} else {
			var wg sync.WaitGroup
			for _, file := range files {
				wg.Add(1)
				go func(fpath string) {
					err := parseFile(fpath, &wg)
					if err != nil {
						fmt.Printf("解析文件：%s 出错了：%s\n", fpath, err.Error())
					}
				}(file)
			}
			wg.Wait()
		}
	}
}
