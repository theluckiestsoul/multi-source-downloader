package main

import (
	"errors"
	"net/http"
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
