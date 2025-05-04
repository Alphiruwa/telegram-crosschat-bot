package entity

import (
	"errors"
	"time"
)

type Link struct {
	SrcChatID int64
	TgtChatID int64
	CreatedAt time.Time
}

type LinkRepository interface {
	CreateLink(srcChatID, tgtChatID int64) error
	IsLinkExists(chatID1, chatID2 int64) (bool, error)
	GetAllChatLinks(chatID int64) ([]*Link, error)
	DeleteLink(chatID1, chatID2 int64) error
}

var (
	ErrLinkExists   = errors.New("link already exists")
	ErrLinkNotFound = errors.New("link not found")
)
