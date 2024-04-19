package gemeni

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gemeni struct {
	context context.Context
	key     string
	client  *genai.Client
}

func NewGemini(ctx context.Context, key string) (*Gemeni, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}
	c, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	return &Gemeni{
		key:     key,
		client:  c,
		context: ctx,
	}, nil
}
func (g *Gemeni) Close() error {
	return g.client.Close()
}
func (g *Gemeni) Do(payload string) (string, error) {

	model := g.client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(g.context, genai.Text(payload))
	if err != nil {
		return "", err
	}

	s, err := json.Marshal(resp.Candidates[0].Content.Parts[0])
	if err != nil {
		return "", err
	}

	return string(s), nil

}

func cleanText(text string) string {
	// search for any consecutive spaces and replace them with a single space
	text = strings.Join(strings.Fields(text), " ")
	text = strings.ReplaceAll(text, "\n\n", "")
	text = strings.ReplaceAll(text, "\n\t", "")
	text = strings.ReplaceAll(text, "\t\t", "")
	text = strings.ReplaceAll(text, "\t", "")
	// split by - and join with new line
	text = strings.Join(strings.Split(text, "-"), "\n")
	return text
}
