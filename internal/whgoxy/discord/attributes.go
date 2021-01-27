package discord

type attributes struct {
	MaxWebhookCount uint `json:"max_webhook_count" bson:"max_webhook_count"`
}

func NewDefaultAttributes() *attributes {
	return &attributes{
		MaxWebhookCount: 20,
	}
}

func (a *attributes) Repair() (updated bool) {
	if a.MaxWebhookCount == 0 {
		a.MaxWebhookCount = 20
		updated = true
	}
	return
}
