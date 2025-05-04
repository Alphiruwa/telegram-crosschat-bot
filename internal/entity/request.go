package entity

import (
	"errors"
	"time"
)

type Request struct {
	SrcChatID    int64
	TgtChatID    int64
	CreatedAt    time.Time
	TgtMessageID int64
	FromUserID   int64
}

type RequestRepository interface {
	IsRequestExists(srcChatID, tgtChatID int64) (bool, error)
	CreateRequest(srcChatID, tgtChatID, tgtMessageID, fromUserID int64) error
	DeleteRequest(srcChatID, tgtChatID int64) error
	GetRequest(srcChatID, tgtChatID int64) (*Request, error)
	GetAllChatOutRequests(chatID int64) ([]*Request, error)
	GetAllChatIncRequests(chatID int64) ([]*Request, error)
}

var (
	ErrRequestExists   = errors.New("request already exists")
	ErrRequestNotFound = errors.New("request not found")
)
