package ali

import (
	cli "github.com/alibabacloud-go/cas-20180713/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

var client = &cli.Client{}

func InitClient(ak string, sk string) (err error) {
	config := &openapi.Config{
		AccessKeyId:     &ak,
		AccessKeySecret: &sk,
	}
	config.Endpoint = tea.String("cas.aliyuncs.com")
	client, err = cli.NewClient(config)
	return err
}
