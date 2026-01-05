package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaProvider struct {
	BaseURL string
	Model   string
}

func NewOllamaProvider() *OllamaProvider {
	return &OllamaProvider{

		BaseURL: "http://localhost:11434/api/embeddings",
		Model:   "nomic-embed-text",
	}
}

func (o *OllamaProvider) GetVector(text string) ([]float32, error) {
	payload := map[string]string{
		"model":  o.Model,
		"prompt": text,
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(o.BaseURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama error: %d", resp.StatusCode)
	}

	var res struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Embedding, nil
}
