package discord

type attributes struct {
	MaxWebhookCount  uint       `json:"max_webhook_count" bson:"max_webhook_count"`
	PermissionsAllow Permission `json:"permissions_allow" bson:"permissions_allow"`
	PermissionsDeny  Permission `json:"permissions_deny" bson:"permissions_deny"`
}

func NewDefaultAttributes() *attributes {
	return &attributes{
		MaxWebhookCount:  20,
		PermissionsAllow: 0,
		PermissionsDeny:  0,
	}
}

func (a *attributes) Repair() (updated bool) {
	if a.MaxWebhookCount == 0 {
		a.MaxWebhookCount = 20
		updated = true
	}
	if a.PermissionsAllow == 0 {
		a.PermissionsAllow = PermissionPackBasic
	}
	// PermissionsDeny: 0
	return
}
