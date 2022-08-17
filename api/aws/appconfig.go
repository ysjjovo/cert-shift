package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appconfig"
)

var appconfigClient *appconfig.Client

func InitAppconfigClient() {
	appconfigClient = appconfig.NewFromConfig(cfg)
}
func GetConfig(appId, configId, clientId, envId string) (*appconfig.GetConfigurationOutput, error) {
	return appconfigClient.GetConfiguration(context.TODO(), &appconfig.GetConfigurationInput{
		Application: &appId,
		Configuration: &configId,
		ClientId: &clientId,
		Environment: &envId,
	})
}
