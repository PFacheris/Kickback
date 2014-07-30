package emailparser

import (
	"fmt"
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

func (e *EmailParser) Parse(ch chan *models.PurchaseData, r io.Reader) {

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

		sellerA := a.NextAllFiltered("a").First()
		fmt.Println(sellerA)
		sellerName := sellerA.Text()
		if sellerName == "" {
			// @TODO(Shrugs) make sure that this is changed based on region and stuff
			sellerName = "Amazon.com"
		}

		priceStr := s.SiblingsFiltered(PRICE_SEL).First().Text()
		priceSlice := strings.Trim(priceStr, " ")[1:]
		priceFloat64, err := strconv.ParseFloat(string(priceSlice), 32)
		priceFloat := float32(priceFloat64)
		if err != nil {
			panic(err)
		}

		purchaseData := models.PurchaseData{
			PurchasePrice: priceFloat,
			// @TODO(Shrugs) remove the AddDate
			PurchaseAt:  time.Now().AddDate(0, 0, -3).Round(time.Hour * 24),
			ProductId:   productID,
			ProductName: name,
			ProductURL:  productURL,
			SellerName:  sellerName,
		}

		ch <- &purchaseData
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
	link = innerURL.String()
	path := innerURL.Path
	id = strings.Split(path, "/")[2]
	return
}

// func getSellerID(href string) (id string, err error) {
// 	parsedURL, err := url.Parse(href)
// 	if err != nil {
// 		panic(err)
// 	}

// 	innerURL, err := url.Parse(parsedURL.Query()["U"][0])
// 	if err != nil {
// 		panic(err)
// 	}
// 	link := innerURL.String()
// 	parsedLink, err := url.Parse(link)
// 	if err != nil {
// 		return
// 	}
// 	id = parsedLink.Query()["seller"][0]
// 	return
// }
