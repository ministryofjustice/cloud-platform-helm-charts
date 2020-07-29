package main

import (
	"fmt"
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
    VpcID string `json:"value"`
}


func main() {
	
	jsonFile, err := os.Open("terraform.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var outputs Outputs

	json.Unmarshal(byteValue, &outputs)

    for i := 0; i < len(outputs.Outputs); i++ {
        fmt.Println("Vpc: " + outputs.Outputs[i].Network.VpcID)
    }

}
