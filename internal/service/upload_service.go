package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"convenienceStore/pkg/uid"
)

// UploadService handles storing uploaded files and returning their accessible path.
type UploadService interface {
	SaveFile(ctx context.Context, originalName string, content io.Reader) (string, error)
}

type uploadService struct {
	basePath string
}

// NewUploadService creates an UploadService writing files under the uploads directory.
func NewUploadService(deps Dependencies) UploadService {
	return &uploadService{basePath: "uploads"}
}

func (s *uploadService) SaveFile(ctx context.Context, originalName string, content io.Reader) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	if err := os.MkdirAll(s.basePath, 0o755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(originalName))
	filename := uid.New("file_") + ext
	fullPath := filepath.Join(s.basePath, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, content); err != nil {
		return "", err
	}

	relative := filepath.ToSlash(filepath.Join(s.basePath, filename))
	relative = strings.TrimPrefix(relative, "./")

	return relative, nil
}
