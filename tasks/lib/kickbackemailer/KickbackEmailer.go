package kickbackemailer

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/gmail/v1"
	"github.com/pfacheris/kickback/config"
	// . "github.com/pfacheris/kickback/db"
	"github.com/pfacheris/kickback/models"
	"net/http"
)

type MailgunResponse struct {
	Tags            []string      `json:"tags"`
	Timestamp       string        `json:"timestamp"`
	Envelope        Envelope      `json:"envelope"`
	RecipientDomain string        `json:"recipient-domain"`
	Event           string        `json:"event"`
	Campaigns       []string      `json:"campaigns"`
	UserVariables   UserVariables `json:"user-variables"`
	Flags           Flags         `json:"flags"`
	Routes          []Route       `Json:"routes"`
	Message         Message       `json:"message"`
	Recipient       string        `json:"recipient"`
	Method          string        `json:"method"`
}

type Envelope struct {
	Targets   string `json:"targets"`
	Transport string `json:"transport"`
	Sender    string `json:"sender"`
}

type UserVariables struct{}

type Flags struct {
	IsAuthenticated bool `json:"is-authenticated"`
	IsSystemTest    bool `json:"is-system-test"`
	IsTestMode      bool `json:"is-test-mode"`
}

type Route struct {
	Priority    int      `json:"priority"`
	Experssion  string   `json:"expression"`
	Description string   `json:"description"`
	Actions     []string `json:"actions"`
}

type Message struct {
	Headers     Headers      `json:"headers"`
	Attachments []Attachment `json:"attachments"`
	Recipients  []string     `json:"recipients"`
	Size        int          `json:"size"`
}

type Headers struct {
	To        string `json:"to"`
	MessageId string `json:"message-id"`
	From      string `json:"from"`
	Subject   string `json:"subject"`
}

type Attachment struct{}

func SendGmailEmailFromUser(user *models.User, purchases *[]models.Purchase) {

	oauthConfig := &oauth.Config{
		ClientId:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		RedirectURL:  config.GOOGLE_REDIRCET_URL,
		Scope:        config.GOOGLE_API_SCOPE,
		AuthURL:      config.GOOGLE_AUTH_URL,
		TokenURL:     config.GOOGLE_TOKEN_URL,
	}

	transport := &oauth.Transport{
		Token:     &oauth.Token{RefreshToken: user.RefreshToken},
		Config:    oauthConfig,
		Transport: http.DefaultTransport,
	}

	err := transport.Refresh()
	if err != nil {
		panic(err)
	}

	oauthHttpClient := transport.Client()

	gmailService, err := gmail.New(oauthHttpClient)
	if err != nil {
		panic(err)
	}

	message := &gmail.Message{}
	gmailService.Users.Messages.Send(user.Email, message)
}
