package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
    "os"
	
)

type Users struct {
    Users []User `json:"users"`
}

type User struct {
    Social Social `json:"social"`
}


type Social struct {
    Facebook string `json:"facebook"`
}


func main() {
	
	jsonFile, err := os.Open("users.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users Users 

	json.Unmarshal(byteValue, &users)

    for i := 0; i < len(users.Users); i++ {
        fmt.Println("Vpc: " + users.Users[i].Social.Facebook)
    }

}