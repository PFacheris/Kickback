package main

import (
	"fmt"
	"github.com/pfacheris/kickback/tasks/kickbacklib"
	"sync"
	// "github.com/Shrugs/goaws"
	"time"
)

func main() {

	// grab all of the products who's ScrapedAt date is older than a day ago
	aDayAgo := time.Now().AddDate(0, 0, -1)

	wg := sync.WaitGroup()
	wg.Add(len(products))

	for _, product := range products {
		go func(product *kickback.Product, wg *sync.WaitGroup) {
			defer wg.Done()

			// foreach product, grab its info from the Amazon Product API

		}(product, &wg)
	}
}
