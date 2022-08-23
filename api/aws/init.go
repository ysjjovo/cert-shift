package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)
var cfg aws.Config
func InitConfig() error{
	var err error
	cfg, err = config.LoadDefaultConfig(context.TODO())
	return err
}