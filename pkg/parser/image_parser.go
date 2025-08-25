package parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

type ImageParser struct {
	Path              string
	SelectionCount    int
	MaxImageBytesSize int
	browserCtx        context.Context
}

func New(path string, selectionCount, maxImageBytesSize int) *ImageParser {
	if selectionCount == 0 {
		selectionCount = 10
	}
	if maxImageBytesSize == 0 {
		maxImageBytesSize = 1 * 1024 * 1024
	}

	return &ImageParser{
		Path:              path,
		SelectionCount:    selectionCount,
		MaxImageBytesSize: maxImageBytesSize,
	}
}

func (i *ImageParser) runBrowser() func() {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true))

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	i.browserCtx = browserCtx

	return func() {
		cancel()
		cancelBrowser()
		i.browserCtx = nil
		println("Browser closed")
	}
}

func (i *ImageParser) Parse(queries []string) {

	cleanUp := i.runBrowser()
	defer cleanUp()

	wg := sync.WaitGroup{}
	errorsChan := make(chan error, len(queries))
	queriesChan := make(chan string, len(queries))

	for idx := 0; idx < runtime.NumCPU(); idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i.parse(queriesChan, errorsChan)
		}()
	}

	for _, query := range queries {
		queriesChan <- query
	}
	close(queriesChan)

	wg.Wait()
	close(errorsChan)

	for err := range errorsChan {
		if err != nil {
			log.Println("Error:", err)
		}
	}
}

func (i *ImageParser) parse(queryChan <-chan string, errChan chan<- error) {
	ctx, cancel := chromedp.NewContext(i.browserCtx)
	defer cancel()

	var html string

	for query := range queryChan {

		fmt.Println("🔍 Картинки по запросу:", query)
		err := chromedp.Run(ctx,
			chromedp.Navigate("https://www.google.com/search?tbm=isch&q="+query),
			chromedp.OuterHTML("html", &html))

		if err != nil {
			errChan <- err
			return
		}

		images, err := i.findImagesInHtml(html)

		if err != nil {
			errChan <- err
			return
		}
		image, err := i.findOptimalImage(images)
		if err != nil {
			errChan <- err
			return
		}
		err = i.saveImage(image, query)
		if err != nil {
			errChan <- err
			return
		}
	}

}

func (i *ImageParser) findImagesInHtml(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return nil, err
	}

	found := []string{}

	doc.Find("img").Each(func(idx int, s *goquery.Selection) {
		if len(found) >= i.SelectionCount {
			return
		}

		src, _ := s.Attr("src")
		if strings.Contains(src, "jpg") || strings.Contains(src, "jpeg") {
			found = append(found, src)
		}
	})

	return found, nil
}

func (i *ImageParser) findOptimalImage(images []string) ([]byte, error) {
	var largestImage []byte
	var largestSize int

	for _, img := range images {
		// Убираем префикс, если есть
		img = strings.TrimPrefix(img, "data:image/jpeg;base64,")
		img = strings.TrimPrefix(img, "data:image/jpg;base64,")

		// Декодируем Base64
		data, err := base64.StdEncoding.DecodeString(img)
		if err != nil {
			return []byte(""), err
		}

		size := len(data)

		if size > largestSize && size <= i.MaxImageBytesSize {
			largestImage = data
		}
	}

	return largestImage, nil
}

func (i *ImageParser) saveImage(img []byte, imageName string) error {
	println("Saving", imageName)
	err := os.MkdirAll(i.Path, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(i.Path+imageName+".jpeg", img, 0644)
	if err != nil {
		return err
	}

	return nil
}
