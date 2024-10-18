package main

import (
    "context"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

// 上传文件到S3的函数
func uploadFileToS3(bucketName, fileName, accessKeyID, secretAccessKey, region string) error {
    customResolver := aws.CredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""))
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(customResolver), config.WithRegion(region))
    if err != nil {
        return fmt.Errorf("无法加载AWS配置: %v", err)
    }

    s3Client := s3.NewFromConfig(cfg)

    file, err := os.Open(fileName)
    if err != nil {
        return fmt.Errorf("无法打开文件 %s: %v", fileName, err)
    }
    defer file.Close()

    key := filepath.Base(fileName)

    _, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
        Body:   file,
    })
    if err != nil {
        return fmt.Errorf("上传文件到S3失败: %v", err)
    }

    fmt.Printf("文件 %s 已成功上传到S3 Bucket: %s\n", fileName, bucketName)
    return nil
}

// 从S3下载文件的函数
func downloadFileFromS3(bucketName, key, accessKeyID, secretAccessKey, region, outputFilePath string) error {
    customResolver := aws.CredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""))
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(customResolver), config.WithRegion(region))
    if err != nil {
        return fmt.Errorf("无法加载AWS配置: %v", err)
    }

    s3Client := s3.NewFromConfig(cfg)

    output, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
    })
    if err != nil {
        return fmt.Errorf("下载文件失败: %v", err)
    }
    defer output.Body.Close()

    file, err := os.Create(outputFilePath)
    if err != nil {
        return fmt.Errorf("无法创建文件 %s: %v", outputFilePath, err)
    }
    defer file.Close()

    _, err = io.Copy(file, output.Body)
    if err != nil {
        return fmt.Errorf("写入文件失败: %v", err)
    }

    fmt.Printf("文件 %s 已成功下载到 %s\n", key, outputFilePath)
    return nil
}

func main() {
    // 定义命令行参数
    action := flag.String("action", "upload", "操作类型: upload 或 download")
    bucketName := flag.String("bucket", "", "S3 Bucket 名称")
    fileName := flag.String("file", "", "上传文件路径或下载的文件名")
    outputFilePath := flag.String("output", "", "下载的文件保存路径")
    accessKeyID := flag.String("access", "", "AWS ACCESS KEY ID")
    secretAccessKey := flag.String("secret", "", "AWS SECRET ACCESS KEY")
    region := flag.String("region", "us-east-1", "AWS 区域")

    // 解析命令行参数
    flag.Parse()

    if *bucketName == "" || *accessKeyID == "" || *secretAccessKey == "" {
        log.Fatalf("Bucket、Access Key 和 Secret Key 是必需的参数")
    }

    // 根据指定的操作执行上传或下载
    switch *action {
    case "upload":
        if *fileName == "" {
            log.Fatalf("上传操作需要提供文件路径")
        }
        err := uploadFileToS3(*bucketName, *fileName, *accessKeyID, *secretAccessKey, *region)
        if err != nil {
            log.Fatalf("文件上传失败: %v", err)
        }
    case "download":
        if *outputFilePath == "" {
            log.Fatalf("下载操作需要提供输出文件路径")
        }
        err := downloadFileFromS3(*bucketName, *fileName, *accessKeyID, *secretAccessKey, *region, *outputFilePath)
        if err != nil {
            log.Fatalf("文件下载失败: %v", err)
        }
    default:
        log.Fatalf("无效的操作类型: %s. 只支持 upload 或 download", *action)
    }
}
