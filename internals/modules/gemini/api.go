package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// VideoMetadata represents the structured output for video generation
type VideoMetadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	VideoPrompt string   `json:"prompt"`
}

// GeminiRequest represents the request structure for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response structure from Gemini API
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

// ProcessQuoteWithGemini takes a quote and generates video metadata using Gemini Flash 2.5
func ProcessQuoteWithGemini(quote string) (*VideoMetadata, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Construct the prompt
	prompt := fmt.Sprintf(`You are a creative content strategist. Given the following quote, generate metadata for creating a compelling video.

Quote: "%s"

Generate a JSON response with the following structure:
{
  "title": "A catchy, engaging title for the video (max 60 characters)",
  "description": "A compelling description that explains the quote's meaning and relevance (2-3 sentences)",
  "tags": ["tag1", "tag2", "tag3", "tag4", "tag5"],
  "prompt": "A detailed prompt for an AI video generation model describing the visual scenes, mood, style, and cinematography that would best represent this quote"
}

Requirements:
- Title should be attention-grabbing and SEO-friendly
- Description should provide context and value
- Tags should be relevant and popular search terms (maximum 5 tags)
- Video prompt should be detailed and specific, describing visuals, atmosphere, colors, movements, and style

Return ONLY the JSON object, no additional text or explanation.`, quote)

	// Prepare the request body
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini API")
	}

	// Extract the JSON text from the response
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Parse the VideoMetadata from the response
	var metadata VideoMetadata
	if err := json.Unmarshal([]byte(responseText), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse video metadata from response: %w", err)
	}

	return &metadata, nil
}
