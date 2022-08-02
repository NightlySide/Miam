package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

func NewRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language": {"fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3"},
		"DNT":             {"1"},
		"Connection":      {"keep-alive"},
	}

	return req
}

func Scrape(url string) (*ScrapedRecipe, error) {
	// Request the page
	req := NewRequest(url)
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// find the data
	jsonMetadata := doc.Find("script[type=\"application/ld+json\"]").First().Text()
	logrus.Debugf("Recipe json+ld: %s", jsonMetadata)
	var recipe ScrapedRecipe
	err = json.Unmarshal([]byte(jsonMetadata), &recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

type ScrapedRecipe struct {
	Context            string    `json:"@context"`
	Type               string    `json:"@type"`
	Name               string    `json:"name"`
	RecipeCategory     string    `json:"recipeCategory"`
	Image              []string  `json:"image"`
	DatePublished      time.Time `json:"datePublished"`
	PrepTime           string    `json:"prepTime"`
	CookTime           string    `json:"cookTime"`
	TotalTime          string    `json:"totalTime"`
	RecipeYield        string    `json:"recipeYield"`
	RecipeIngredient   []string  `json:"recipeIngredient"`
	RecipeInstructions []struct {
		Type string `json:"@type"`
		Text string `json:"text"`
	} `json:"recipeInstructions"`
	Author          string `json:"author"`
	Description     string `json:"description"`
	Keywords        string `json:"keywords"`
	RecipeCuisine   string `json:"recipeCuisine"`
	AggregateRating struct {
		Type        string  `json:"@type"`
		ReviewCount int     `json:"reviewCount"`
		RatingValue float64 `json:"ratingValue"`
		WorstRating int     `json:"worstRating"`
		BestRating  int     `json:"bestRating"`
	} `json:"aggregateRating"`
	Video struct {
		Type         string   `json:"@type"`
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		ThumbnailURL []string `json:"thumbnailUrl"`
		ContentURL   string   `json:"contentUrl"`
		EmbedURL     string   `json:"embedUrl"`
		UploadDate   string   `json:"uploadDate"`
	} `json:"video"`
}
