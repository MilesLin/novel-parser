package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	//DownloadPageToFile()

}

func DownloadPageToFile() {
	var result Model
	var nextUrl string = "https://www.banxia.co/112_112058/25677257.html"
	var exist bool
	for nextUrl != "" {
		req, err := http.NewRequest("GET", nextUrl, nil)
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

		if resp.StatusCode != 200 {
			log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
		}
		transformData := DetermineEncoding(resp.Body)

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(transformData))
		if err != nil {
			log.Printf("Could not NewDocumentFromReader because %v\n", err)
		}

		title := doc.Find("#nr_title").Text()
		title = strings.ReplaceAll(title, "免費閱讀", "")
		content := doc.Find("#nr1").Text()

		result.Data = append(result.Data, Content{
			Title:   title,
			Content: strings.Split(content, "\n"),
		})

		resp.Body.Close()
		fmt.Printf("%s Done\n", title)

		nextUrl, exist = doc.Find(".next a").Attr("href")
		if !exist {
			nextUrl = ""
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	b, _ := json.Marshal(result)
	err := ioutil.WriteFile("temp.txt", []byte(b), 0644)
	if err != nil {
		log.Printf("Could not write to file %v", err)
	}
}

//func main() {
//	result := Model{
//		Data: make([]Content, 0),
//	}
//
//	for i := 25677251; i <= 25677281; i++ {
//		//https://www.banxia.co/112_112058/25677251.html
//		//https://www.banxia.co/112_112058/25676926.html
//		req, err := http.NewRequest("GET", fmt.Sprintf("https://www.banxia.co/112_112058/%d.html", i), nil)
//		if err != nil {
//			log.Printf("Could not create request because %v\n", err)
//		}
//		client := http.Client{
//			Timeout: 10 * time.Second,
//		}
//		resp, err := client.Do(req)
//		if err != nil {
//			log.Printf("Could not send request because %v\n", err)
//		}
//
//		if resp.StatusCode != 200 {
//			log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
//		}
//		transformData := DetermineEncoding(resp.Body)
//
//		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(transformData))
//		if err != nil {
//			log.Printf("Could not NewDocumentFromReader because %v\n", err)
//		}
//
//		title := doc.Find("#nr_title").Text()
//		content := doc.Find("#nr1").Text()
//		result.Data = append(result.Data, Content{
//			Title:   title,
//			Content: strings.Split(content, "\n"),
//		})
//
//		resp.Body.Close()
//		fmt.Printf("%d Done\n", i)
//		time.Sleep(500 * time.Millisecond)
//	}
//
//	//content = strings.ReplaceAll(content, string(32), "")
//	//content = strings.ReplaceAll(content, string(), "")
//
//	fmt.Println("輸出 PDF")
//
//	r := pdfGenerator.NewRequestPdf("")
//
//	//html template path
//	templatePath := "doctemplate.html"
//
//	//path for download pdf
//	outputPath := "example.pdf"
//
//	if err := r.ParseTemplate(templatePath, result); err == nil {
//		ok, _ := r.GeneratePDF(outputPath)
//		fmt.Println(ok, "pdf generated successfully")
//	} else {
//		fmt.Println(err)
//	}
//
//	// 移除免費閱讀字樣
//	// 加大字體 28 好像還太小的
//
//}

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
