package main

import (
	"fmt"
	"os/exec"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
//	"log"
	"github.com/aws/aws-sdk-go/service/s3"
)



func main() {
	
//	err := ListS3Buckets()
//	if err != nil {
//		log.Fatal(err)
//	}

	cmd := exec.Command("terraform", "init")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))


	cmd2 := exec.Command("terraform", "state", "show", "module.webops_ecr_scan_repos_s3_bucket.aws_s3_bucket.bucket")
    stdout2, err := cmd2.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

	fmt.Println(string(stdout2))



}


func ListS3Buckets() error {

	var awsRegion string
	awsRegion = "eu-west-1"
	s3svc := s3.New(session.New(&aws.Config{Region: aws.String(awsRegion)},))
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		fmt.Println("Failed to list buckets", err)
		return err
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	}
	
	return nil
}

