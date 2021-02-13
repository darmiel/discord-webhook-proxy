package discord

type Permission uint64

// Permissions
const (
	// Basic Permissions
	PermissionLogin Permission = 1 << iota

	// Webhook
	PermissionWebhookCreate
	PermissionWebhookEdit
	PermissionWebhookDelete

	// CMS
	PermissionCMSCreatePage
	PermissionCMSEditPage
	PermissionCMSViewPageUpdates
)

// Permission Packs
const (
	PermissionPackWebhook  = PermissionWebhookCreate | PermissionWebhookEdit | PermissionWebhookDelete
	PermissionPackBasic    = PermissionLogin | PermissionPackWebhook
	PermissionPackCMSAdmin = PermissionCMSCreatePage | PermissionCMSEditPage | PermissionCMSViewPageUpdates
	PermissionPackAdmin    = PermissionPackBasic | PermissionPackCMSAdmin
)

func (p Permission) Has(u *DiscordUser) bool {
	return u.HasPermission(p)
}

func (u *DiscordUser) HasPermission(perm Permission) bool {
	if u == nil {
		return false
	}
	if u.Attributes == nil {
		u.Repair()
	}
	attr := u.Attributes
	if (attr.PermissionsDeny & perm) == perm {
		return false
	}
	return (attr.PermissionsAllow & perm) == perm
}
