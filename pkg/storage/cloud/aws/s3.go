package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Client is the client to access S3 service.
type S3Client struct {
	// The S3 client.
	client *s3.S3
}

// NewS3Client creates a new S3 client.
func NewS3Client(region, accessKey, secureKey, bucket string) *S3Client {
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secureKey, ""),
	}))

	return &S3Client{
		client: s3.New(session),
	}
}
