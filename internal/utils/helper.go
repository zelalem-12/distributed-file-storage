package utils

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sync"

	"github.com/zelalem-12/distributed-file-storage/internal/domain"
)

const chunkSize = 1024 * 1024 // 1MB

func GetFileInfo(path string) (int64, string, string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, "", "", err
	}

	fileSize := fileInfo.Size()
	fileExtension := filepath.Ext(path)
	mimeType := mime.TypeByExtension(fileExtension)

	return fileSize, mimeType, fileExtension, nil
}

func DownloadFileInParallel(fileData *domain.File, writer io.Writer) error {

	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	partCh := make(chan []byte, 10)

	fileSize := fileData.GetSize()
	numParts := int(fileSize / chunkSize)
	if fileSize%chunkSize != 0 {
		numParts++
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	absoluteFilePath := filepath.Join(workingDirectory, fileData.GetPath())

	openFile, err := os.Open(absoluteFilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer openFile.Close()

	for i := 0; i < numParts; i++ {
		wg.Add(1)

		start := int64(i) * chunkSize
		end := start + chunkSize
		if end > fileSize {
			end = fileSize
		}

		go func(partNum int, start, end int64) {
			defer wg.Done()

			part := make([]byte, end-start)
			_, err := openFile.ReadAt(part, start)
			if err != nil && err != io.EOF {
				errCh <- fmt.Errorf("error reading file part %d: %w", partNum, err)
				return
			}

			partCh <- part
		}(i, start, end)
	}

	go func() {
		wg.Wait()
		close(partCh)
		close(errCh)
	}()

	for part := range partCh {
		_, err := writer.Write(part)
		if err != nil {
			return fmt.Errorf("error writing file part: %w", err)
		}
	}

	select {
	case err := <-errCh:
		return err
	default:
	}

	return nil
}
