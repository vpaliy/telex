package model

import (
	"github.com/jinzhu/gorm"
)

const (
	Dialog  = "dialog"
	Group   = "group"
	General = "general"
)

type Channel struct {
	gorm.Model
	Name          string `gorm:"unique_index;not null"`
	Tags          []string
	Creator       User `gorm:"foreignkey:CreatorID"`
	CreatorID     uint
	Image         *string
	Description   string `gorm:"size:2048"`
	Type          string
	Archived      bool
	Subscriptions []Subscription
}

type Subscription struct {
	gorm.Model
	Alert     bool
	Unread    int16
	Channel   Channel
	ChannelID uint
	User      User
	UserID    uint
	Type      string
}
