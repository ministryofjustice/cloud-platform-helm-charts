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

	// Now that we have the VPC Id, we can use this to get the actual resources within this VPC using golang aws sdk.
	// We can then compare the actual resources with those in the state file downloaded above.


	// This function outputs NatGateway details using the golang sdk for aws. However more details outputted than required. 
	// Would be good to have function that simply outputs a single value (i.e just the nat id) for a given vpc id specified. 
	DescribeNatGateway(vpcIdStr)
	

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

	//Output the NatGateway description to a file so we can read.

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
	
	jq2 := gojsonq.New().File("describeNatGateways.json")
	natGatewayId := jq2.Find("NatGateways.NatGatewayId")
	fmt.Println(natGatewayId)

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