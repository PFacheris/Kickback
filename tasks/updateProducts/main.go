package main

import (
	"fmt"
	"github.com/DDRBoxman/go-amazon-product-api"
	"github.com/pfacheris/kickback/config"
	"github.com/pfacheris/kickback/tasks/lib/emailparser"
	"math"
	"sync"
	"time"
)

func main() {

	var api amazonproduct.AmazonProductAPI

	api.AccessKey = config.AMAZON_ACCESS_KEY
	api.SecretKey = config.AMAZON_SECRET_KEY
	api.Host = "webservices.amazon.com"
	api.AssociateTag = config.AMAZON_ASSOCIATE_TAG

	// grab all of the products who's ScrapedAt date is older than a day ago
	aDayAgo := time.Now().AddDate(0, 0, -1)

	wg := sync.WaitGroup()
	wg.Add(len(products))

	for _, product := range products {
		go func(product *kickback.Product, wg *sync.WaitGroup) {
			defer wg.Done()

			// foreach product, grab its info from the Amazon Product API
			result, err := api.ItemSearchByKeyword("sgt+frog")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(result)

		}(product, &wg)
	}

	wg.Wait()
}

func splitSlice(ids []string, idsPerSplit int) (ret [][]string) {

	numIter := math.Ceil(len(ids) / idsPerSplit)

	// @TODO fix this
	ret = make([][]string)

	for i := 1; i <= numIter; i++ {
		ret[i-1] = ids[(i-1)*idsPerSplit : i*idsPerSplit]
	}

	// for (var i = 0; i < numIter; i++) {
	//     // take the first <itemsPerRow> and add to a row
	//     r += rowTemplate({projects: ps.splice(0, itemsPerRow)});
	// }
}
