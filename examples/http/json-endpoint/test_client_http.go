//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	var url string = "http://localhost:7000/echo"
	var msg string = "hello world"

	flag.StringVar(&url, "endpoint", url, "endpoint")
	flag.StringVar(&msg, "msg", msg, "echo message")

	flag.Parse()

	body, err := json.Marshal(map[string]any{
		"message": msg,
	})
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Printf("response: %s\n", body)
}
