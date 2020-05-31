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
	"strings"
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
	result := Model{
		Data: make([]Content, 0),
	}
	title := doc.Find("#nr_title").Text()
	content := doc.Find("#nr1").Text()
	cc := strings.Split(content, "\n")
	//content = strings.ReplaceAll(content, string(32), "")
	//content = strings.ReplaceAll(content, string(), "")
	result.Data = append(result.Data, Content{
		Title:   title,
		Content: cc,
	})

	r := pdfGenerator.NewRequestPdf("")

	//html template path
	templatePath := "doctemplate.html"

	//path for download pdf
	outputPath := "example.pdf"

	if err := r.ParseTemplate(templatePath, result); err == nil {
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
