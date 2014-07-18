package main

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/gmail/v1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/pfacheris/kickback/tasks/kickbacklib"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	EMAIL           = "matthewc267@gmail.com"
	ACCESS_TOKEN    = "ya29.RgD-zFO5txK7QCEAAAC47m0nvaG5VMH8aAK3km81-ql31jhzRmSpFDBtMBK16SgkujDnwR0h96YU22YVfxo"
	REFRESH_TOKEN   = "1/2FFEwLRrPWiL5MfnTxsDCT9AfgaoEmaDvi3bbImNc5c"
	AMAZON_QUERY    = "amazon from:auto-confirm"
	LAST_MESSAGE_ID = "1450fd07eb5c4272"
)

type User struct {
	Email        string
	AccessToken  string
	RefreshToken string
}

func main() {

	aWeekAgo := time.Now().AddDate(0, 0, -14)
	fmt.Println(aWeekAgo)

	config := &oauth.Config{
		ClientId:     "692789787338-t9f5805ou1uec14gl1l4fttohtkld54e.apps.googleusercontent.com",
		ClientSecret: "mfFPXn7ZZDlXw8TR2bwtgexD",
		RedirectURL:  "http://localhost:3000/oauth2callback",
		Scope:        "https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/gmail.readonly,https://www.googleapis.com/auth/gmail.compose",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
	}

	users := make([]User, 0)
	users = append(users, User{Email: EMAIL, AccessToken: ACCESS_TOKEN, RefreshToken: REFRESH_TOKEN})
	for _, user := range users {

		transport := &oauth.Transport{
			Token:     &oauth.Token{AccessToken: user.AccessToken, RefreshToken: user.RefreshToken},
			Config:    config,
			Transport: http.DefaultTransport,
		}

		oauthHttpClient := transport.Client()

		gmailService, err := gmail.New(oauthHttpClient)
		if err != nil {
			panic(err)
		}
		// for each user, get their refresh token
		// use that token to get the access token if isExpired

		// get list of all 'amazon' messages
		listResponse, err := gmailService.Users.Messages.List(EMAIL).Q(AMAZON_QUERY).MaxResults(int64(100)).Do()
		if err != nil {
			panic(err)
		}
		messages := listResponse.Messages

		// ignore all IDs after user.LastEmailID
		if LAST_MESSAGE_ID != "" {
			messages = filterBeforeID(messages, LAST_MESSAGE_ID)
		}

		// create waitgroup
		var wg sync.WaitGroup

		out := make(chan kickback.Purchase, len(messages))
		wg.Add(len(messages))

		for _, message := range messages {
			// for each message, grab its details and return if it's within the last week
			go func(messageID string, w *sync.WaitGroup, out chan kickback.Purchase) {
				defer w.Done()
				msg, err := gmailService.Users.Messages.Get(EMAIL, messageID).Format("full").Do()
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
				ep := kickback.EmailParser{}
				ep.Parse(out, bodyReader)
				return
			}(message.Id, &wg, out)

		}
		fmt.Println("Waiting...")
		wg.Wait()
		close(out)
		fmt.Println("Done Waiting...")

		for purchase := range out {
			// @TODO: ADD EVERYTHING TO THE DB, YAY
			fmt.Printf("&#v", purchase)
			fmt.Println("")
		}
	}
}

func filterBeforeID(messages []*gmail.Message, lastID string) (ret []*gmail.Message) {
	lastIndex := 0
	for i, msg := range messages {
		if msg.Id == lastID {
			lastIndex = i
		}
	}
	ret = messages[0 : lastIndex-1]
	return ret
}

// func filterMessages(messages []*Message, fn func(*Message) bool) (ret []*Message) {
// 	ret = make([]*Message, 0)
// 	for _, msg := range messages {
// 		if fn(msg) {
// 			ret = append(ret, msg)
// 		}
// 	}
// 	return ret
// }

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
