package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	req, err := http.NewRequest("GET", "https://www.banxia.co/112_112058/25676926.html", nil)
	if err != nil {
		log.Printf("Could not create request because %v\n", err)
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Could not send request because %v\n", err)
	}
	defer resp.Body.Close()
	transformData := DetermineEncoding(resp.Body)
	fmt.Print(string(transformData))
}

func DetermineEncoding(r io.Reader) []byte {
	OldReader := bufio.NewReader(r)
	bytes, err := OldReader.Peek(1024)
	if err != nil {
		panic(err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	reader := transform.NewReader(OldReader, e.NewDecoder())
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return all
}
