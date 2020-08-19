package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

type PageProperties struct {
	Title string
}

func GetLatestBlogTitles(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	titles := ""
	doc.Find(".post-title").Each(func(i int, s *goquery.Selection) {
		titles += "- " + s.Text() + "\n"
	})
	return titles, nil
}

func getUrl(url string) {
	blogTitles, err := GetLatestBlogTitles(url)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Blog Titles:")
	fmt.Printf(blogTitles)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	p := PageProperties{Title: "Golang web scraper"}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}
func submitUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	furl := r.FormValue("url")

	data := struct {
		url string
	}{
		url: furl,
	}
	getUrl(data.url)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/submit", submitUrl)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequest()
}
