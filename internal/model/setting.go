package model

const (
	SINGLE = iota
	SITE
	STYLE
	PREVIEW
	GLOBAL
	_ // 5: reserved (was OFFLINE_DOWNLOAD), kept to preserve the group numbering
	INDEX
	SSO
	LDAP
	_ // 9: reserved (was S3), kept to preserve the group numbering
	_ // 10: reserved (was FTP), kept to preserve the group numbering
	TRAFFIC
)

const (
	PUBLIC = iota
	PRIVATE
	READONLY
	DEPRECATED
)

type SettingItem struct {
	Key            string `json:"key" gorm:"primaryKey" binding:"required"` // unique key
	Value          string `json:"value"`                                    // value
	MigrationValue string `json:"-" gorm:"-:all"`                           // deprecated value
	Help           string `json:"help"`                                     // help message
	Type           string `json:"type"`                                     // string, number, bool, select
	Options        string `json:"options"`                                  // values for select
	Group          int    `json:"group"`                                    // use to group setting in frontend
	Flag           int    `json:"flag"`                                     // 0 = public, 1 = private, 2 = readonly, 3 = deprecated, etc.
	Index          uint   `json:"index"`
}

func (s SettingItem) IsDeprecated() bool {
	return s.Flag == DEPRECATED
}
