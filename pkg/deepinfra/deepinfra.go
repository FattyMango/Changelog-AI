package deepinfra

import (
	"changer/pkg/http"
	"errors"
	"strings"
)

type Deepinfra struct {
	key string
}

func NewDeepinfra(key string) (*Deepinfra, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}
	return &Deepinfra{
		key: key,
	}, nil
}

type input struct {
	Input string `json:"input"`
}

func newInput(data string) *input {
	return &input{
		Input: data,
	}
}
func (d *Deepinfra) Do(payload string) (string, error) {
	req := http.NewRequest("POST", "https://api.deepinfra.com/v1/inference/mistralai/Mixtral-8x7B-Instruct-v0.1", map[string]string{"input": payload})
	req.SetAuthorizationHeaders(d.key)
	req.SetHeaders("Content-Type", "application/json")
	resp, err := req.Execute()
	if err != nil {
		return "", err
	}

	return cleanText(resp.Results[0].GeneratedText), nil

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
