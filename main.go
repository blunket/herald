package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/nlopes/slack"
)

const tpl = `Hey, everybody. Have you seen the <#{{ .ID }}|{{ .Name }}> channel recently?

{{ if .Purpose.Value -}}
Purpose: {{ .Purpose.Value }}
{{ else }}
Purpose: _(not set)_
{{ end }}`

var t *template.Template

func init() {
	t = template.Must(template.New("msg").Parse(tpl))
	rand.Seed(time.Now().UnixNano())
}

func main() {
	tkn := os.Getenv("SLACK_TOKEN")
	if tkn == "" {
		log.Fatal("Env var SLACK_TOKEN was not set")
	}

	api := slack.New(tkn)

	c, err := getRandomChannel(api)
	if err != nil {
		log.Print(err)
	}

	if err := announceChannel(api, c); err != nil {
		log.Print(err)
	}
}

func announceChannel(api *slack.Client, c slack.Channel) error {
	var b bytes.Buffer
	if err := t.Execute(&b, c); err != nil {
		return err
	}

	log.Printf("Sending message\n%s", b.String())

	params := slack.NewPostMessageParameters()
	params.AsUser = true
	params.EscapeText = false
	_, _, err := api.PostMessage("herald-test", b.String(), params)

	return err
}

func getRandomChannel(api *slack.Client) (slack.Channel, error) {
	ch, err := api.GetChannels(true)
	if err != nil {
		return slack.Channel{}, err
	}

	return ch[rand.Intn(len(ch))], nil
}
