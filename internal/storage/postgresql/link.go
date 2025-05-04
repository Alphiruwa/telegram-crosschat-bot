package postgresql

import (
	"context"

	"github.com/Alphiruwa/telegram-crosschat-bot/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db}
}

func (repo *LinkRepository) IsLinkExists(chatID1, chatID2 int64) (bool, error) {
	row := repo.db.QueryRow(context.Background(), "SELECT FROM links WHERE src_chat_id=$1 AND tgt_chat_id=$2 OR src_chat_id=$2 AND tgt_chat_id=$1", chatID1, chatID2)
	if err := row.Scan(); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *LinkRepository) CreateLink(srcChatID, tgtChatID int64) error {
	if exists, err := repo.IsLinkExists(srcChatID, tgtChatID); exists {
		return entity.ErrLinkExists
	} else if err != nil {
		return err
	}
	_, err := repo.db.Exec(context.Background(), "INSERT INTO links (src_chat_id, tgt_chat_id) VALUES ($1, $2)", srcChatID, tgtChatID)
	return err
}

func (repo *LinkRepository) GetAllChatLinks(chatID int64) ([]*entity.Link, error) {
	rows, err := repo.db.Query(context.Background(), "SELECT * FROM links WHERE src_chat_id=$1 UNION SELECT * FROM links WHERE tgt_chat_id=$1", chatID)
	if err != nil {
		return []*entity.Link{}, err
	}
	var links []*entity.Link
	defer rows.Close()
	for rows.Next() {
		link := &entity.Link{}
		if err := rows.Scan(&link.SrcChatID, &link.TgtChatID, &link.CreatedAt); err != nil {
			return []*entity.Link{}, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (repo *LinkRepository) DeleteLink(chatID1, chatID2 int64) error {
	ct, err := repo.db.Exec(context.Background(), "DELETE FROM links WHERE src_chat_id=$1 AND tgt_chat_id=$2 OR src_chat_id=$2 AND tgt_chat_id=$1", chatID1, chatID2)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return entity.ErrLinkNotFound
	}
	return nil
}
