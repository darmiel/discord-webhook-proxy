package discord

type attributes struct {
	MaxWebhookCount uint   `json:"max_webhook_count" bson:"max_webhook_count"`
	Permissions     uint64 `json:"permissions" bson:"permissions"`
}

func NewDefaultAttributes() *attributes {
	return &attributes{
		MaxWebhookCount: 20,
		Permissions:     0,
	}
}

func (a *attributes) Repair() (updated bool) {
	if a.MaxWebhookCount == 0 {
		a.MaxWebhookCount = 20
		updated = true
	}
	// Permissions: 0
	return
}
