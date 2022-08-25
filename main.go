package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/pkg/errors"
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
		CertId   string `yaml:"certId"`
		SnsTopic string `yaml:"snsTopic"`
	}
}

var cfg Config

func getFuncArn(ctx context.Context) ([]string, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, errors.Errorf("could not get lambda context")
	}
	return strings.Split(lc.InvokedFunctionArn, ":"), nil
}
func getAccountId(ctx context.Context) (string, error) {
	arn, err := getFuncArn(ctx)
	if err != nil {
		return "", errors.New("getFuncArn in getAccountId failed!")
	}
	return arn[4], nil
}
func getRegion(ctx context.Context) (string, error) {
	arn, err := getFuncArn(ctx)
	if err != nil {
		return "", errors.New("getFuncArn in getRegion failed!")
	}
	return arn[3], nil
}
func getRegionAndAcc(ctx context.Context) (string, error) {
	region, err := getRegion(ctx)
	if err != nil {
		return "", err
	}
	acc, err := getAccountId(ctx)
	if err != nil {
		return "", err
	}
	return region + ":" + acc, nil
}
func getCertARN(regionAndAcc, certId string) string {
	var cn string
	if regionAndAcc[0] == 'c' && regionAndAcc[1] == 'n' {
		cn = "-cn"
	}
	return fmt.Sprintf("arn:aws%s:acm:%s:certificate/%s", cn, regionAndAcc, certId)
}
func getSNSTopicArn(regionAndAcc, topic string) string {
	var cn string
	if regionAndAcc[0] == 'c' && regionAndAcc[1] == 'n' {
		cn = "-cn"
	}
	return fmt.Sprintf("arn:aws%s:sns:%s:%s", cn, regionAndAcc, topic)
}
func init() {
	var err error
	var dir string
	var cfgBytes []byte

	aws.InitConfig()
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
	// tmp := make(map [string]interface{})
	// if err = yaml.Unmarshal(cfgBytes, tmp); err != nil {
	// 	println("yaml convert to map error", err)
	// 	os.Exit(1)
	// }
	// s, _ := json.Marshal(tmp)
	// println("s", string(s))

	if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		println("yaml convert to config error", err)
		os.Exit(1)
	}
	
	ali.InitClient(cfg.Ali.Ak, cfg.Ali.Sk)
	aws.InitACMClient()
	if cfg.Aws.SnsTopic != "" {
		aws.InitSNSClient()
	}
	println("init completed!")
}
func handler(ctx context.Context) error {
	res, err := ali.GetCert(cfg.Ali.CertId)
	if err != nil {
		msg := "getAliCert error" + err.Error()
		println(msg)
		log(ctx, msg)
		return err
	}
	cert := *res.Body.Cert

	sep := "-----END CERTIFICATE-----"
	sp := strings.Split(cert, sep)
	pub := sp[0] + sep + "\n"
	chain := sp[1] + sep + "\n"
	key := *res.Body.Key

	// println("certArn", cfg.Aws.CertArn, "pub", pub, "key", key, "chain", chain)
	regionAndAcc, err := getRegionAndAcc(ctx)
	if err != nil {
		println("getRegionAndAcc error", err)
		return err
	}
	if err := aws.ImportCert(pub, key, chain, getCertARN(regionAndAcc, cfg.Aws.CertId)); err != nil {
		msg := "importCert error" + err.Error()
		println(msg)
		log(ctx, msg)
		return err
	}
	println("import cert completed!")
	return nil
}

func main() {
	// handler(context.TODO())
	lambda.Start(handler)
}

func log(ctx context.Context, msg string) {
	if cfg.Aws.SnsTopic != "" {
		regionAndAcc, err := getRegionAndAcc(ctx)
		if err != nil {
			println("getRegionAndAcc error", err)
			return
		}
		if err := aws.SnsPublish(msg, getSNSTopicArn(regionAndAcc, cfg.Aws.SnsTopic)); err != nil {
			println("publish error", err)
		}
	}
}