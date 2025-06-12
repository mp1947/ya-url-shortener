package usecase

import "fmt"

func ExampleGenerateIDFromURL() {
	url := "https://example.com/some/long/url"
	id := GenerateIDFromURL(url)

	// Note: The actual output may vary based on the URL provided.

	fmt.Println(id)

	// Output: e40ad0b5
}
