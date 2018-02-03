package slackhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Message to send to Slack's Incoming WebHook API.
//
// See https://api.slack.com/incoming-webhooks
type Message struct {
	Text        string        `json:"text"`
	Channel     string        `json:"channel,omitempty"`
	UserName    string        `json:"username,omitempty"`
	IconURL     string        `json:"icon_url,omitempty"`
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

// Attachment is a type to provide richly formatted messages
// See https://api.slack.com/docs/attachments
type Attachment struct {
	Fallback   string  `json:"fallback,omitempty"` // plain text summary
	Color      string  `json:"color,omitempty"`    // {good|warning|danger|hex}
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Title      string  `json:"title,omitempty"` // larger, bold text at top of attachment
	TitleLink  string  `json:"title_link,omitempty"`
	Text       string  `json:"text,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
	ThumbURL   string  `json:"thumb_url,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	Timestamp  int     `json:"ts,omitempty"` // Unix timestamp
}

// Field is the small footer in a Slack message
type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// Attach an attachment to a Slack message.
func (m *Message) Attach(a *Attachment) {
	m.Attachments = append(m.Attachments, a)
}

// TSSet allow for setting of an arbitrary timestamp
func (a *Attachment) TSSet(t time.Time) {
	a.Timestamp = int(t.Unix())
}

// TSNow sets the attachment timestamp with the current time.
func (a *Attachment) TSNow() {
	a.Timestamp = int(time.Now().Unix())
}

// AddField attaches a new Field strcut to an attachemnt
func (a *Attachment) AddField(f Field) {
	a.Fields = append(a.Fields, f)
}

// Poster interface is the methods of http.Client required by Client to ease
// testing.
type Poster interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

// Client for Slack's Incoming WebHook API.
type Client struct {
	url        string
	HTTPClient Poster
}

// NewHook Slack Incoming WebHook Client using http.DefaultClient for its Poster.
func NewHook(url string) *Client {
	return &Client{url: url, HTTPClient: http.DefaultClient}
}

// NewAttachmentGood is a courtesy feature for a good attachment
func NewAttachmentGood() *Attachment {
	return &Attachment{Color: "good"}
}

// NewAttachmentWarning is a courtesy feature for a warning attachment
func NewAttachmentWarning() *Attachment {
	return &Attachment{Color: "warning"}
}

// NewAttachmentDanger is a courtesy feature for a danger attachment
func NewAttachmentDanger() *Attachment {
	return &Attachment{Color: "danger"}
}

// Simple text message.
func (c *Client) Simple(msg string) error {
	return c.Send(&Message{Text: msg})
}

// Send a Message.
func (c *Client) Send(msg *Message) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Post(c.url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Discard response body to reuse connection
	io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
