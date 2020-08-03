package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    //"strconv"
)

// Users struct which contains
// an array of users
type Outputs struct {
    Outputs []Output `json:"outputs"`
}

// User struct which contains a name
// a type and a list of social links
type Output struct {
    NetworkID NetworkID `json:"network_id"`
}

// Social struct which contains a
// list of links
type NetworkID struct {
    Value string `json:"value"`

}

func main() {
    // Open our jsonFile
    jsonFile, err := os.Open("teraform.tfstate")
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Successfully Opened terraform.json")
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    // read our opened xmlFile as a byte array.
    byteValue, _ := ioutil.ReadAll(jsonFile)

 
    var outputs Outputs

 
    json.Unmarshal(byteValue, &outputs)

 
    for i := 0; i < len(outputs.Outputs); i++ {
        fmt.Println("Network ID: " + outputs.Outputs[i].NetworkID.Value)
    }

}