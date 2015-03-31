package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	refresh := os.Getenv("REFRESH_TOKEN")
	client, err := NewClient(refresh)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("export REFRESH_TOKEN=" + client.Credentials.RefreshToken + "\n\n")

	a, err := client.GetCandles(8049, time.Now().AddDate(-1, 0, 0), time.Now(), "OneDay")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a)
}
