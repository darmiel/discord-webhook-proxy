package cms

import (
	"time"
)

type CMSPageUpdate struct {
	// Time of update
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	// UserID of updater, -1 = System
	UpdaterUserID string `bson:"updater_user_id" json:"updater_user_id"`
	// Patch of content
	Patch string `bson:"patch" json:"patch"`
}

func (p *CMSPage) GetLastUpdate() (res *CMSPageUpdate) {
	for _, u := range p.Updates {
		unix := u.UpdatedAt.Unix()
		if res == nil || unix > res.UpdatedAt.Unix() {
			res = &u
		}
	}
	return
}

type CMSUpdateInfo struct {
	FormattedTime string
	UpdaterID     string
}

func (u *CMSPageUpdate) GetInfo() *CMSUpdateInfo {
	return &CMSUpdateInfo{
		FormattedTime: u.UpdatedAt.Format("02.01.2006 15:04"),
		UpdaterID:     u.UpdaterUserID,
	}
}
