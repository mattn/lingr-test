package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var nickname = flag.String("nickname", "mattn", "nickname on lingr")
var speaker_id = flag.String("speakerid", "mattn", "speaker_id on lingr")
var room = flag.String("room", "vim", "room on lingr")
var ptype = flag.String("type", "human", "type of poster")
var icon_url = flag.String("iconurl", "http://mattn.kaoriya.net/images/logo.png", "icon_url of poster")

type Event struct {
	Id      int      `json:"event_id"`
	Message *Message `json:"message"`
}

type Message struct {
	Id              string `json:"id"`
	Room            string `json:"room"`
	PublicSessionId string `json:"public_session_id"`
	IconUrl         string `json:"icon_url"`
	Type            string `json:"type"`
	SpeakerId       string `json:"speaker_id"`
	Nickname        string `json:"nickname"`
	Text            string `json:"text"`
	Timestamp       string `json:"timestamp"`
	Mine            bool   `json:"mine"`
}

var reUrl = regexp.MustCompile(`(?:^|[^a-zA-Z0-9])(https?://[a-zA-Z][a-zA-Z0-9_-]*(\.[a-zA-Z0-9][a-zA-Z0-9_-]*)*(:\d+)?(?:/[a-zA-Z0-9_/.\-+%#?&=;@$,!*~]*)?)`)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: lingr-test [options] [bot_id|URL] [text]")
		flag.PrintDefaults()
		return
	}

	uri := flag.Arg(0)
	if !reUrl.MatchString(uri) {
		uri = fmt.Sprintf("http://lingr.com/bot/%s", flag.Arg(0))
		doc, err := goquery.NewDocument(uri)
		if err != nil {
			log.Fatal(err)
		}
		uri = ""
		doc.Find("#property .left").Each(func(_ int, s *goquery.Selection) {
			if strings.TrimSpace(s.Text()) == "Endpoint:" {
				uri = strings.TrimSpace(s.Next().Text())
			}
		})
		if uri == "" {
			log.Fatal("404 Bot Not Found")
		}
	}

	events := map[string][]Event{
		"events": {
			{
				Id: 12345,
				Message: &Message{
					Id:              "",
					Room:            *room,
					PublicSessionId: "VIM",
					IconUrl:         *icon_url,
					Type:            *ptype,
					SpeakerId:       *speaker_id,
					Nickname:        *nickname,
					Text:            strings.Join(flag.Args()[1:], " "),
					Timestamp:       time.Now().Format(time.RFC822),
				},
			},
		},
	}
	b, err := json.Marshal(&events)
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.Post(uri, "application/json", strings.NewReader(string(b)))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	io.Copy(os.Stdout, r.Body)
}
