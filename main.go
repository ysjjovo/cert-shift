package main

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/yaml.v3"
	"ysjjovo.ml/cert-shift/api/ali"
	"ysjjovo.ml/cert-shift/api/aws"
)

type Config struct {
	Ali struct {
		Ak     string `yaml:"ak"`
		Sk     string `yaml:"sk"`
		CertId int64  `yaml:"certId"`
	}
	Aws struct {
		Ak          string `yaml:"ak"`
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
	var cfgBytes []byte

	aws.InitAppconfigClient()
	res, err := aws.GetConfig(
		os.Getenv("APP_ID"),
		os.Getenv("CONFIG_ID"),
		os.Getenv("CLIENT_ID"),
		os.Getenv("ENV_ID"),
	)
	if err != nil {
		println("getConfig err", err.Error())
		println("now tring local env!")
		dir, err = os.Getwd()
		if err != nil {
			println("get config dir error", err)
			os.Exit(1)
		}
		cfgBytes, err = ioutil.ReadFile(dir + "/config.yml")
		if err != nil {
			println("read config file error", err.Error())
			os.Exit(1)
		}
	} else {
		cfgBytes = res.Content
	}
	if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		println("yaml convert to map error", err)
		os.Exit(1)
	}
	ali.InitClient(cfg.Ali.Ak, cfg.Ali.Sk)
	aws.InitConfig(cfg.Aws.Ak, cfg.Aws.Sk, cfg.Aws.Region)
	aws.InitACMClient()
	aws.InitSNSClient()
	println("init completed!")
}
func log(msg string) {
	if err := aws.SnsPublish(msg, cfg.Aws.SnsTopicArn); err != nil {
		println("publish error", err)
	}
}
func handler(ctx context.Context) (string, error) {
	res, err := ali.GetCert(cfg.Ali.CertId)
	if err != nil {
		msg := "getAliCert error" + err.Error()
		println(msg)
		log(msg)
		return "", err
	}
	cert := *res.Body.Cert

	sep := "-----END CERTIFICATE-----"
	sp := strings.Split(cert, sep)
	pub := sp[0] + sep + "\n"
	chain := sp[1] + sep + "\n"
	key := *res.Body.Key

	// println("certArn", cfg.Aws.CertArn, "pub", pub, "key", key, "chain", chain)
	if err := aws.ImportCert(pub, key, chain, cfg.Aws.CertArn); err != nil {
		msg := "importCert error" + err.Error()
		println(msg)
		log(msg)
		return "", err
	}
	println("import cert completed!")
	return "success", nil
}

func main() {
	// handler(context.TODO())
	lambda.Start(handler)
}
