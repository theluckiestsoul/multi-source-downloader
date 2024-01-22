package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetFileDetails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *file
		wantErr bool
		before  func(url string)
		after   func()
	}{
		{
			name: "Should return file details",
			args: args{
				url: "https://examplefile.com/file-download/25",
			},
			want: &file{
				name:         "test.txt",
				size:         100,
				etag:         "123456789",
				acceptRanges: true,
			},
			wantErr: false,
			before: func(url string) {
				httpmock.RegisterResponder(http.MethodHead, url, func(r *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, "")
					resp.Header.Set("Content-Disposition", "attachment; filename=test.txt")
					resp.Header.Set("Content-Length", "100")
					resp.Header.Set("Accept-Ranges", "bytes")
					resp.Header.Set("ETag", "123456789")
					return resp, nil
				})
			},
			after: func() {
				httpmock.Reset()
			},
		},
		{
			name: "Should return error when http request fails",
			args: args{
				url: "https://examplefile.com/file-download/26",
			},
			want:    nil,
			wantErr: true,
			before: func(url string) {
				httpmock.RegisterResponder(http.MethodHead, url, func(r *http.Request) (*http.Response, error) {
					return nil, errors.New("internal server error")
				})
			},
			after: func() {
				httpmock.Reset()
			},
		},
		{
			name: "Should return error when content length is not an integer",
			args: args{
				url: "https://examplefile.com/file-download/27",
			},
			want:    nil,
			wantErr: true,
			before: func(url string) {
				httpmock.RegisterResponder(http.MethodHead, url, func(r *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, "")
					resp.Header.Set("Content-Length", "abc")
					return resp, nil
				})
			},
			after: func() {
				httpmock.Reset()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(tt.args.url)
			defer tt.after()
			got, err := getFileDetails(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) == tt.wantErr {
				return
			}
			if got.name != tt.want.name {
				t.Errorf("getFileDetails() got = %v, want %v", got.name, tt.want.name)
			}
			if got.size != tt.want.size {
				t.Errorf("getFileDetails() got = %v, want %v", got.size, tt.want.size)
			}
			if got.etag != tt.want.etag {
				t.Errorf("getFileDetails() got = %v, want %v", got.etag, tt.want.etag)
			}
			if got.acceptRanges != tt.want.acceptRanges {
				t.Errorf("getFileDetails() got = %v, want %v", got.acceptRanges, tt.want.acceptRanges)
			}
		})
	}
}

func TestGenerateRandomFileName(t *testing.T) {
	tests := []struct {
		name string
		file *file
		want string
	}{
		{
			name: "Should return a bin file extension",
			file: &file{
				contentType: "application/octet-stream",
			},
			want: ".bin",
		},
		{
			name: "Should return a pdf file extension",
			file: &file{
				contentType: "application/pdf",
			},
			want: ".pdf",
		},
		{
			name: "Should return a txt file extension",
			file: &file{
				contentType: "application/json",
			},
			want: ".json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := generateRandomFileName(tt.file)
			if filepath.Ext(n) != tt.want {
				t.Errorf("generateRandomFileName() got = %v, want %v", filepath.Ext(n), tt.want)
			}
		})
	}
}

func TestCleanupFilesSuccess(t *testing.T) {
	var fileNames []string
	for i := 0; i < 10; i++ {
		f, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		fileNames = append(fileNames, f.Name())
		f.Close()
	}

	errs := make(chan error, len(fileNames))
	done := make(chan struct{})

	go func() {
		cleanupFiles(fileNames, errs)
		close(errs)
		done <- struct{}{}
	}()

	go func() {
		for err := range errs {
			t.Errorf("cleanupFiles() error = %v", err)
		}
	}()

	_ = <-done
	close(done)

	for _, file := range fileNames {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			t.Errorf("cleanupFiles() file %s still exists", file)
		}
	}
}

func TestCleanupFilesError(t *testing.T) {
	var fileNames []string
	for i := 0; i < 10; i++ {
		f, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		fileNames = append(fileNames, f.Name())
		f.Close()
	}
	fileNames = append(fileNames, "non-existent-file")

	errs := make(chan error, len(fileNames))
	done := make(chan struct{})

	go func() {
		cleanupFiles(fileNames, errs)
		close(errs)
		done <- struct{}{}
	}()

	go func() {
		for err := range errs {
			t.Errorf("cleanupFiles() error = %v", err)
		}
	}()

	_ = <-done
	close(done)

	for _, file := range fileNames {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			t.Errorf("cleanupFiles() file %s still exists", file)
		}
	}
}

func TestMergeFiles(t *testing.T) {
	var fileNames []string
	expectedContent := ""

	for i := 0; i < 11; i++ {
		f, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}

		fileNames = append(fileNames, f.Name())

		content := fmt.Sprintf("%d", i)
		_, err = f.WriteString(content)
		if err != nil {
			t.Fatal(err)
		}
		expectedContent += content
		f.Close()
	}

	dst, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	dst.Close()

	err = mergeFiles(dst.Name(), fileNames)
	if err != nil {
		t.Errorf("mergeFiles() error = %v", err)
	}

	//check destination file content
	content, err := os.ReadFile(dst.Name())
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expectedContent {
		t.Errorf("mergeFiles() got = %v, want %v", string(content), expectedContent)
	}
}
