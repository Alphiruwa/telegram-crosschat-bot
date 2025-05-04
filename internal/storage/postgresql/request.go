package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Alphiruwa/telegram-crosschat-bot/internal/entity"
)

type RequestRepository struct {
	db *pgxpool.Pool
}

func NewRequestRepository(db *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{db}
}

func (repo *RequestRepository) IsRequestExists(srcChatID, tgtChatID int64) (bool, error) {
	row := repo.db.QueryRow(context.Background(), "SELECT FROM requests WHERE src_chat_id=$1 AND tgt_chat_id=$2", srcChatID, tgtChatID)
	if err := row.Scan(); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *RequestRepository) CreateRequest(srcChatID, tgtChatID, tgtMessageID, fromUserID int64) error {
	_, err := repo.db.Exec(context.Background(), "INSERT INTO requests (src_chat_id, tgt_chat_id, tgt_message_id, from_user_id) VALUES ($1, $2, $3, $4)", srcChatID, tgtChatID, tgtMessageID, fromUserID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return entity.ErrRequestExists
		}
		return err
	}
	return err
}

func (repo *RequestRepository) GetRequest(srcChatID, tgtChatID int64) (*entity.Request, error) {
	row := repo.db.QueryRow(context.Background(), "SELECT * FROM requests WHERE src_chat_id=$1 AND tgt_chat_id=$2", srcChatID, tgtChatID)
	req := &entity.Request{}
	if err := row.Scan(&req.SrcChatID, &req.TgtChatID, &req.CreatedAt, &req.TgtMessageID, &req.FromUserID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.ErrRequestNotFound
		}
		return nil, err
	}
	return req, nil
}

func (repo *RequestRepository) GetAllChatOutRequests(chatID int64) ([]*entity.Request, error) {
	rows, err := repo.db.Query(context.Background(), "SELECT * FROM requests WHERE src_chat_id=$1", chatID)
	if err != nil {
		return []*entity.Request{}, err
	}
	defer rows.Close()
	var requests []*entity.Request
	for rows.Next() {
		req := &entity.Request{}
		if err := rows.Scan(&req.SrcChatID, &req.TgtChatID, &req.CreatedAt, &req.TgtMessageID, &req.FromUserID); err != nil {
			return []*entity.Request{}, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (repo *RequestRepository) GetAllChatIncRequests(chatID int64) ([]*entity.Request, error) {
	rows, err := repo.db.Query(context.Background(), "SELECT * FROM requests WHERE tgt_chat_id=$1", chatID)
	if err != nil {
		return []*entity.Request{}, err
	}
	defer rows.Close()
	var requests []*entity.Request
	for rows.Next() {
		req := &entity.Request{}
		if err := rows.Scan(&req.SrcChatID, &req.TgtChatID, &req.CreatedAt, &req.TgtMessageID, &req.FromUserID); err != nil {
			return []*entity.Request{}, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (repo *RequestRepository) DeleteRequest(srcChatID, tgtChatID int64) error {
	ct, err := repo.db.Exec(context.Background(), "DELETE FROM requests WHERE src_chat_id=$1 AND tgt_chat_id=$2", srcChatID, tgtChatID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return entity.ErrRequestNotFound
	}
	return nil
}
