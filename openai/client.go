package openai

import (
	"context"
	"errors"
	"math/rand"
	"time"

	vendor "github.com/sashabaranov/go-openai"
)

type ModerationModel string

var singleTokenWords = []string{"Hello", "Dark", "My", "Old", "Friend"}

const (
	TextModerationLatest ModerationModel = "text-moderation-latest"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Client is a wrapper around the go-openai client
type Client struct {
	*vendor.Client
}

// NewClient creates a new client
func NewClient(apiKey string) *Client {
	return &Client{
		vendor.NewClient(apiKey),
	}
}

func (c *Client) Moderate(text string) ([]vendor.Result, error) {
	model := string(TextModerationLatest)
	response, err := c.Client.Moderations(context.Background(), vendor.ModerationRequest{
		Input: text,
		Model: &model,
	})
	if err != nil {
		return nil, err
	}
	return response.Results, nil
}

func (c *Client) TestConnection() error {
	testWord := singleTokenWords[rand.Intn(len(singleTokenWords))]
	results, err := c.Moderate(testWord)
	if err != nil {
		return err
	}
	for _, result := range results {
		if result.Flagged {
			return errors.New("OpenAI inappropriately flagged the text '" + testWord + "' as inappropriate")
		}
	}
	return nil
}
