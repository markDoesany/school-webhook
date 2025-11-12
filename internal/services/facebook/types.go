package facebook

type UserProfile struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Picture   struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture,omitempty"`
}

type WebhookRequest struct {
	Object string `json:"object"`
	Entry  []struct {
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Message *struct {
				Text string `json:"text"`
			} `json:"message,omitempty"`
			Postback *struct {
				Payload string `json:"payload"`
			} `json:"postback,omitempty"`
		} `json:"messaging"`
	} `json:"entry"`
}
