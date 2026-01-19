package gemini

import (
	"context"
	"os"

	"google.golang.org/genai"
)

// NewClient cria um novo cliente para interagir com a API do Gemini.
func NewClient(ctx context.Context) (*genai.Client, error) {
	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
}

// GenerateText gera texto a partir de um prompt de texto.
func GenerateText(client *genai.Client, ctx context.Context, prompt string) (string, error) {
	model := genai.NewGenerativeModel(client, "gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var content string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					content += string(txt)
				}
			}
		}
	}
	return content, nil
}
