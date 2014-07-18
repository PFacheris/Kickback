package kickback

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestEmailParser(t *testing.T) {
	body, err := ioutil.ReadFile("email.html")
	if err != nil {
		panic(err)
	}

	e := &EmailParser{}
	ch := make(chan []Purchase)
	go e.Parse(ch, strings.NewReader(string(body)))

	fmt.Println(<-ch)

}
