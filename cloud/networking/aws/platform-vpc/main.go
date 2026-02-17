package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	ctx := context.Background()
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ec2c := ec2.NewFromConfig(cfg)

	// 1) Create VPC
	cidr := "10.0.0.0/16"
	createVpcOut, err := ec2c.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidr),
	})
	if err != nil {
		log.Fatalf("CreateVpc failed: %v", err)
	}
	vpcId := *createVpcOut.Vpc.VpcId
	fmt.Println("Created VPC:", vpcId)

	// Tag VPC
	_, _ = ec2c.CreateTags(ctx, &ec2.CreateTagsInput{
		Resources: []string{vpcId},
		Tags:      []types.Tag{{Key: aws.String("Name"), Value: aws.String("example-vpc")}},
	})

	// Wait a short moment for the VPC to be available
	time.Sleep(2 * time.Second)

	// 2) Create a subnet
	subCidr := "10.0.1.0/24"
	createSubnetOut, err := ec2c.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		VpcId:     aws.String(vpcId),
		CidrBlock: aws.String(subCidr),
	})
	if err != nil {
		log.Fatalf("CreateSubnet failed: %v", err)
	}
	subnetId := *createSubnetOut.Subnet.SubnetId
	fmt.Println("Created Subnet:", subnetId)

	// 3) Create Internet Gateway
	igwOut, err := ec2c.CreateInternetGateway(ctx, &ec2.CreateInternetGatewayInput{})
	if err != nil {
		log.Fatalf("CreateInternetGateway failed: %v", err)
	}
	igwId := *igwOut.InternetGateway.InternetGatewayId
	fmt.Println("Created Internet Gateway:", igwId)

	// Attach IGW to VPC
	_, err = ec2c.AttachInternetGateway(ctx, &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
		VpcId:             aws.String(vpcId),
	})
	if err != nil {
		log.Fatalf("AttachInternetGateway failed: %v", err)
	}
	fmt.Println("Attached IGW to VPC")

	// 4) Create a route table and route to IGW
	rtOut, err := ec2c.CreateRouteTable(ctx, &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcId),
	})
	if err != nil {
		log.Fatalf("CreateRouteTable failed: %v", err)
	}
	rtId := *rtOut.RouteTable.RouteTableId
	fmt.Println("Created Route Table:", rtId)

	_, err = ec2c.CreateRoute(ctx, &ec2.CreateRouteInput{
		RouteTableId:         aws.String(rtId),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(igwId),
	})
	if err != nil {
		log.Fatalf("CreateRoute failed: %v", err)
	}
	fmt.Println("Created route to IGW")

	// Associate route table with subnet
	_, err = ec2c.AssociateRouteTable(ctx, &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rtId),
		SubnetId:     aws.String(subnetId),
	})
	if err != nil {
		log.Fatalf("AssociateRouteTable failed: %v", err)
	}
	fmt.Println("Associated route table with subnet")

	// 5) (Optional) Enable auto-assign public IP on subnet
	_, err = ec2c.ModifySubnetAttribute(ctx, &ec2.ModifySubnetAttributeInput{
		SubnetId:            aws.String(subnetId),
		MapPublicIpOnLaunch: &types.AttributeBooleanValue{Value: aws.Bool(true)},
	})
	if err != nil {
		log.Fatalf("ModifySubnetAttribute failed: %v", err)
	}
	fmt.Println("Enabled auto-assign public IP on subnet")

	fmt.Println("All resources created. VPC ID:", vpcId)
}
