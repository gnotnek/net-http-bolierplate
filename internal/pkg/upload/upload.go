package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type ChunkInfo struct {
	UploadID    string
	ChunkIndex  int
	TotalChunks int
	Filename    string
}

type ChunkedUploader struct {
	storagePath string
	mu          sync.Mutex
}

func NewChunkedUploader(storagePath string) *ChunkedUploader {
	return &ChunkedUploader{storagePath: storagePath}
}

// SaveChunk saves a single chunk to a temporary file
func (u *ChunkedUploader) SaveChunk(info ChunkInfo, chunkData io.Reader) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	uploadDir := filepath.Join(u.storagePath, info.UploadID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

	chunkPath := filepath.Join(uploadDir, fmt.Sprintf("chunk_%d", info.ChunkIndex))
	out, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, chunkData)
	return err
}

// AssembleChunks combines all chunks into the final file
func (u *ChunkedUploader) AssembleChunks(info ChunkInfo) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	uploadDir := filepath.Join(u.storagePath, info.UploadID)
	finalPath := filepath.Join(u.storagePath, info.Filename)
	out, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	for i := 0; i < info.TotalChunks; i++ {
		chunkPath := filepath.Join(uploadDir, fmt.Sprintf("chunk_%d", i))
		in, err := os.Open(chunkPath)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(out, in)
		in.Close()
		if err != nil {
			return "", err
		}
	}

	// Optionally, clean up chunks
	os.RemoveAll(uploadDir)
	return finalPath, nil
}

func (u *ChunkedUploader) StoragePath() string {
	return u.storagePath
}
