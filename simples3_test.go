package gos3

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestS3_FileUpload(t *testing.T) {
	testTxt, err := os.Open("test.txt")
	if err != nil {
		return
	}
	defer testTxt.Close()
	testPng, err := os.Open("avatar.png")
	if err != nil {
		return
	}
	defer testPng.Close()

	type fields struct {
		AccessKey string
		SecretKey string
		Region    string
	}
	type args struct {
		u UploadInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Upload test.txt",
			fields: fields{
				AccessKey: os.Getenv("AWS_S3_ACCESS_KEY"),
				SecretKey: os.Getenv("AWS_S3_SECRET_KEY"),
				Region:    os.Getenv("AWS_S3_REGION"),
			},
			args: args{
				UploadInput{
					Bucket:      os.Getenv("AWS_S3_BUCKET"),
					ObjectKey:   "test.txt",
					ContentType: "text/plain",
					FileName:    "test.txt",
					Body:        testTxt,
				},
			},
			wantErr: false,
		},
		{
			name: "Upload avatar.png",
			fields: fields{
				AccessKey: os.Getenv("AWS_S3_ACCESS_KEY"),
				SecretKey: os.Getenv("AWS_S3_SECRET_KEY"),
				Region:    os.Getenv("AWS_S3_REGION"),
			},
			args: args{
				UploadInput{
					Bucket:      os.Getenv("AWS_S3_BUCKET"),
					ObjectKey:   "xyz/image.png",
					ContentType: "image/png",
					FileName:    "avatar.png",
					Body:        testPng,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s3 := New(tt.fields.Region, tt.fields.AccessKey, tt.fields.SecretKey)
			resp, err := s3.FileUpload(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3.FileUpload() error = %v, wantErr %v", err, tt.wantErr)
			}
			// check for empty response
			if (resp == UploadResponse{}) {
				t.Errorf("S3.FileUpload() returned empty response, %v", resp)
			}
		})
	}
}

func TestS3_FileDelete(t *testing.T) {
	type fields struct {
		AccessKey string
		SecretKey string
		Region    string
	}
	type args struct {
		u DeleteInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete test.txt",
			fields: fields{
				AccessKey: os.Getenv("AWS_S3_ACCESS_KEY"),
				SecretKey: os.Getenv("AWS_S3_SECRET_KEY"),
				Region:    os.Getenv("AWS_S3_REGION"),
			},
			args: args{
				DeleteInput{
					Bucket:    os.Getenv("AWS_S3_BUCKET"),
					ObjectKey: "test.txt",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s3 := New(tt.fields.Region, tt.fields.AccessKey, tt.fields.SecretKey)
			if err := s3.FileDelete(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("S3.FileDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestS3_NewUsingIAM(t *testing.T) {
	var (
		iam  = `test-new-s3-using-iam`
		resp = `{"Code" : "Success","LastUpdated" : "2018-12-24T10:18:01Z",
				"Type" : "AWS-HMAC","AccessKeyId" : "abc",
				"SecretAccessKey" : "abc","Token" : "abc",
				"Expiration" : "2018-12-24T16:24:59Z"}`
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
		if r.URL.EscapedPath() == "/" {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, iam)
		}
		if r.URL.EscapedPath() == "/"+iam {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, resp)
		}
	}))
	defer ts.Close()

	s3, err := newUsingIAMImpl(ts.URL, "abc")
	if err != nil {
		t.Errorf("S3.FileDelete() error = %v", err)
	}
	if s3.AccessKey != "abc" && s3.SecretKey != "abc" && s3.Region != "abc" {
		t.Errorf("S3.FileDelete() got = %v", s3)
	}
}
