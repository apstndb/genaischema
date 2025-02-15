package genaischema

import (
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"github.com/swaggest/jsonschema-go"
	"google.golang.org/genai"
)

func GenerateForValue(value any) (*genai.Schema, error) {
	reflector := jsonschema.Reflector{}
	schema, err := reflector.Reflect(value, jsonschema.InlineRefs)
	if err != nil {
		return nil, err
	}

	return convertSchema(schema)
}

func GenerateForType[T any]() (*genai.Schema, error) {
	var v T
	return GenerateForValue(v)
}

func convertSchema(schema jsonschema.Schema) (*genai.Schema, error) {
	var err error

	var anyOf []*genai.Schema
	for _, elem := range schema.AnyOf {
		s, err := convertSchemaOrBool(elem)
		if err != nil {
			return nil, err
		}
		anyOf = append(anyOf, s)
	}

	properties := make(map[string]*genai.Schema)
	for k, v := range schema.Properties {
		s, err := convertSchemaOrBool(v)
		if err != nil {
			return nil, err
		}
		properties[k] = s
	}

	var typ genai.Type
	var items *genai.Schema
	if schema.Items != nil {
		typ = genai.TypeArray
		if schema.Items.SchemaOrBool == nil {
			return nil, fmt.Errorf("type is array, but schema.Items.SchemaOrBool is nil")
		}

		items, err = convertSchemaOrBool(*schema.Items.SchemaOrBool)
		if err != nil {
			return nil, err
		}
	} else {
		typ, err = convertType(schema.Type)
		if err != nil {
			return nil, err
		}
	}

	enum, err := convertEnum(schema.Enum)
	if err != nil {
		return nil, err
	}

	var def any
	if schema.Default != nil {
		def = schema.Default
	}

	var example any
	if len(schema.Examples) > 0 {
		example = schema.Examples
	}

	return &genai.Schema{
		MinItems:         lo.EmptyableToPtr(schema.MinItems),
		Example:          example,
		PropertyOrdering: nil, // TODO
		Pattern:          lo.FromPtr(schema.Pattern),
		Minimum:          schema.Minimum,
		Default:          def,
		AnyOf:            anyOf,
		MaxLength:        schema.MaxLength,
		Title:            lo.FromPtr(schema.Title),
		MinLength:        lo.EmptyableToPtr(schema.MinLength),
		MinProperties:    lo.EmptyableToPtr(schema.MinProperties),
		MaxItems:         schema.MaxItems,
		Maximum:          schema.Maximum,
		Nullable:         false, // TODO
		MaxProperties:    schema.MaxProperties,
		Type:             typ,
		Description:      lo.FromPtr(schema.Description),
		Enum:             enum,
		Format:           lo.FromPtr(schema.Format),
		Items:            items,
		Properties:       properties,
		Required:         schema.Required,
	}, nil
}

func convertEnum(enum []interface{}) ([]string, error) {
	var result []string
	b, err := json.Marshal(enum)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func convertType(t *jsonschema.Type) (genai.Type, error) {
	if t == nil {
		return "", fmt.Errorf("invalid argument: type is nil")
	}

	switch lo.FromPtr(t.SimpleTypes) {
	case jsonschema.Array:
		return genai.TypeArray, nil
	case jsonschema.String:
		return genai.TypeString, nil
	case jsonschema.Number:
		return genai.TypeNumber, nil
	case jsonschema.Integer:
		return genai.TypeInteger, nil
	case jsonschema.Boolean:
		return genai.TypeBoolean, nil
	case jsonschema.Object:
		return genai.TypeObject, nil
	default:
		// TODO
		return genai.TypeObject, nil
	}
}

func convertSchemaOrBool(schema jsonschema.SchemaOrBool) (*genai.Schema, error) {
	switch {
	case schema.TypeObject != nil:
		return convertSchema(*schema.TypeObject)
	default:
		return nil, fmt.Errorf("jsonschema.SchemaOrBool.TypeObject is only supported")
	}
}
