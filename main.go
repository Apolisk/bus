package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	if err := getTimetable("tbody"); err != nil {
		return
	}
}

func getTimetable(tag string) (err error) {
	c := colly.NewCollector()
	c.OnHTML(tag, func(e *colly.HTMLElement) {
		link := e.ChildAttrs("a", "href")
		err = writeToFile(link)
		if err != nil {
			return
		}

	})
	err = c.Visit("http://obukhivtrans.com.ua/assets/")
	if err != nil {
		return err
	}
	return nil
}

func writeToFile(link []string) (err error) {
	var wg sync.WaitGroup
	start := time.Now()
	//for j := 0; j < runtime.NumCPU()/8; j++ {
	for i := 5; i < len(link); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Nanosecond)
			filetype := link[i][len(link[i])-3:]
			value := "http://obukhivtrans.com.ua/assets/" + link[i]
			err = saveFile(filetype, value)
			if err != nil {
				log.Println(err)
			}
		}(i)
	}
	//}
	wg.Wait()
	fmt.Println(time.Now().Sub(start).Seconds())
	return nil
}

func saveFile(filetype, value string) error {
	switch filetype {
	case "jpg":
		err := getFile(value, "images")
		if err != nil {
			return err
		}
	case "pdf":
		err := getFile(value, "pdf")
		if err != nil {
			return err
		}
	}
	return nil
}

func getFile(value, path string) error {
	decodedValue, err := url.QueryUnescape(value)
	if err != nil {
		return err
	}
	response, err := http.Get(decodedValue)
	if err != nil || response.StatusCode != 404 {
		return err
	}
	filename := decodedValue[strings.LastIndex(decodedValue, "/")+1:]

	file, err := os.Create(fmt.Sprintf("%v/%v", path, filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
