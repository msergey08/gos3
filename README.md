# gos3 : Simple no frills AWS S3 Library using REST with V4 Signing

## Overview [![GoDoc](https://godoc.org/github.com/animber-coder/gos3?status.svg)](https://godoc.org/github.com/animber-coder/gos3) [![Go Report Card](https://goreportcard.com/badge/github.com/animber-coder/gos3)](https://goreportcard.com/report/github.com/animber-coder/gos3) [![GoCover](https://gocover.io/_badge/github.com/animber-coder/gos3)](https://gocover.io/_badge/github.com/animber-coder/gos3)

GoS3 is a golang library for uploading and deleting objects on S3 buckets using REST API calls or Presigned URLs signed using AWS Signature Version 4.

## Install

```sh
go get github.com/animber-coder/gos3
```

## Example

```go
testTxt, _ := os.Open("testdata/test.txt")
defer testTxt.Close()

// Create an instance of the package
// You can either create by manually supplying credentials
// (preferably using Environment vars)
s3 := gos3.New(Region, AWSAccessKey, AWSSecretKey)
// or you can use this on an EC2 instance to 
// obtain credentials from IAM attached to the instance.
s3 := gos3.NewUsingIAM(Region)

// You can also set a custom endpoint to a compatible s3 instance. 
s3.SetEndpoint(CustomEndpoint)

// Note: Consider adding a testTxt.Seek(0, 0)
// in case you have read 
// the body, as the pointer is shared by the library.

// File Upload is as simple as providing the following
// details.
resp, err := s3.FileUpload(gos3.UploadInput{
    Bucket:      AWSBucket,
    ObjectKey:   "test.txt",
    ContentType: "text/plain",
    FileName:    "test.txt",
    Body:        testTxt,
})

// Similarly, Files can be deleted.
err := s3.FileDelete(gos3.DeleteInput{
    Bucket:    os.Getenv("AWS_S3_BUCKET"),
    ObjectKey: "test.txt",
})

// You can also use this library to generate
// Presigned URLs that can for eg. be used to
// GET/PUT files on S3 through the browser.
var time, _ = time.Parse(time.RFC1123, "Fri, 24 May 2013 00:00:00 GMT")

url := s.GeneratePresignedURL(PresignedInput{
    Bucket:        "examplebucket",
    ObjectKey:     "test.txt",
    Method:        "GET",
    Timestamp:     time,
    ExpirySeconds: 86400,
})
```

## Contributing

You are more than welcome to contribute to this project. Fork and make a Pull Request, or create an Issue if you see any problem or want to propose a feature.

## Author

Rohan Verma <hello@rohanverma.net>

## License

MIT.
