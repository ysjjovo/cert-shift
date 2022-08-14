// This file is auto-generated, don't edit it. Thanks.
package ali

import (
	cli "github.com/alibabacloud-go/cas-20180713/client"
	util "github.com/alibabacloud-go/tea-utils/service"
)

func GetAliCert (certId int64) (*cli.DescribeUserCertificateDetailResponse, error) {
  describeUserCertificateDetailRequest := &cli.DescribeUserCertificateDetailRequest{
	CertId: &certId,
  }
  runtime := &util.RuntimeOptions{}
    return client.DescribeUserCertificateDetailWithOptions(describeUserCertificateDetailRequest, runtime)
}
