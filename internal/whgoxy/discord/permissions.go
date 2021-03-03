package discord

import "log"

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
	PermissionCMSViewPageAuthor
	PermissionCMSViewHistory

	// Admin
	PermissionAdminDashboardView
)

// Permission Packs
const (
	PermissionPackWebhook        = PermissionWebhookCreate | PermissionWebhookEdit | PermissionWebhookDelete
	PermissionPackBasic          = PermissionLogin | PermissionPackWebhook
	PermissionPackCMSAdmin       = PermissionCMSCreatePage | PermissionCMSEditPage | PermissionCMSViewPageUpdates | PermissionCMSViewPageAuthor | PermissionCMSViewHistory
	PermissionPackAdminDashboard = PermissionAdminDashboardView
	PermissionPackAdmin          = PermissionPackBasic | PermissionPackCMSAdmin | PermissionAdminDashboardView
)

func (p Permission) Func() func(u *DiscordUser) bool {
	return func(u *DiscordUser) bool {
		log.Println("User", u.GetFullName(), "requested permission", p)
		res := p.Has(u)
		log.Println("  -> Res:", res)
		return res
	}
}

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
