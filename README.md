## 命令行格式

### 安装依赖
```shell
go mod tidy
```

### 上传文件
```shell
go run main.go -action upload -bucket your-s3-bucket-name -file path/to/your/file -access your-access-key-id -secret your-secret-key -region ap-northeast-1
```

### 下载文件
```shell
go run main.go -action download -bucket your-s3-bucket-name -file path/to/your/object -output path/to/local/file -access your-access-key-id -secret your-secret-key -region ap-northeast-1
```
