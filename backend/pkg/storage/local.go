package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{basePath: absPath}, nil
}

// sanitizePath validates and sanitizes the key to prevent path traversal attacks
func (s *LocalStorage) sanitizePath(key string) (string, error) {
	// Clean the path to remove any .. or . elements
	cleanKey := filepath.Clean(key)
	
	// Prevent absolute paths
	if filepath.IsAbs(cleanKey) {
		return "", errors.New("absolute paths are not allowed")
	}
	
	// Prevent paths that start with .. after cleaning
	if strings.HasPrefix(cleanKey, "..") {
		return "", errors.New("path traversal detected")
	}
	
	// Join with base path and verify it's still within base directory
	fullPath := filepath.Join(s.basePath, cleanKey)
	
	// Ensure the full path is within the base directory
	if !strings.HasPrefix(fullPath, s.basePath+string(filepath.Separator)) && fullPath != s.basePath {
		return "", errors.New("invalid path: outside base directory")
	}
	
	return fullPath, nil
}

func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader) error {
	filePath, err := s.sanitizePath(key)
	if err != nil {
		return err
	}
	
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	filePath, err := s.sanitizePath(key)
	if err != nil {
		return nil, err
	}
	return os.Open(filePath)
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	filePath, err := s.sanitizePath(key)
	if err != nil {
		return err
	}
	return os.Remove(filePath)
}
