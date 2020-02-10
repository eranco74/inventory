package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	parser := argparse.NewParser("health", "Health check")
	// Create string flag
	ip := parser.String("i", "ip", &argparse.Options{Required: true, Help: "Ip to connect"})
	port := parser.String("p", "port", &argparse.Options{Required: true, Help: "port to connect"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}
	fmt.Println("AAAAA")
	for {
		time.Sleep(10 * time.Second)
		go looping(fmt.Sprintf("http://%s:%s", *ip, *port))
		fmt.Println("BBBBB")
	}
}

func looping(connectionString string) string {
	resp, err := http.Get(connectionString)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Println("AAAAA")
		fmt.Println(err.Error())
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(fmt.Sprintf("RESPONSE %s", body))
	return string(body)
}
