package fileupload

import (
	"errors"
	"mime/multipart"
)

var (
	ErrrNotAllowExt = errors.New("不支持文件类型")
	ErrFileTooLarge = errors.New("文件太大")
)

// FileUploader 保存文件接口
type FileUploader interface {
	SaveFile(multipart.File, *multipart.FileHeader) (string, error)
}

func IsUploaderError(err error) bool {
	if errors.Is(err, ErrrNotAllowExt) || errors.Is(err, ErrFileTooLarge) {
		return true
	}
	return false
}
