package main

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"ysjjovo.ml/cert-shift/api/ali"
	"ysjjovo.ml/cert-shift/api/aws"
	// "github.com/aws/aws-lambda-go/lambda"
	// "fmt"
)

type Config struct {
	Ali struct {
		AK     string `yaml:"ak"`
		Sk     string `yaml:"sk"`
		CertId int64  `yaml:"certId"`
	}
	Aws struct {
		AK          string `yaml:"ak"`
		Sk          string `yaml:"sk"`
		Region      string `yaml:"region"`
		CertArn     string `yaml:"certArn"`
		SnsTopicArn string `yaml:"snsTopicArn"`
	}
}

var cfg Config

func init() {
	var err error
	var dir string
	dir, err = os.Getwd()
	if err != nil {
		println("get config dir error", err)
	}
	var cfgBytes []byte
	cfgBytes, err = ioutil.ReadFile(dir + "/config.yml")
	if err != nil {
		println("read config file error", err.Error())
	}
	if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		println("yaml convert to map error", err)
	}
	ali.InitAliClient(cfg.Ali.AK, cfg.Ali.Sk)
	aws.InitConfig(cfg.Aws.AK, cfg.Aws.Sk, cfg.Aws.Region)
	aws.InitACMClient()
	aws.InitSNSClient()
}
func log(msg string) {
	if err := aws.SnsPublish(msg, cfg.Aws.SnsTopicArn); err != nil {
		println("publish error", err)
	}
}
func handler(ctx context.Context) (string, error) {
	res, err := ali.GetAliCert(cfg.Ali.CertId)
	if err != nil {
		msg := "getAliCert error" + err.Error()
		println(msg)
		log(msg)
	}
	cert := *res.Body.Cert

	sep := "-----END CERTIFICATE-----"
	sp := strings.Split(cert, sep)
	pub := sp[0] + sep + "\n"
	chain := sp[1] + sep + "\n"
	key := *res.Body.Key

	println("certArn", cfg.Aws.CertArn, "pub", pub, "key", key, "chain", chain)
	if err := aws.ImportCert(pub, key, chain, cfg.Aws.CertArn); err != nil {
		msg := "importCert error" + err.Error()
		println(msg)
		log(msg)
	}
	return "success", nil
}

func main() {
	handler(context.TODO())
	// lambda.Start(HandleRequest)
}
