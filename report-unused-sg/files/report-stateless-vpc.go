package main

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"log"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	
)

func main() {

	DownloadFromS3Bucket("cloud-platform-terraform-state/cloud-platform-network/iawan-test/", "terraform.tfstate")

	jq := gojsonq.New().File("terraform.tfstate")
	vpcId := jq.Find("outputs.network_id.value")
	fmt.Println(vpcId)
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