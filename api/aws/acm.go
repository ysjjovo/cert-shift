package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"
)
var acmClient *acm.Client
func InitACMClient() {
	acmClient = acm.NewFromConfig(cfg)
}

func ImportCert(pub, key, chain, certArn string) error {
	_, err := acmClient.ImportCertificate(context.TODO(), &acm.ImportCertificateInput{
		Certificate: []byte(pub),
		PrivateKey: []byte(key),
		CertificateChain: []byte(chain),
		CertificateArn: &certArn,
	})
	return err
}