package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)
var snsClient *sns.Client
func InitSNSClient() {
	snsClient = sns.NewFromConfig(cfg)
}
func SnsPublish(msg, snsTopicArn string) error{
	_, err := snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: &snsTopicArn,
		Message:  &msg,
	})
	return err
}