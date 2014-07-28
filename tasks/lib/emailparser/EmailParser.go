package emailparser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pfacheris/kickback/models"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	NAME_SEL  = `.name`
	PRICE_SEL = `.price`
)

type EmailParser struct{}

func (e *EmailParser) Parse(ch chan *models.Purchase, r io.Reader) {

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic(err)
	}

	doc.Find(NAME_SEL).Each(func(i int, s *goquery.Selection) {
		a := s.Find("a").First()
		name := a.Text()
		productLink, _ := a.Attr("href")
		// get product id
		productID, productURL, err := getProductID(productLink)
		if err != nil {
			panic(err)
		}
		priceStr := s.SiblingsFiltered(PRICE_SEL).First().Text()
		priceSlice := strings.Trim(priceStr, " ")[1:]
		priceFloat64, err := strconv.ParseFloat(string(priceSlice), 32)
		priceFloat := float32(priceFloat64)
		if err != nil {
			panic(err)
		}

		product := models.Product{
			Id:        productID,
			Name:      name,
			URL:       productURL,
			ScrapedAt: time.Now().AddDate(0, 0, -7),
		}

		purchase := models.Purchase{
			PurchasePrice: priceFloat,
			PurchaseAt:    time.Now().Round(time.Hour * 24),
			ProductId:     product.Id,
		}

		ch <- &purchase
	})
	return
}

func getProductID(href string) (id, link string, err error) {
	parsedURL, err := url.Parse(href)
	if err != nil {
		panic(err)
	}

	innerURL, err := url.Parse(parsedURL.Query()["U"][0])
	if err != nil {
		panic(err)
	}
	path := innerURL.Path
	id = strings.Split(path, "/")[2]
	return
}
