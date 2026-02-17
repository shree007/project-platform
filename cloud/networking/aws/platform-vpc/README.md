# aws-vpc-go

This example creates a VPC, subnet, internet gateway, and route table using the AWS SDK for Go v2.

Prerequisites:

- Go 1.20+
- AWS credentials configured (e.g., `~/.aws/credentials` or environment variables `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
- Set `AWS_REGION` env var, e.g. `export AWS_REGION=us-east-1`

Usage:

1. From the project directory:

```bash
cd /Users/hero/src/project-platform/cloud/networking/aws/platform-vpc
go mod tidy
go run main.go
```

2. The program will print created resource IDs.

Important:

- Running the program will create real AWS resources and may incur costs. Review the code before running and delete resources in the AWS Console when finished.
