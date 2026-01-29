package main

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type geminiProvider struct{}

func init() {
	registerProvider(geminiProvider{})
}

func (geminiProvider) Name() string { return "gemini" }

func (geminiProvider) ResolveModel(alias string) string {
	aliases := map[string]string{
		"flash":      "gemini-2.5-flash",
		"pro":        "gemini-2.5-pro",
		"flash-lite": "gemini-2.0-flash-lite",
	}
	if id, ok := aliases[alias]; ok {
		return id
	}
	return alias
}

func (geminiProvider) ModelAliases() []string {
	return []string{"flash", "pro", "flash-lite"}
}

func (geminiProvider) DefaultModel() string { return "flash" }

func (geminiProvider) EnvKey() string { return "GEMINI_API_KEY" }

func (p geminiProvider) Run(ctx context.Context, prompt, model, apiKey, _ string) error {
	modelID := p.ResolveModel(model)
	if modelID == "" {
		modelID = p.ResolveModel(p.DefaultModel())
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return runStreaming(func(emit func(string)) error {
		for result, err := range client.Models.GenerateContentStream(
			ctx,
			modelID,
			genai.Text(prompt),
			nil,
		) {
			if err != nil {
				return fmt.Errorf("Gemini API error: %v", err)
			}
			emit(result.Text())
		}
		return nil
	})
}
