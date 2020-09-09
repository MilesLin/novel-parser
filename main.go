package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"github.com/mileslin/novel-parser/pdfGenerator"
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
type Chap struct {
	Url string
	Title string
}
func main() {

	filename := "temp.txt"

	//chaps := make([]Chap, 0)
	//
	//// 取得所有文章連結
	//// https://m.23ts.net/book/208125/asc/index_1.html
	//for i := 1; i < 50; i++ {
	//	url := fmt.Sprintf("https://m.23ts.net/book/208125/asc/index_%d.html", i)
	//	resp, err := http.Get(url)
	//	if err != nil {
	//		log.Println(err)
	//		continue
	//	}
	//	transformData := determineEncoding(resp.Body)
	//	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(transformData))
	//
	//	doc.Find("#jieqi_page_contents").Find("a").Each(func(i int, selection *goquery.Selection) {
	//		attr, exists := selection.Attr("href")
	//		if exists {
	//			chaps = append(chaps, Chap{
	//				Url:   attr,
	//				Title:  selection.Text(),
	//			})
	//			log.Println(attr, selection.Text())
	//		}
	//	})
	//
	//	resp.Body.Close()
	//}
	//
	//marshal, _ := json.Marshal(chaps)
	//err := ioutil.WriteFile("url.txt", marshal, 0660)
	//if err != nil {
	//	log.Fatal(err)
	//}


	// 先把網站內容下載成 txt ，這樣之後讀取檔案的時候只要從 txt 讀取就好
	//downloadPageToFile(filename)
	// 建立 epub
	createEpubFromFile(filename)
	// 建立 pdf
	//createPDFFromFile(filename)

}

func createEpubFromFile(filename string) {
	e := epub.NewEpub("梁知語完結大結局")
	e.SetAuthor("諾小穎")
	e.SetDescription("梁知語完結大結局")
	e.SetTitle("梁知語完結大結局")

	// 讀取下載好的資料
	b, err := ioutil.ReadFile(filename)
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
	err = e.Write("梁知語完結大結局.epub")
	if err != nil {
		log.Fatalf("Could not wirte to epub %v", err)
	}
	fmt.Println("epub 建立完成")
}

func createPDFFromFile(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Could not read the file %v", err)
	}
	var model Model
	json.Unmarshal(b, &model)
	r := pdfGenerator.NewRequestPdf("")

	templatePath := "pdfGenerator/doctemplate.html"

	outputPath := "example.pdf"

	if err := r.ParseTemplate(templatePath, model); err == nil {
		_, err := r.GeneratePDF(outputPath)
		if err != nil {
			log.Fatalf("Could not generate pdf %v", err)
		}
		fmt.Println("pdf 建立成功")
	} else {
		log.Fatalf("Could not ParseTemplate %v", err)
	}
}

// 將網站內容儲存到 txt
func downloadPageToFile(filename string) {

	// get all url
	file, _ := ioutil.ReadFile("url.txt")
	urls := make([]Chap, 0)
	json.Unmarshal(file, &urls)

	var result Model
	for i := 0; i < len(urls) ; i++ {
		req, err := http.NewRequest("GET", urls[i].Url, nil)
		if err != nil {
			log.Printf("Could not create request because %v\n", err)
			continue
		}
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Could not send request because %v\n", err)
			continue
		}

		if resp.StatusCode != 200 {
			log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
			continue
		}

		transformData := determineEncoding(resp.Body)

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(transformData))
		if err != nil {
			log.Fatalf("Could not NewDocumentFromReader because %v\n", err)
		}

		title := urls[i].Title
		//title = strings.ReplaceAll(urls[i].Title, "免費閱讀", "")
		content := doc.Find("#acontent").Text()

		result.Data = append(result.Data, Content{
			Title:   title,
			Content: strings.Split(content, "\n"),
		})

		resp.Body.Close()
		fmt.Printf("%s Done\n", title)

		// 防止被以為是 DDOS 攻擊
		time.Sleep(50 * time.Millisecond)
	}

	b, _ := json.Marshal(result)
	err := ioutil.WriteFile(filename, []byte(b), 0644)
	if err != nil {
		log.Fatalf("Could not write to file %v", err)
	}
	fmt.Printf("%s 檔案建立成功\n", filename)
}

// 用來修正中文字體是亂碼問題
func determineEncoding(r io.Reader) []byte {
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
