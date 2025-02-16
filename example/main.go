package main

// Port of controlled generation examples
// https://cloud.google.com/vertex-ai/generative-ai/docs/multimodal/control-generated-output

import (
	"context"
	"fmt"
	"log"

	"github.com/k0kubun/pp/v3"
	"google.golang.org/genai"

	"github.com/apstndb/genaischema"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

type T6 string

func (T6) Enum() []any {
	return []any{"drama", "comedy", "documentary"}
}

func run(ctx context.Context) error {
	const model = "gemini-2.0-flash"

	client, err := genai.NewClient(ctx, &genai.ClientConfig{HTTPOptions: genai.HTTPOptions{APIVersion: "v1"}})
	if err != nil {
		return err
	}

	{
		fmt.Println("Example: Send a prompt with a response schema")

		type T1 struct {
			RecipeName string `json:"recipe_name" required:"true"`
		}

		prompt := "List a few popular cookie recipes"

		ret, err := genaischema.GenerateTextContents[[]T1](ctx, client, model, genai.Text(prompt),
			&genai.GenerateContentConfig{ResponseMIMEType: "application/json", CandidateCount: genai.Ptr[int64](3), Temperature: genai.Ptr(0.6)})
		if err != nil {
			return err
		}

		pp.Println(ret)
	}

	{
		fmt.Println("Example: Send a prompt with a response schema")

		type T1 struct {
			RecipeName string `json:"recipe_name" required:"true"`
		}

		prompt := "List a few popular cookie recipes"

		ret, err := genaischema.GenerateObjectContent[[]T1](ctx, client, model, genai.Text(prompt), nil)
		if err != nil {
			return err
		}

		pp.Println(ret)
	}

	{
		fmt.Println("Example: Summarize review ratings")

		type T2 struct {
			Rating int    `json:"rating"`
			Flavor string `json:"flavor"`
		}

		prompt := `
			Reviews from our social media:

			- "Absolutely loved it! Best ice cream I've ever had." Rating: 4, Flavor: Strawberry Cheesecake
			- "Quite good, but a bit too sweet for my taste." Rating: 1, Flavor: Mango Tango
			`

		ret, err := genaischema.GenerateObjectContent[[]T2](ctx, client, model, genai.Text(prompt), nil)
		if err != nil {
			return err
		}

		pp.Println(ret)
	}

	{
		fmt.Println("Example: Forecast the weather for each day of the week")

		type T3Forecast struct {
			Day         string `json:"Day" required:"true"`
			Forecast    string `json:"Forecast" required:"true"`
			Humidity    string `json:"Humidity"`
			Temperature string `json:"Temperature" required:"true"`
			WindSpeed   string `json:"Wind Speed"`
		}

		type T3 struct {
			Forecast []T3Forecast `json:"forecast"`
		}

		prompt := `
			The week ahead brings a mix of weather conditions.
			Sunday is expected to be sunny with a temperature of 77°F and a humidity level of 50%. Winds will be light at around 10 km/h.
			Monday will see partly cloudy skies with a slightly cooler temperature of 72°F and humidity increasing to 55%. Winds will pick up slightly to around 15 km/h.
			Tuesday brings rain showers, with temperatures dropping to 64°F and humidity rising to 70%. Expect stronger winds at 20 km/h.
			Wednesday may see thunderstorms, with a temperature of 68°F and high humidity of 75%. Winds will be gusty at 25 km/h.
			Thursday will be cloudy with a temperature of 66°F and moderate humidity at 60%. Winds will ease slightly to 18 km/h.
			Friday returns to partly cloudy conditions, with a temperature of 73°F and lower humidity at 45%. Winds will be light at 12 km/h.
			Finally, Saturday rounds off the week with sunny skies, a temperature of 80°F, and a humidity level of 40%. Winds will be gentle at 8 km/h.
			`

		ret, err := genaischema.GenerateObjectContent[T3](ctx, client, model, genai.Text(prompt), nil)
		if err != nil {
			return err
		}

		pp.Println(ret)
	}

	{
		fmt.Println("Example: Classify a product")

		type T4 struct {
			ToDiscard    int    `json:"to_discard"`
			Subcategory  string `json:"subcategory"`
			SafeHandling string `json:"safe_handling"`
			ItemCategory string `json:"item_category" enum:"clothing,winter apparel,specialized apparel,furniture,decor,tableware,cookware,toys"`
			ForResale    int    `json:"for_resale"`
			Condition    string `json:"condition" enum:"new in package,like new,gently used,used,damaged,soiled"`
		}

		prompt := `
			Item description:
			The item is a long winter coat that has many tears all around the seams and is falling apart.
			It has large questionable stains on it.
			`

		ret, err := genaischema.GenerateObjectContent[[]T4](ctx, client, model, genai.Text(prompt), nil)
		if err != nil {
			return err
		}

		pp.Println(ret)
	}

	{
		fmt.Println("Example: Classify a product")

		type T5 struct {
			Object string `json:"object"`
		}

		img1 := genai.NewPartFromURI(
			"gs://cloud-samples-data/generative-ai/image/office-desk.jpeg",
			"image/jpeg",
		)

		img2 := genai.NewPartFromURI(
			"gs://cloud-samples-data/generative-ai/image/gardening-tools.jpeg",
			"image/jpeg",
		)

		prompt := "Generate a list of objects in the images."

		contents := []*genai.Content{{Parts: []*genai.Part{img1, img2, genai.NewPartFromText(prompt)}}}

		res, err := genaischema.GenerateObjectContent[T5](ctx, client, model, contents, nil)
		if err != nil {
			return err
		}

		pp.Println(res)
	}

	{
		fmt.Println("Example: Example schema for enum output")

		prompt := `
The film aims to educate and inform viewers about real-life subjects, events, or people.
It offers a factual record of a particular topic by combining interviews, historical footage,
and narration. The primary purpose of a film is to present information and provide insights
into various aspects of reality.
`

		res, err := genaischema.GenerateEnumContent[T6](ctx, client, model, genai.Text(prompt), nil)
		if err != nil {
			return err
		}

		fmt.Println(res)
	}
	return nil
}
