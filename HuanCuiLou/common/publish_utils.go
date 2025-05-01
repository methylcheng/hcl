package common

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"
)

// OSS配置
const (
	endpoint        = "https://oss-cn-hangzhou.aliyuncs.com"
	accessKeyID     = "你的AccessKeyId"
	accessKeySecret = "你的AccessKeySecret"
	bucketName      = "你的Bucket名"
)

// UploadToOSS 上传图片到 OSS, 返回图片的url路径和error
func UploadToOSS(file *multipart.FileHeader, filePath string) (string, error) {
	// 初始化 OSS 客户端
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Println("OSS客户端初始化失败:", err)
		return "", fmt.Errorf("OSS客户端初始化失败:%v", err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Println("OSS客户端初始化失败:", err)
		return "", fmt.Errorf("OSS客户端初始化失败:%v", err)
	}
	// 上传到 OSS
	filename := fmt.Sprintf("cover/%d_%s", time.Now().Unix(), filepath.Base(file.Filename)) // 命名文件
	if err := bucket.PutObjectFromFile(filename, filePath); err != nil {
		return "", fmt.Errorf("OSS上传失败:%v", err)
	}
	// 构造文件访问地址
	url := fmt.Sprintf("https://%s.%s/%s", bucketName, endpoint[8:], filename)
	return url, nil
}
