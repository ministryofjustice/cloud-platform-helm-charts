package main

import (
	"os"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "bytes"
    "encoding/json"
    "errors"
    "log"
    "net/http"
	"time"
	"strings"
	"strconv"
)

type SlackRequestBody struct {
    Text string `json:"text"`
}

func main() {
	var awsRegion string
	awsRegion = "eu-west-2"
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(awsRegion)},))
	var maxResults int64
	maxResults = 500
	var nextToken string
	var sgSlice []string
	var numOfUnusedSg string

	sgSlice = append(sgSlice, "| SG GROUP ID          | SG GROUP NAME ")
	
	webhookUrl := "https://hooks.slack.com/services/T02DYEB3A/<TOKEN>"
	
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
				sgSlice = append(sgSlice, "| "+*sg.GroupId+"     |     "+*sg.GroupName+" |")
		
			}
		}
		
		if descSGOut.NextToken != nil {
			nextToken = *descSGOut.NextToken
		} else {
			break
		}
	}
	

	numOfUnusedSg = strconv.Itoa(len(sgSlice))

	err := SendSlackNotification(webhookUrl, numOfUnusedSg+" currently unused SG groups in "+awsRegion+".\n"+ strings.Join(sgSlice, " \n"))
	if err != nil {
		log.Fatal(err)
	}

}


// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func SendSlackNotification(webhookUrl string, msg string) error {

    slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
    req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    if buf.String() != "ok" {
        return errors.New("Non-ok response returned from Slack")
    }
    return nil
}



