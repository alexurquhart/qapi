package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"alexurquhart.com/goquestrade"
)

func main() {
	refresh := os.Getenv("REFRESH_TOKEN")
	client, err := qapi.NewClient(refresh)

	if err != nil {
		fmt.Println(err)
		return
	}

	file, _ := ioutil.ReadFile("sp500.json")
	symbols := make(map[string]string)
	err = json.Unmarshal(file, &symbols)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("export REFRESH_TOKEN=" + client.Credentials.RefreshToken + "\n\n")

}
