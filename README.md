## 程式碼說明
此專案是將網站內容轉成 epub 或 pdf 的程式範例，程式碼的執行內容都放在 main.go，主就是三件事情
1. downloadPageToFile
  * 使用 [net/http](https://golang.org/pkg/net/http/) 發 Request 到網站
  * 使用 [goquery](https://github.com/PuerkitoBio/goquery) 解析 Response 的內容
  * 將網站需要的資訊，儲存到 temp.txt ，以便重新製作 pdf 跟 epub 的時候不用重新撈資料
2. createEpubFromFile
  * 使用 [go-epub](https://github.com/bmaupin/go-epub) 建立 epub 檔案
3. createPDFFromFile
  * 參考 [Golang-HTML-TO-PDF-Converter](https://github.com/Mindinventory/Golang-HTML-TO-PDF-Converter) 建立 pdf 檔案
  * *範例是用 [wkhtmltopdf](https://wkhtmltopdf.org/) 轉 pdf，檔案放在 `pdfGenerator/wkhtmltopdf.exe`*

**執行方式**  
執行 `go run .` 就會產生 pdf 與 epub 的範例檔案。

## 參考
* 亂碼解決方案: [golang学习笔记之-采集gbk乱码的问题?](https://www.codercto.com/a/60635.html)
* HTML Parser: [goquery](https://github.com/PuerkitoBio/goquery)
* HTML TO PDF: [Golang-HTML-TO-PDF-Converter](https://github.com/Mindinventory/Golang-HTML-TO-PDF-Converter)
* 寫入 txt 檔案: [Go by Example: Writing Files](https://gobyexample.com/writing-files)
* 產生電子書 epub 檔案: [go-epub](https://github.com/bmaupin/go-epub)
  * [Example Code](https://github.com/bmaupin/go-docs-epub)

## 其他參考
* 能夠 parser 支援 reader mode 網站的套件: [go-readability](https://github.com/go-shiori/go-readability)
