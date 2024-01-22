package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type (
	chunk struct {
		index        int
		tempFileName string
	}

	file struct {
		name         string
		size         int
		etag         string
		acceptRanges bool
		contentType  string
	}
)

var (
	client = http.DefaultClient
)

func downloadFile(ctx context.Context, chunkSize int, remainder int, errs chan error, f *file) {
	fmt.Printf("Downloading %s \n", f.name)
	chunks := make(chan chunk, *numberOfChunks)

	var eg errgroup.Group

	for i := 0; i < *numberOfChunks; i++ {
		index := i
		eg.Go(func() error {
			start := index * chunkSize
			end := start + chunkSize - 1
			if index == *numberOfChunks-1 {
				end += remainder
			}
			tmpFileName, err := downloadChunk(ctx, *url, start, end)
			if err != nil {
				return err
			}
			select {
			case chunks <- chunk{index: index, tempFileName: tmpFileName}:
			case <-ctx.Done():
				fmt.Println("Operation canceled")
				return ctx.Err()
			}
			return nil
		})
	}
	go func() {
		err := eg.Wait()
		close(chunks)
		if err != nil {
			errs <- err
		}
	}()

	fileNames := make([]string, *numberOfChunks)
	defer cleanupFiles(fileNames, errs)

	for chunk := range chunks {
		fileNames[chunk.index] = chunk.tempFileName
		fmt.Printf("Chunk %s downloaded\n", chunk.tempFileName)
	}
	if err := mergeFiles(f.name, fileNames); err != nil {
		errs <- err
	} else {
		fmt.Println("File downloaded successfully.")
	}

}

func getFileDetails(url string) (*file, error) {
	resp, err := client.Head(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	cd := resp.Header.Get("Content-Disposition")
	index := strings.Index(cd, "filename=")
	var filename string
	if index > -1 {
		filename = cd[index+len("filename="):]
	}

	etag := resp.Header.Get("ETag")

	acceptRanges := resp.Header.Get("Accept-Ranges") == "bytes"

	ct := resp.Header.Get("Content-Type")

	f := &file{
		name:         filename,
		size:         size,
		etag:         etag,
		acceptRanges: acceptRanges,
		contentType:  ct,
	}

	return f, nil
}

func downloadChunk(ctx context.Context, url string, start, end int) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Add("Range", rangeHeader)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return saveChunk(resp.Body)
}

func saveChunk(body io.ReadCloser) (string, error) {
	tmpFile, err := os.CreateTemp("", "chunk")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func mergeFiles(dst string, src []string) error {
	fmt.Printf("Merging %d files into %s\n", len(src), dst)
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	closers := make([]io.Closer, 0, len(src))
	defer func() {
		for _, closer := range closers {
			closer.Close()
		}
	}()

	for _, tmpFileName := range src {
		tmpFile, err := os.OpenFile(tmpFileName, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		closers = append(closers, tmpFile)

		_, err = io.Copy(file, tmpFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupFiles(files []string, errs chan error) {
	fmt.Println("Cleaning up temporary files...")
	for _, file := range files {
		if file == "" {
			continue
		}
		fmt.Printf("Removing %s\n", file)
		err := os.Remove(file)
		if err != nil {
			errs <- err
		}
	}
}

func generateRandomFileName(f *file) string {
	ext, err := mime.ExtensionsByType(f.contentType)
	if err != nil || len(ext) == 0 {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%d", time.Now().UnixNano()) + ext[0]
}
