package worker

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"

	"crawler/pkg/lib"
)

// Page structure
type Page struct {
	Link  string
	Title string
}

// Worker pool
type Worker struct {
	urlQueue    chan *url.URL
	maxDispatch int
	pageChan    chan Page
}

// NewWorker constructor
func NewWorker(maxDispatch int) Worker {
	return Worker{
		urlQueue:    make(chan *url.URL),
		maxDispatch: maxDispatch,
		pageChan:    make(chan Page),
	}
}

// Run the loop
func (w *Worker) Run() {
	for i := 0; i < w.maxDispatch; i++ {
		go func() {

			for url := range w.urlQueue {
				err := parseHtml(url, w.pageChan)
				if err != nil {
					log.Printf("Critical error parse Html: %s", err.Error())
				}
			}
		}()
	}
}

// Add the link
func (w *Worker) Add(link string) error {

	lnk, err := url.Parse(link)
	if !(err == nil && lnk.Scheme != "" && lnk.Host != "") {
		return errors.New("Link not valid")
	}

	go func() {
		w.urlQueue <- lnk
	}()

	return nil
}

// PrintResult method
func (w *Worker) PrintResult() {
	go func() {
		for page := range w.pageChan {
			log.Printf("Href: %s, Title: %s", page.Link, page.Title)
		}
	}()
}

// TODO
// func (w *Worker) Close()

func parseHtml(url *url.URL, pageChan chan Page) error {

	response, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return err
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}

	lnks, err := lib.GetElementsNode(doc, "a")
	if err != nil {
		return err
	}

	var href, title string
	for _, lnk := range lnks {

		hrefLink, err := lib.GetElementNodeAttrValueString(&lnk, "href")
		if err != nil {
			log.Printf("error get href %s: %s", lnk, err.Error())
			continue
		}

		urlLnk, err := url.Parse(hrefLink)
		if err != nil {
			log.Printf("error parse href %s: %s", lnk, err.Error())
			continue
		}

		href = urlLnk.String()

		if urlLnk.Host != url.Host {
			continue
		}

		{
			response, err := http.Get(href)
			if err != nil {
				log.Printf("error get html by url %s: %s", href, err.Error())
				continue
			}
			defer response.Body.Close()

			subDoc, err := html.Parse(response.Body)
			if err != nil {
				log.Printf("error parse html by url %s: %s", href, err.Error())
				continue
			}
			err = response.Body.Close()
			if err != nil {
				log.Printf("error clode body %s: %s", href, err.Error())
				continue
			}

			headNode, err := lib.GetElementNode(subDoc, "head")
			if err != nil {
				log.Printf("error get head by url %s: %s", href, err.Error())
				continue
			}

			titleNode, err := lib.GetElementNode(&headNode, "title")
			if err != nil {
				log.Printf("error get title by url %s: %s", href, err.Error())
				continue
			}

			if titleNode.FirstChild != nil {
				title = titleNode.FirstChild.Data
			}
		}

		pageChan <- Page{
			Link:  href,
			Title: title,
		}
	}

	return nil
}
