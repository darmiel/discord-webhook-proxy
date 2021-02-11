package cms

import (
	"time"
)

type CMSPageUpdate struct {
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`           // Time of update
	UpdaterUserID string    `bson:"updater_user_id" json:"updater_user_id"` // UserID of updater, -1 = System
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
