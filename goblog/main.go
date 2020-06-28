package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	css         = "style.css"
	imgMaxWidth = "max-width: 800px !important;"
	blogRoot    = "https://blog.golang.org"
	font        = blogRoot + "/fonts.css"
	cover       = "cover.png"
	titleTag    = "h2"
	sectionTag  = "h4"
	codePreTag  = "pre"
	contentTag  = "p"
)

type sectionBody struct {
	index int
	title string
	body string
}
func main() {

	blogs := getLink()
	e := epub.NewEpub("The Go Blog")
	e.SetAuthor("https://blog.golang.org")
	e.SetTitle("The Go Blog")
	coverImg, err := e.AddImage(cover, "")
	if err != nil {
		log.Printf("Could not add image: %v", err)
	}
	e.SetCover(coverImg, "")
	e.AddFont(font, "")
	cssName, _ := e.AddCSS(css, "")

	sectionCh := make(chan sectionBody)
	sections := make([]sectionBody, len(blogs))
	queue := make(chan struct{}, 10)
	var wg sync.WaitGroup
	wg.Add(len(blogs))
	var wgSectionCh sync.WaitGroup
	wgSectionCh.Add(1)
	go func(arr []sectionBody){
		for ch := range sectionCh {
			arr[ch.index] = sectionBody{
				title: ch.title,
				body:  ch.body,
			}
		}
		wgSectionCh.Done()
	}(sections)

	for index, url := range blogs {
		queue <- struct{}{}
		go func(index int, url string){
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("Could not create request because %v\n", err)
				return
			}
			client := http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Could not send request because %v\n", err)
				return
			}

			if resp.StatusCode != 200 {
				log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
				return
			}
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatalf("Could not NewDocumentFromReader because %v\n", err)
			}

			contentNode := doc.Find("#content")
			articleNode := contentNode.Find(".article")
			title := articleNode.Find(".title").Text()
			body := ""
			articleNode.Children().Each(func(i int, selection *goquery.Selection) {
				tag := selection.Nodes[0].Data
				classes, _ := selection.Attr("class")
				switch tag {
				case titleTag:
					body += fmt.Sprintf(`<h2>%s</h2>`, selection.Text())
				case contentTag:
					ret, _ := selection.Html()
					content := fmt.Sprintf(`<p class="%s">%s</p>`, classes, ret)
					body += content

				case sectionTag:
					body += fmt.Sprintf(`<h4>%s</h4>`, selection.Text())
				case "div":
					if classes == "image" {
						src, _ := selection.Find("img").Attr("src")
						imgPath := blogRoot + "/" + src
						alt, _ := selection.Find("img").Attr("alt")
						path, err := e.AddImage(imgPath, "")
						if err != nil {
							log.Printf("failed to add image: %v", err)
						}
						body += fmt.Sprintf(`<p><img style="%s" src="%s" alt="%s"/></p>`, imgMaxWidth, path, alt)
					}
				case codePreTag:

					bodeBody, _ := selection.Find("code").Html()
					body += fmt.Sprintf(`<pre><code>%s</code></pre>`, bodeBody)
				case "ol":
					start, _ := selection.Attr("start")
					body += fmt.Sprintf(`<p>%s. %s</p>`, start, selection.Text())
				case "ul":
					ret, _ := selection.Html()
					content := fmt.Sprintf(`<ul>%s</ul>`, ret)
					body += content
				default:
					fmt.Printf("---------- In Default: %s ----------\n", tag)
					ret, _ := selection.Html()
					content := fmt.Sprintf(`<%s>%s</%s>`, tag, ret, tag)
					body += content
				}

			})

			sectionCh <- sectionBody{
				body: body,
				title: title,
				index: index,
			}

			resp.Body.Close()
			<- queue
			wg.Done()
			fmt.Println(title)
		}(index, url)
	}

	wg.Wait()

	close(queue)
	close(sectionCh)

	fmt.Println("Waiting SectionCh")

	wgSectionCh.Wait()

	fmt.Println("Adding sections")

	for _, section := range sections {
		e.AddSection(section.body, section.title, "", cssName)
	}

	fmt.Println("Writing epub")

	err = e.Write("Go Blog.epub")
	if err != nil {
		log.Fatalf("Could not wirte to epub %v", err)
	}

	fmt.Println("epub 建立完成")

}

func getLink() []string {
	url := "https://blog.golang.org/index"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Could not create request because %v\n", err)
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Could not send request because %v\n", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Could not NewDocumentFromReader because %v\n", err)
	}
	//blogLinkNode := doc.Find(".blogtitle").FilterFunction(func(i int, selection *goquery.Selection) bool {
	//	return  strings.Contains(selection.Text(), "Using Go Modules") || strings.Contains(selection.Text(), "Go 2017 Survey Results") || strings.Contains(selection.Text(), "Go 2016 Survey Results") || strings.Contains(selection.Text(), "Seven years of Go") || strings.Contains(selection.Text(), "Text normalization in Go")
	//})
	blogLinkNode := doc.Find(".blogtitle")

	len := blogLinkNode.Length()

	result := make([]string, len)

	blogLinkNode.Each(func(i int, selection *goquery.Selection) {
		a := selection.Find("a")
		href, exists := a.Attr("href")
		if exists {
			result[len-1-i] = blogRoot + href
		}
	})

	return result
}
