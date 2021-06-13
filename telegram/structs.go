package telegram

type TGMessage struct {
	Message struct {
		Chat struct {
			FirstName string `json:"first_name"`
			ID        int64  `json:"id"`
			Type      string `json:"type"`
			Username  string `json:"username"`
		} `json:"chat"`
		Date     int64 `json:"date"`
		Entities []struct {
			Length int64  `json:"length"`
			Offset int64  `json:"offset"`
			Type   string `json:"type"`
		} `json:"entities"`
		From struct {
			FirstName    string `json:"first_name"`
			ID           int64  `json:"id"`
			IsBot        bool   `json:"is_bot"`
			LanguageCode string `json:"language_code"`
			Username     string `json:"username"`
		} `json:"from"`
		MessageID int64  `json:"message_id"`
		Text      string `json:"text"`
	} `json:"message"`
	UpdateID int64 `json:"update_id"`
}

type Webhookres struct {
	Ok          bool   `json: "ok"`
	Description string `json: "description"`
}
