package discord

type Permission uint64

const (
	_ Permission = (1 << iota) - 1

	// CMS
	PermissionCMSEditPage
	PermissionCMSViewPageUpdates
)

func (p Permission) Has(u *DiscordUser) bool {
	return u.HasPermission(p)
}

func (u *DiscordUser) HasPermission(perm Permission) bool {
	return (u.Attributes.Permissions & uint64(perm)) == uint64(perm)
}
