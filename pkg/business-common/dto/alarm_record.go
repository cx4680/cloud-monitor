package dto

type UserContactInfo struct {
	ContactId   string `gorm:"COLUMN:contactId" json:"contactId"`
	ContactName string `gorm:"COLUMN:contactName" json:"contactName"`
	Phone       string `gorm:"COLUMN:phone" json:"phone"`
	Mail        string `gorm:"COLUMN:mail" json:"mail"`
}

type ContactGroupInfo struct {
	GroupId   string            `gorm:"COLUMN:groupId" json:"groupId"`
	GroupName string            `gorm:"COLUMN:groupName" json:"groupName"`
	Contacts  []UserContactInfo `gorm:"-" json:"contactList"`
}

type AutoScalingData struct {
	TenantId        string
	RuleId          string
	ResourceGroupId string
	Param           string
}
