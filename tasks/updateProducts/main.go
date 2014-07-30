package main

import (
	"encoding/xml"
	"fmt"
	"github.com/Shrugs/goaws/productadvertising/v1"
	"math"
	"sync"
	"time"

	// app
	"github.com/Shrugs/go-amazon-product-api"
	"github.com/pfacheris/kickback/config"
	. "github.com/pfacheris/kickback/db"
	"github.com/pfacheris/kickback/models"
)

func main() {

	var api amazonproduct.AmazonProductAPI

	api.AccessKey = config.AMAZON_ACCESS_KEY
	api.SecretKey = config.AMAZON_SECRET_KEY
	api.Host = "webservices.amazon.com"
	api.AssociateTag = config.AMAZON_ASSOCIATE_TAG

	// grab all of the products who's ScrapedAt date is older than a day ago
	aDayAgo := time.Now().AddDate(0, 0, -1)
	aWeekAgo := time.Now().AddDate(0, 0, -7)

	products := []models.Product{}

	// lol, fuck this ORM bullshit
	DB.Raw("SELECT products.id, products.product_id, name, u_r_l, scraped_at, products.created_at, products.updated_at, products.deleted_at FROM products JOIN purchases ON products.id=purchases.product_id WHERE purchases.purchase_at > ? AND scraped_at < ? AND kickback_amount < ?", aWeekAgo, aDayAgo, config.KICKBACK_THRESHOLD).Scan(&products)

	productIDs := make([]string, len(products))
	for _, product := range products {
		productIDs = append(productIDs, product.ProductId)
	}

	setOfProducts := splitSlice(productIDs, 10)

	var wg sync.WaitGroup
	wg.Add(len(setOfProducts))

	for _, someProducts := range setOfProducts {
		go func(someProducts []string, wg *sync.WaitGroup) {
			defer wg.Done()

			// foreach product, grab its info from the Amazon Product API
			result, err := api.MultipleItemLookup(someProducts)
			if err != nil {
				fmt.Println(err)
			}

			var response productadvertising.ItemLookupResponse
			xml.Unmarshal([]byte(result), &response)

			for _, item := range response.Items {
				var thisProduct models.Product
				DB.Where("product_id = ?", item.ASIN).First(&thisProduct)
				thisProduct.SmallImageURL = item.SmallImage.URL
				for _, offer := range item.Offers {
					// for each offer, if one of the purchases in the DB matches ProductId and SellerName, update its CurrentSellerPrice
					var purchase models.Purchase
					query := DB.Where("product_id = ?", thisProduct.Id).Where("seller_name = ?", offer.Merchant.Name).Find(&purchase)

					if query.Error != nil {
						fmt.Println("NO GOOD PURCHASES, IGNORE")
						continue
					}

					// now we have purchase, update price and kickback amount
					purchase.CurrentSellerPrice = float32(offer.OfferListing.Price.Amount) / float32(100)
					purchase.KickbackAmount = purchase.PurchasePrice - purchase.CurrentSellerPrice
					DB.Save(&purchase)
				}

			}

			// update the product price in the db
			// now for every product, get purchases
			// for each purchase
			// KickbackAmount = originalPrice - currentPrice
			// if KickbackAmount is posistive, yay, schedule an email and mark as kickbacked
			// else, do nothing

		}(someProducts, &wg)
	}
	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("Done...")
}

func splitSlice(ids []string, idsPerSplit int) (ret [][]string) {

	numIter := int(math.Ceil(float64(len(ids)) / float64(idsPerSplit)))

	ret = make([][]string, numIter)
	for i := 0; i < numIter; i++ {
		ret[i] = make([]string, idsPerSplit)
	}

	for i := 1; i <= numIter; i++ {
		ret[i-1] = ids[(i-1)*idsPerSplit : int(math.Min(float64(i*idsPerSplit), float64(len(ids))))]
	}

	return
}
