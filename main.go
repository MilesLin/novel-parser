package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	//DownloadPageToFile()

	// Create a new EPUB
	e := epub.NewEpub("施落衛琮曦-錦繡小福妻")

	// Set the author
	e.SetAuthor("三妖")
	e.SetDescription("穿越成宰相庶女的施落，還沒過一天好日子，就被老皇帝賜給貶為庶人雙腿殘廢的衛家小王爺……醒來後，沒有丫環，沒有錦衣玉食，只有一個雙腿殘廢的帥哥陰鷙的看著她。“休書我不會寫，你想死隨時都可以，生是衛家人，死也是衛家的鬼！”施落看著窮的只剩老鼠的家，為了能吃飽穿暖活下去，只能想辦法賺錢養家，賺錢養夫，賺錢養娃…衛小王爺多疑敏感不好養，原以為養成了一只白眼狼，誰知道一朝功成名就，他居然帶人堵上門。“听說你跟別人說我死了？”施落︰“呵呵……誤會！”某王爺將門一關，“那我這個死人再娶你一次可好？”")
	e.SetTitle("施落衛琮曦-錦繡小福妻")

	// Add a section
	b, err := ioutil.ReadFile("temp.txt")
	if err != nil {
		log.Fatal("Could not read the file %v", err)
	}
	var model Model
	json.Unmarshal(b, &model)
	for _, v := range model.Data {

		section1Body := `<h2>{{.Title}}</h2>{{range .Content}}<p>{{.}}</p>{{end}}`
		var tpl bytes.Buffer
		temp := template.New("novel")
		t, err := temp.Parse(section1Body)
		if err != nil {
			log.Fatalf("Could not parse %v", err)
		}
		err = t.Execute(&tpl, v)
		if err != nil {
			log.Fatalf("Could not Execute %v", err)
		}

		e.AddSection(string(tpl.Bytes()), v.Title, "", "")
		fmt.Println(v.Title)
	}

	// Write the EPUB
	err = e.Write("施落衛琮曦-錦繡小福妻.epub")
	if err != nil {
		log.Fatalf("Could not wirte to epub %v", err)
	}

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
