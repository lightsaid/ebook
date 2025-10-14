package fileupload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/lightsaid/ebook/pkg/random"
)

// LocalUploader 本地文件上传实现结构体
type LocalUploader struct {
	uploadDir string   // 保存文件地址
	allowExts []string // 允许文件类型
	maxBytes  int64    // 文件大小限制
}

// 类型检查
var _ FileUploader = (*LocalUploader)(nil)

// NewLocalUplader 创建一个本地上传图片实例, dir 文件存储地址, allowExts 允许文件格式，max 文件最大限制
func NewLocalUplader(dir string, allowExts []string, max int64) FileUploader {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	if max <= 0 {
		max = 2 << 20 // 1MB
	}
	return &LocalUploader{
		uploadDir: dir,
		allowExts: allowExts,
		maxBytes:  max,
	}
}

// IsAllowExt 允许文件类型 fileExt => .xxx
func (l *LocalUploader) IsAllowExt(fileExt string) bool {
	for _, ext := range l.allowExts {
		if ext == fileExt {
			return true
		}
	}
	return false
}

// SaveFile 本地上传，实现接口
func (l *LocalUploader) SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := path.Ext(header.Filename)
	if allow := l.IsAllowExt(ext); !allow {
		return "", fmt.Errorf("%w:%s; allow ext: %v", ErrrNotAllowExt, ext, l.allowExts)
	}

	if header.Size > l.maxBytes {
		return "", ErrFileTooLarge
	}
	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 构建本地存储路径
	dstPath := l.uploadDir + random.RandomString(10) + ext

	// 创建目标文件
	newFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	// 复制文件内容
	_, err = io.Copy(newFile, src)
	if err != nil {
		return "", err
	}

	return dstPath, nil
}
