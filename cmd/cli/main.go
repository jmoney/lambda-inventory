package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"gopkg.in/yaml.v3"
)

var (
	configFile = flag.String("config", "config.yaml", "config file")
)

type Config struct {
	Accounts map[string]AccountConfig `yaml:"accounts"`
}

type AccountConfig struct {
	Profile string   `yaml:"profile"`
	Regions []string `yaml:"regions"`
}

type FunctionInventory struct {
	FunctionArn string             `json:"functionArn"`
	Runtime     *string            `json:"runtime,omitempty"`
	PackageType string             `json:"packageType"`
	Tags        map[string]*string `json:"tags,omitempty"`
}

func main() {
	flag.Parse()

	output := make(map[string]map[string]map[string][]FunctionInventory)
	config := Config{}
	err := yaml.Unmarshal(Must(os.ReadFile(*configFile)), &config)
	if err != nil {
		panic(err)
	}

	for account, config := range config.Accounts {
		for _, region := range config.Regions {
			svc := lambda.New(session.Must(session.NewSession(
				aws.NewConfig().
					WithRegion(region).
					WithCredentials(credentials.NewSharedCredentials("", config.Profile))),
			))

			lfo := Must(svc.ListFunctions(&lambda.ListFunctionsInput{}))

			for {
				for _, f := range lfo.Functions {
					function := Must(svc.GetFunction(&lambda.GetFunctionInput{
						FunctionName: f.FunctionName,
					}))

					functionInventory := FunctionInventory{
						FunctionArn: *function.Configuration.FunctionArn,
						PackageType: *function.Configuration.PackageType,
						Runtime:     function.Configuration.Runtime,
						Tags:        function.Tags,
					}

					name := function.Tags["Name"]
					if name == nil {
						name = function.Configuration.FunctionName
					}
					lambdaInventory, ok := output[*name]
					if ok {
						accountInventory, ok := lambdaInventory[account]
						if ok {
							regionalInventory, ok := accountInventory[region]
							if ok {
								accountInventory[region] = append(regionalInventory, functionInventory)
							} else {
								accountInventory[region] = []FunctionInventory{functionInventory}
							}
						} else {
							accountInventory = map[string][]FunctionInventory{
								region: {functionInventory},
							}
						}
					} else {
						output[*name] = map[string]map[string][]FunctionInventory{
							account: {
								region: {functionInventory},
							},
						}
					}
				}

				if lfo.NextMarker == nil {
					break
				}
				lfo = Must(svc.ListFunctions(&lambda.ListFunctionsInput{
					Marker: lfo.NextMarker,
				}))
			}
		}
	}
	fmt.Println(string(Must(json.Marshal(output))))
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
