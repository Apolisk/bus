package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gocolly/colly/v2"
)

type timetable struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func main() {
	if err := getTimetable("#u23501-89"); err != nil {
		return
	}
}

func getTimetable(id string) (err error) {

	c := colly.NewCollector()

	c.OnHTML(id, func(e *colly.HTMLElement) {
		url := "https://obukhivtrans.com.ua/"
		text := e.ChildTexts("h6")
		link := e.ChildAttrs("a", "href")

		err = writeToFile(url, text, link)
		if err != nil {
			return
		}

	})
	err = c.Visit("https://obukhivtrans.com.ua/%d1%80%d0%be%d0%b7%d0%ba%d0%bb%d0%b0%d0%b4%d0%b8-%d1%80%d1%83%d1%85%d1%83.html")
	if err != nil {
		return err
	}
	return nil
}

func writeToFile(url string, text, link []string) (err error) {
	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	for i := range text {
		r, _ := http.Get(url + link[i])
		if r.StatusCode == 404 {
			link = append(link[:i], link[i+1:]...)
		}

		data, err := json.MarshalIndent(timetable{Title: text[i], Link: url + link[i]}, "", "  ")
		if err != nil {
			return err
		}

		_, err = file.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}
