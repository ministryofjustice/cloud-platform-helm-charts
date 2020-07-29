package main

import (
	"os"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "log"
)

type SlackRequestBody struct {
    Text string `json:"text"`
}

func main() {
	var awsRegion string
	//awsRegion = os.Getenv("AWS_DEFAULT_REGION")
	awsRegion = "us-west-1"
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(awsRegion)},))
	var maxResults int64
	maxResults = 500
	var nextToken string

	for {
		descSGInput := &ec2.DescribeSecurityGroupsInput{
			MaxResults: aws.Int64(maxResults),
		}

		if nextToken != "" {
			descSGInput.NextToken = aws.String(nextToken)
		}

		descSGOut, err := svc.DescribeSecurityGroups(descSGInput)
		if err != nil {
			log.Println("Unable to get security groups. Err:", err)
			os.Exit(1)
		}

		
        // Loop through all ec2s and get describe the network interface as per SG group id
		for _, sg := range descSGOut.SecurityGroups {
			niInput := &ec2.DescribeNetworkInterfacesInput{
				Filters: []*ec2.Filter{&ec2.Filter{
					Name:   aws.String("group-id"),
					Values: []*string{sg.GroupId},
				}},
			}

			niOut, err := svc.DescribeNetworkInterfaces(niInput)
			if err != nil {
				log.Println("Unable to get network interfaces for the security group. Err:", err)
				os.Exit(1)
			}
            
			if len(niOut.NetworkInterfaces) == 0 { //If any of the SGs are not attached to a network interface then add them to slice
				log.Println("Unused Security Group: ", *sg.GroupId, *sg.GroupName)
			}
		}
		
		if descSGOut.NextToken != nil {
			nextToken = *descSGOut.NextToken
		} else {
			break
		}
	}
	
}

