package webfile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// MimeType returns the media type of a multipart.File object.
func MimeType(f multipart.File, h *multipart.FileHeader) (string, error) {
	mimeType := mime.TypeByExtension(filepath.Ext(h.Filename))
	if len(mimeType) > 0 {
		return mimeType, nil
	}

	buff := make([]byte, 512)
	_, err := f.Read(buff)
	if err != nil {
		return "", fmt.Errorf("error reading the first 512 bytes: %w", err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("error seeking file: %w", err)
	}

	return http.DetectContentType(buff), nil
}

// Hash computes a hash of a given multipart.File.
// Apply an optional salt to create a unique hash if needed.
func Hash(f multipart.File, h *multipart.FileHeader, salt []byte) (string, error) {
	buff := make([]byte, h.Size)
	_, err := f.Read(buff)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("error seeking file: %w", err)
	}

	buff = append(buff, salt...)
	hasher := sha256.New()
	_, err = hasher.Write(buff)
	if err != nil {
		return "", fmt.Errorf("error hashing file: %w", err)
	}
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}
