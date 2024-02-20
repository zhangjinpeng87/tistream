package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SnsClient is the client to access SNS service.
type SnsClient struct {
	// The SNS client.
	client *sns.SNS

	// The subscription ARN.
	subscriptionArn string
}

// NewSnsClient creates a new SNS client.
func NewSnsClient(region, accessKey, secureKey string) *SnsClient {
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secureKey, ""),
	}))

	return &SnsClient{
		client: sns.New(session),
	}
}

// ListenTopic listens to a topic.
func (c *SnsClient) ListenTopic(topic string) error {
	_, err := c.client.Subscribe(&sns.SubscribeInput{
		Protocol: aws.String("http"),
		TopicArn: aws.String(topic),
		Endpoint: aws.String(c.subscriptionArn),
	})
	return err
}

// SqsClient is the client to access SQS service.
type SqsClient struct {
	// The SQS client.
	client *sqs.SQS

	// The queue URL.
	queueURL string
}