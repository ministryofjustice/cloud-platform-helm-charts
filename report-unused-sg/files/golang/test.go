package main

import (
	"fmt"
	//"os/exec"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
//	"log"
	"github.com/aws/aws-sdk-go/service/s3"
	//"encoding/json"
//	"github.com/tidwall/gjson"
 //   "bytes"
	"encoding/json"
	"io/ioutil"
    "os"
	
)

type Outputs struct {
    Outputs []Output `json:"outputs"`
}

type Output struct {
    Network Network `json:"network_id"`
}


type Network struct {
    VpcId string `json:"value"`
}


func main() {
	
//	err := ListS3Buckets()
//	if err != nil {
//		log.Fatal(err)
//	}

//	cmd := exec.Command("terraform", "init")
//	stdout, err := cmd.Output()

//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println(string(stdout))


//	cmd2 := exec.Command("terraform", "show", "-json")
//	stdout2, err := cmd2.Output()
	

//    if err != nil {
//        fmt.Println(err.Error())
 //       return
 //   }

//	fmt.Println(string(stdout2))

	//b, _ := json.Marshal(stdout2)
	//s := string(b)
    //fmt.Println(s)



//	const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

    
//	vpcid := gjson.Get(string(stdout2), "network_id.value")
//	fmt.Println(vpcid)


	jsonFile, err := os.Open("terraform.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened terraform.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	//var result map[string]interface{}
	//json.Unmarshal([]byte(byteValue), &result)
	
	var outputs Outputs 

	json.Unmarshal(byteValue, &outputs)

    for i := 0; i < len(outputs.Outputs); i++ {
        fmt.Println("Vpc: " + outputs.Outputs[i].Network.VpcId)
    }

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

