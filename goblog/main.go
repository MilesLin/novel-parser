package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"log"
	"net/http"
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

func main() {

	url := "https://blog.golang.org/slices-intro"
	e := epub.NewEpub("The Go Blog")
	e.SetAuthor("https://blog.golang.org")
	//e.SetDescription("")
	e.SetTitle("The Go Blog")
	coverImg, err := e.AddImage(cover, "")
	if err != nil {
		log.Printf("Could not add image: %v", err)
	}
	e.SetCover(coverImg, "")
	e.AddFont(font, "")
	cssName, _ := e.AddCSS(css, "")

	nextUrl := url
	for nextUrl != "" {
		req, err := http.NewRequest("GET", nextUrl, nil)
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
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalf("Could not NewDocumentFromReader because %v\n", err)
		}

		contentNode := doc.Find("#content")
		articleNode := contentNode.Find(".article")
		title := articleNode.Find(".title").Text()
		sectionBody := ""
		articleNode.Children().Each(func(i int, selection *goquery.Selection) {
			tag := selection.Nodes[0].Data
			classes, _ := selection.Attr("class")

			if tag == titleTag {
				sectionBody += fmt.Sprintf(`<h2>%s</h2>`, selection.Text())
			}

			if tag == contentTag {
				ret, _ := selection.Html()
				content := fmt.Sprintf(`<p class="%s">%s</p>`, classes, ret)

				sectionBody += content
			}

			if tag == sectionTag {
				sectionBody += fmt.Sprintf(`<h4>%s</h4>`, selection.Text())
			}

			if tag == "div" && classes == "image" {
				src, _ := selection.Find("img").Attr("src")
				imgPath := blogRoot + "/" + src
				alt, _ := selection.Find("img").Attr("alt")
				path, err := e.AddImage(imgPath, "")
				if err != nil {
					log.Printf("failed to add image: %v", err)
				}
				sectionBody += fmt.Sprintf(`<p><img style="%s" src="%s" alt="%s"/></p>`, imgMaxWidth, path, alt)
			}

			if tag == codePreTag {
				sectionBody += fmt.Sprintf(`<pre><code>%s</code></pre>`, selection.Text())
			}

		})
		e.AddSection(sectionBody, title, "", cssName)

		resp.Body.Close()

		fmt.Printf("%s Done\n", title)

		nextUrl = ""
	}

	err = e.Write("Go Blog.epub")
	if err != nil {
		log.Fatalf("Could not wirte to epub %v", err)
	}
	fmt.Println("epub 建立完成")

}
