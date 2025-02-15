package genaischema

import (
	"context"
	"encoding/json"
	"errors"

	"google.golang.org/genai"
)

func GenerateObjectContent[T any](ctx context.Context, client *genai.Client, contents []*genai.Content) (T, error) {
	var zero T

	res, err := GenerateTextContent[T](ctx, client, contents, "application/json")
	if err != nil {
		return zero, err
	}

	var result T
	if err = json.Unmarshal([]byte(res), &result); err != nil {
		return zero, err
	}

	return result, nil
}

func GenerateEnumContent[T interface {
	~string
	Enum() []any
}](ctx context.Context, client *genai.Client, contents []*genai.Content) (T, error) {
	res, err := GenerateTextContent[T](ctx, client, contents, "text/x.enum")
	if err != nil {
		return "", err
	}

	return T(res), nil
}

func GenerateTextContent[T any](ctx context.Context, client *genai.Client, contents []*genai.Content, mimeType string) (string, error) {
	schema, err := ForType[T]()
	if err != nil {
		return "", err
	}

	res, err := client.Models.GenerateContent(ctx,
		"gemini-2.0-flash",
		contents,
		&genai.GenerateContentConfig{
			ResponseMIMEType: mimeType,
			ResponseSchema:   schema,
		})
	if err != nil {
		return "", err
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}

	return res.Candidates[0].Content.Parts[0].Text, nil
}
