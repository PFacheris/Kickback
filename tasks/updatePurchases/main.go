package main

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/gmail/v1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	// App-level
	"github.com/pfacheris/kickback/config"
	. "github.com/pfacheris/kickback/db"
	"github.com/pfacheris/kickback/models"
	"github.com/pfacheris/kickback/tasks/lib/emailparser"
)

func main() {

	aWeekAgo := time.Now().AddDate(0, 0, -14)
	fmt.Println(aWeekAgo)

	oauthConfig := &oauth.Config{
		ClientId:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		RedirectURL:  config.GOOGLE_REDIRCET_URL,
		Scope:        config.GOOGLE_API_SCOPE,
		AuthURL:      config.GOOGLE_AUTH_URL,
		TokenURL:     config.GOOGLE_TOKEN_URL,
	}

	users := []models.User{}
	DB.Where("refresh_token != ''").Find(&users)

	fmt.Println(users)

	// create waitgroup
	var wg sync.WaitGroup
	out := make(chan *models.Purchase, len(users))

	wg.Add(len(users))
	for _, u := range users {
		go func(user *models.User, w *sync.WaitGroup) {
			defer w.Done()

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

			// ignore all IDs after user.LastEmailID
			if config.LAST_MESSAGE_ID != "" {
				messages = filterBeforeID(messages, config.LAST_MESSAGE_ID)
			}

			wg.Add(len(messages))

			for _, message := range messages {
				// for each message, grab its details and return if it's within the last week
				go func(messageID string) {
					defer w.Done()
					msg, err := gmailService.Users.Messages.Get(user.Email, messageID).Format("full").Do()
					if err != nil {
						return
					}

					receivedHeader, err := whereHeader(msg.Payload.Headers, func(header *gmail.MessagePartHeader) bool {
						return header.Name == "Date"
					})
					receivedDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", receivedHeader.Value)
					if err != nil {
						panic(err)
					}

					// ignore if the message was send after a week ago from now
					if !receivedDate.After(aWeekAgo) {
						return
					}

					bodyReader, err := func() (r io.Reader, err error) {
						b64body := func() string {
							fmt.Println(msg.Id)
							for _, part := range msg.Payload.Parts {
								for _, p := range part.Parts {
									headerMap := mapFromHeaders(p.Headers)
									val, ok := headerMap["Content-Type"]
									if ok && strings.Index(val, "text/html") != -1 {
										// yay
										fmt.Println(len(p.Body.Data))
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

						return strings.NewReader(""), errors.New("NO EMAIL BODY???")
					}()

					if err != nil {
						panic(err)
					}
					ep := emailparser.EmailParser{}
					ep.Parse(out, bodyReader)
					return
				}(message.Id)

			}
		}(&u, &wg)
	}
	fmt.Println("Waiting...")
	wg.Wait()
	close(out)
	fmt.Println("Done Waiting...")
	fmt.Printf("&#v", <-out)
	for purchase := range out {
		// @TODO: ADD EVERYTHING TO THE DB, YAY
		fmt.Printf("&#v", purchase)
		fmt.Println("")
	}

}

func filterBeforeID(messages []*gmail.Message, lastID string) (ret []*gmail.Message) {
	fmt.Println(messages)
	if len(messages) == 0 {
		return messages
	}
	lastIndex := 1
	for i, msg := range messages {
		if msg.Id == lastID {
			lastIndex = i
		}
	}
	ret = messages[0 : lastIndex-1]
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
