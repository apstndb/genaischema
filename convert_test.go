package genaischema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/genai"
)

func TestGenerateForValue(t *testing.T) {
	tests := []struct {
		desc  string
		value any
		want  *genai.Schema
	}{
		{
			"simple",
			struct {
				Text    string  `json:"text"`
				Number  float64 `json:"number"`
				Integer int     `json:"integer"`
			}{},
			&genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"text":    {Type: genai.TypeString},
					"number":  {Type: genai.TypeNumber},
					"integer": {Type: genai.TypeInteger},
				}},
		},
		{
			"nested",
			struct {
				Inner struct {
					Text string `json:"text"`
				} `json:"inner"`
				Array []struct {
					Text string `json:"text"`
				} `json:"array"`
			}{},
			&genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"inner": {Type: genai.TypeObject, Properties: map[string]*genai.Schema{
						"text": {Type: genai.TypeString},
					}},
					"array": {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"text": {Type: genai.TypeString},
						}}},
				},
			},
		},
		{
			"complex",
			struct {
				Text          string   `json:"text" description:"text field" minLength:"1" maxLength:"100"`
				Direction     string   `json:"direction" title:"Direction" description:"Direction of target" enum:"NORTH,SOUTH,EAST,WEST" required:"true"`
				Email         string   `json:"email" format:"email" required:"true"`
				Number        float64  `json:"number" maximum:"100.0"`
				Integer       int      `json:"integer" minimum:"1" default:"-1"`
				Abc           string   `json:"abc" pattern:"^[abc]$"`
				ArrayOfString []string `json:"arrayOfString" minItems:"1" maxItems:"10"`
				_             any      `minProperties:"3" maxProperties:"4"`
				_             any      `title:"Complex" description:"Example of complex schema generation"`
			}{},
			&genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"text":      {Type: genai.TypeString, Description: "text field", MinLength: genai.Ptr[int64](1), MaxLength: genai.Ptr[int64](100)},
					"direction": {Type: genai.TypeString, Enum: []string{"NORTH", "SOUTH", "EAST", "WEST"}, Title: "Direction", Description: "Direction of target"},
					"abc":       {Type: genai.TypeString, Pattern: "^[abc]$"},
					"email":     {Type: genai.TypeString, Format: "email"},
					"number":    {Type: genai.TypeNumber, Maximum: genai.Ptr(100.0)},
					"integer":   {Type: genai.TypeInteger, Minimum: genai.Ptr(1.0), Default: genai.Ptr[any](int64(-1))},
					"arrayOfString": {Type: genai.TypeArray,
						Items:    &genai.Schema{Type: genai.TypeString},
						MinItems: genai.Ptr[int64](1),
						MaxItems: genai.Ptr[int64](10),
					},
				},
				MinProperties: genai.Ptr[int64](3),
				MaxProperties: genai.Ptr[int64](4),
				Required:      []string{"direction", "email"},
				Title:         "Complex",
				Description:   "Example of complex schema generation",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := ForValue(test.value)
			if err != nil {
				t.Errorf("ForValue() error = %v", err)
			}

			if diff := cmp.Diff(got, test.want, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ForValue() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
