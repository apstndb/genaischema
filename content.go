package genaischema

import (
	"context"
	"errors"
	"iter"

	"google.golang.org/genai"
)

// GenerateObjectContent returns the first candidate object from the model using T as response schema of controlled generation.
func GenerateObjectContent[T any](ctx context.Context, client *genai.Client, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (T, error) {
	next, stop := iter.Pull2(generateObjectContents[T](ctx, client, model, contents, config))
	defer stop()

	if v, err, ok := next(); ok {
		return v, err
	} else {
		return v, errors.New("empty response from model")
	}
}

// GenerateObjectContents returns the candidate objects from the model using T as response schema of controlled generation.
func GenerateObjectContents[T any](ctx context.Context, client *genai.Client, model string, contents []*genai.Content, config *genai.GenerateContentConfig) ([]T, error) {
	var results []T
	for v, err := range generateObjectContents[T](ctx, client, model, contents, config) {
		if err != nil {
			return nil, err
		}

		results = append(results, v)
	}

	return results, nil
}

func generateObjectContents[T any](ctx context.Context, client *genai.Client, model string, contents []*genai.Content, config *genai.GenerateContentConfig) iter.Seq2[T, error] {
	if config == nil {
		config = &genai.GenerateContentConfig{}
	}

	config.ResponseMIMEType = "application/json"

	if texts, err := GenerateTextContents[T](ctx, client, model, contents, config); err != nil {
		return func(yield func(T, error) bool) {
			yield(empty[T](), err)
		}
	} else {
		return func(yield func(T, error) bool) {
			for _, text := range texts {
				if !yield(unmarshal[T]([]byte(text))) {
					return
				}
			}
		}
	}
}

// GenerateEnumContent returns the first enum candidate from the model using T as response schema of controlled generation.
func GenerateEnumContent[T interface {
	~string
	Enum() []any
}](ctx context.Context, client *genai.Client, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (T, error) {
	if config == nil {
		config = &genai.GenerateContentConfig{}
	}

	config.ResponseMIMEType = "text/x.enum"

	res, err := GenerateTextContent[T](ctx, client, model, contents, config)
	return T(res), err
}

// GenerateTextContent returns the first text candidate from the model using T as response schema of controlled generation.
func GenerateTextContent[T any](ctx context.Context, client *genai.Client, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (string, error) {
	res, err := GenerateTextContents[T](ctx, client, model, contents, config)
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("empty response from model")
	}

	return res[0], nil
}

// GenerateTextContents returns the text candidates from the model using T as response schema of controlled generation.
func GenerateTextContents[T any](ctx context.Context, client *genai.Client, model string, contents []*genai.Content, config *genai.GenerateContentConfig) ([]string, error) {
	res, err := GenerateRawContents[T](ctx, client, model, contents, config)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, candidate := range res.Candidates {
		if len(candidate.Content.Parts) == 0 {
			continue
		}

		result = append(result, candidate.Content.Parts[0].Text)
	}
	return result, nil
}

// GenerateRawContents returns the raw response from the model using T as response schema of controlled generation.
func GenerateRawContents[T any](ctx context.Context, client *genai.Client, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	schema, err := ForType[T]()
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &genai.GenerateContentConfig{}
	}
	config.ResponseSchema = schema

	return client.Models.GenerateContent(ctx, model, contents, config)
}
