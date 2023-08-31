# lambda-inventory

Go program to inventory AWS Lambda functions across multiple accounts and regions.  It executes based off of a configuration file written in YAML to determine which accounts and regions to scan.  It currently only catalogs the `Runtime`, `PackageType`, and `Tags` but more can be added easily.

## Overview

| Arguemment | Description | Default Value |
| --- | --- | --- |
| config | The path to the config file to configure the inventory run | config.yaml |

## Example Config File

```yaml
accounts:
    dev:
        profile: dev-admin
        regions:
            - us-east-1
            - us-west-2
    production:
        profile: production-admin
        regions:
            - us-east-1
            - us-west-2
```

The above configuration would inventory the lambdas in the `dev` and `prod` account using the profile `dev-admin` and `production-admin` in the regions `us-east-1` and `us-west-2`.

## Example Output

```json
{
    "lambda1": {
        "dev": {
            "us-east-1": [
                {
                    "functionArn": "arn:aws:lambda:us-east-1:1234567891234:function:lambda1",
                    "runtime": "go1.x",
                    "packageType": "Zip",
                    "tags": {
                        "Key": "Value",
                    }
                }
            ]
        },
        "production": {
            "us-east-1": [
                {
                    "functionArn": "arn:aws:lambda:us-east-1:3434567891234:function:lambda1",
                    "runtime": "go1.x",
                    "packageType": "Zip",
                    "tags": {
                        "Key": "Value",
                    }
                }
            ],
            "us-west-2": [
                {
                    "functionArn": "arn:aws:lambda:us-west-2:3434567891234:function:lambda1",
                    "runtime": "go1.x",
                    "packageType": "Zip",
                    "tags": {
                        "Key": "Value",
                    }
                }
            ]
        }
    },
}
```
