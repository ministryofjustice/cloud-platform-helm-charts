package main

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"log"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	
)

func main() {

	//Downlaod the state for 'cloud-platform-network' 
	DownloadFromS3Bucket("cloud-platform-terraform-state/cloud-platform-network/iawan-test/", "terraform.tfstate")

    //From the downloaded state file extract the VPC ID
	jq := gojsonq.New().File("terraform.tfstate")
	vpcId := jq.Find("outputs.network_id.value")


    //Convert the VPC ID into a string
	var vpcIdStr string
	vpcIdStr = fmt.Sprintf("%v", vpcId)
	fmt.Println(vpcIdStr)

	//DescribeNatGateway(vpcIdStr)
	
	DescribeVPC(vpcIdStr)

		
}

func DescribeVPC(vpcId string) {


		svc := ec2.New(session.New(&aws.Config{
		Region: aws.String("eu-west-2")},))

		input := &ec2.DescribeVpcsInput{
			VpcIds: []string{
				"vpc-0c16457fd570a1f0b",
			},
		}
		
		req := svc.DescribeVpcsRequest(input)
		result, err := req.Send(context.Background())
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		
		fmt.Println(result)
}

func DescribeNatGateway(vpcId string) {

	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String("eu-west-2")},))
	
	input := &ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcId),
				},
			},
		},
	}
	
	result, err := svc.DescribeNatGateways(input)
	if err != nil {
			fmt.Println(err.Error())
			return
	}
	
	var descNatGateways string
	descNatGateways = fmt.Sprintf("%v", result)
	fmt.Println(descNatGateways)

	//*********************

	f, err := os.Create("describeNatGateways.json")
    if err != nil {
        fmt.Println(err)
        return
    }
    l, err := f.WriteString(descNatGateways)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
	}
	fmt.Println(l, "bytes written successfully")
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
	
	//jq2 := gojsonq.New().File("describeNatGateways.json")
	//natGatewayId := jq2.Find("NatGateways.NatGatewayId")
	//fmt.Println(natGatewayId)

}

func DownloadFromS3Bucket(bucket string, item string) {
	// 2) Create an AWS session
	sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1")},
	)

	// 3) Create a new AWS S3 downloader 
	downloader := s3manager.NewDownloader(sess)

	// 4) Download the item from the bucket. If an error occurs, log it and exit. Otherwise, notify the user that the download succeeded.
	file, err := os.Create(item)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
	})

	if err != nil {
		log.Fatalf("Unable to download item %q, %v", item, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

}