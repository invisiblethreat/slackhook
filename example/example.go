package main

import (
	"os"
	"time"

	"github.com/invisiblethreat/slackhook"
)

func main() {
	hook := slackhook.NewHook(os.Getenv("SLACK_HOOK"))
	msg := slackhook.Message{
		Channel:  "#hook-testing",
		Text:     "testing text",
		UserName: "slackhook tester",
		IconURL:  "https://slack.global.ssl.fastly.net/66f9/img/avatars/ava_0002-48.png",
	}

	att := slackhook.Attachment{
		Color:      "good",
		Text:       "text",
		Fallback:   "fallback",
		AuthorName: "author",
		AuthorLink: "https://github.com/invisiblethreat/slackhook",
		AuthorIcon: "https://avatars1.githubusercontent.com/u/2525006?s=40&v=4",
		Title:      "attachment title",
		TitleLink:  "https://github.com/invisiblethreat/slackhook",
		ImageURL:   "http://i.imgur.com/50NA7vr.gif",
		Footer:     "footer text",
		FooterIcon: "http://simpleicon.com/wp-content/uploads/foot.png",
	}

	att.TSSet(time.Now())

	field := slackhook.Field{
		Title: "field title",
		Value: "field value",
		Short: false,
	}

	att.AddField(field)

	msg.Attach(&att)
	hook.Send(&msg)
}
