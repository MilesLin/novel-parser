package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mileslin/novel-parser/pdfGenerator"
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

	if resp.StatusCode != 200 {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	transformData := DetermineEncoding(resp.Body)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(transformData))
	//doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Could not NewDocumentFromReader because %v\n", err)
	}
	content := doc.Find("#nr1").Text()
	title := doc.Find("#nr_title").Text()

	r := pdfGenerator.NewRequestPdf("")
	// 32
	// 160

	//html template path
	templatePath := "sample.html"

	//path for download pdf
	outputPath := "example.pdf"

	//html template data
	templateData := struct {
		Title       string
		Description string
	}{
		Title:       title,
		Description: content,
	}

	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		ok, _ := r.GeneratePDF(outputPath)
		fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
	}

	// 剩下組合資料到 HTML 上
	// 還有調整文字大小

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
