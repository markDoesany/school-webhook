package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"school-assistant-wh/internal/config"
)

type Service struct {
	config config.FacebookConfig
}

func NewService(cfg config.FacebookConfig) *Service {
	return &Service{
		config: cfg,
	}
}

func (s *Service) SendTextMessage(recipientID, text string) error {
	url := "https://graph.facebook.com/v18.0/me/messages"

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": text,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("access_token", s.config.PageAccessToken)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("facebook API error: %s - %s", resp.Status, string(body))
	}

	log.Printf("Message sent to %s: %s", recipientID, text)
	return nil
}

func (s *Service) SendWithPayload(recipientID string, payload map[string]interface{}) error {
	if _, ok := payload["recipient"]; !ok {
		payload["recipient"] = map[string]string{
			"id": recipientID,
		}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	url := "https://graph.facebook.com/v18.0/me/messages"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("access_token", s.config.PageAccessToken)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("facebook API error: %s - %s", resp.Status, string(body))
	}

	log.Printf("Payload sent to %s", recipientID)
	return nil
}

type QuickReply struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title,omitempty"`
	Payload     string `json:"payload,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

func (s *Service) SendQuickReplies(recipientID, text string, quickReplies []QuickReply) error {
	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"messaging_type": "RESPONSE",
		"message": map[string]interface{}{
			"text":          text,
			"quick_replies": quickReplies,
		},
	}

	return s.SendWithPayload(recipientID, payload)
}

func (s *Service) GetUserProfile(userID string) (*UserProfile, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s", userID)

	fields := "id,name,first_name,last_name,email,picture.type(large)"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add("fields", fields)
	q.Add("access_token", s.config.PageAccessToken)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Facebook API error: %s", string(body))
	}

	var profile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &profile, nil
}

func (s *Service) SetGetStartedButton(payload string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/me/messenger_profile?access_token=%s", s.config.PageAccessToken)

	requestBody := map[string]interface{}{
		"get_started": map[string]string{
			"payload": payload,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling get started button: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error setting get started button: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Facebook API error: %s", string(body))
	}

	return nil
}

func (s *Service) SetGreetingText(text string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/me/messenger_profile?access_token=%s", s.config.PageAccessToken)

	requestBody := map[string]interface{}{
		"greeting": []map[string]string{
			{
				"locale": "default",
				"text":   text,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling greeting text: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error setting greeting text: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Facebook API error: %s", string(body))
	}

	return nil
}

// SendImage sends an image to the specified recipient using the image URL
func (s *Service) SendImage(recipientID, imageURL string) error {
	if imageURL == "" {
		return fmt.Errorf("image URL cannot be empty")
	}

	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		return fmt.Errorf("image URL must start with http:// or https://")
	}

	if _, err := url.ParseRequestURI(imageURL); err != nil {
		return fmt.Errorf("invalid image URL format: %v", err)
	}

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "image",
				"payload": map[string]string{
					"url":         imageURL,
					"is_reusable": "true",
				},
			},
		},
	}

	return s.SendWithPayload(recipientID, payload)
}
