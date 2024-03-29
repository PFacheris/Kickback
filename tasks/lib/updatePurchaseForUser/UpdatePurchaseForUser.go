package updatePurchaseForUser

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/gmail/v1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"github.com/pfacheris/kickback/config"
	"github.com/pfacheris/kickback/tasks/lib/emailparser"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	// App
	. "github.com/pfacheris/kickback/db"
	"github.com/pfacheris/kickback/models"
)

func Do(user *models.User) {

	aWeekAgo := time.Now().AddDate(0, -1, 0)
	oauthConfig := &oauth.Config{
		ClientId:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		RedirectURL:  config.GOOGLE_REDIRCET_URL,
		Scope:        config.GOOGLE_API_SCOPE,
		AuthURL:      config.GOOGLE_AUTH_URL,
		TokenURL:     config.GOOGLE_TOKEN_URL,
	}

	// create waitgroup
	var wg sync.WaitGroup
	out := make(chan *models.PurchaseData, 1000)

	transport := &oauth.Transport{
		Token:     &oauth.Token{RefreshToken: user.RefreshToken},
		Config:    oauthConfig,
		Transport: http.DefaultTransport,
	}

	err := transport.Refresh()
	if err != nil {
		panic(err)
	}

	// fmt.Println("%#v", transport.Token)

	oauthHttpClient := transport.Client()

	gmailService, err := gmail.New(oauthHttpClient)
	if err != nil {
		panic(err)
	}

	// get list of all 'amazon' messages
	listResponse, err := gmailService.Users.Messages.List(user.Email).Q(config.AMAZON_QUERY).MaxResults(int64(100)).Do()
	if err != nil {
		panic(err)
	}
	messages := listResponse.Messages
	// fmt.Println(messages)
	// ignore all IDs after user.LastMessageId
	// fmt.Println(user.LastMessageId)
	if user.LastMessageId != "" {
		messages = filterBeforeID(messages, user.LastMessageId)
	}

	mostRecentMessageID := user.LastMessageId
	// make a time that's definitely last
	mostRecentMessageTime := time.Now().AddDate(-1, 0, 0)

	wg.Add(len(messages))
	for _, message := range messages {
		// for each message, grab its details and return if it's within the last week
		go func(messageID string) {
			defer wg.Done()

			msg, err := gmailService.Users.Messages.Get(user.Email, messageID).Format("full").Do()
			if err != nil {
				panic(err)
				return
			}

			receivedHeader, err := whereHeader(msg.Payload.Headers, func(header *gmail.MessagePartHeader) bool {
				return header.Name == "Date"
			})
			receivedDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", receivedHeader.Value)
			if err != nil {
				// fmt.Println(receivedHeader.Value)
				fmt.Println(err)
				// panic(err)
			}

			// ignore if the message was send after a week ago from now
			if receivedDate.Before(aWeekAgo) {
				// fmt.Println("IGNORING MESSAGE")
				return
			}
			if receivedDate.After(mostRecentMessageTime) {
				mostRecentMessageTime = receivedDate
				mostRecentMessageID = messageID
			}

			bodyReader, err := func() (r io.Reader, err error) {
				b64body := func() string {
					// fmt.Println(msg.Id)
					for _, part := range msg.Payload.Parts {
						for _, p := range part.Parts {
							headerMap := mapFromHeaders(p.Headers)
							val, ok := headerMap["Content-Type"]
							if ok &&
								strings.Index(val, "text/html") != -1 &&
								len(p.Body.Data) > 0 {
								// yay
								return p.Body.Data
							}
						}
					}
					return ""
				}()

				if len(b64body) > 0 {
					r, _ := base64.URLEncoding.DecodeString(b64body)
					return strings.NewReader(string(r[:])), err
				}
				pretty.Println(msg)
				return strings.NewReader(""), errors.New("NO EMAIL BODY???")
			}()

			if err != nil {
				// panic(err)
				return
			}
			ep := emailparser.EmailParser{}
			ep.Parse(out, bodyReader)
			return
		}(message.Id)

	}

	fmt.Println("Waiting")
	wg.Wait()
	fmt.Println("Done Waiting...")
	close(out)
	for purchaseData := range out {

		// grab the Product from the DB or create a new one
		product := &models.Product{}
		if err := product.Get(purchaseData.ProductId); err != nil {
			product = &models.Product{
				ProductId: purchaseData.ProductId,
				Name:      purchaseData.ProductName,
				URL:       purchaseData.ProductURL,
				ScrapedAt: purchaseData.PurchaseAt,
			}
		}
		if DB.NewRecord(product) {
			DB.Create(product)
		}

		purchase := models.Purchase{
			PurchasePrice:      purchaseData.PurchasePrice,
			KickbackAmount:     0.0,
			PurchaseAt:         purchaseData.PurchaseAt,
			UserId:             user.Id,
			ProductId:          product.Id,
			SellerName:         purchaseData.SellerName,
			CurrentSellerPrice: purchaseData.PurchasePrice,
			WasKickbacked:      false,
		}

		product.Purchases = append(product.Purchases, purchase)
		DB.Save(product)
	}

	// update the user
	user.LastMessageId = mostRecentMessageID
	DB.Save(user)
	return

}

func filterBeforeID(messages []*gmail.Message, lastID string) (ret []*gmail.Message) {
	// fmt.Println(messages)
	if len(messages) == 0 {
		return messages
	}
	lastIndex := len(messages)
	for i, msg := range messages {
		if msg.Id == lastID {
			lastIndex = i
			break
		}
	}
	ret = messages[0:lastIndex]
	return ret
}

func whereHeader(headers []*gmail.MessagePartHeader, fn func(*gmail.MessagePartHeader) bool) (h *gmail.MessagePartHeader, err error) {
	for _, header := range headers {
		if fn(header) {
			return header, err
		}
	}
	return nil, errors.New("NO HEADER FOUND")
}

func mapFromHeaders(headers []*gmail.MessagePartHeader) (ret map[string]string) {
	ret = make(map[string]string)
	for _, header := range headers {
		ret[header.Name] = header.Value
	}
	return
}
