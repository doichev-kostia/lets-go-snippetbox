package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (uuid.UUID, error)
	Get(id uuid.UUID) (Snippet, error)
	Latest() ([]Snippet, error)
}

type Snippet struct {
	ID         uuid.UUID
	Title      string
	Content    string
	CreateTime time.Time
	ExpireTime time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (uuid.UUID, error) {
	stmt := `insert into "snippets" (id, title, content, create_time, expire_time)
	values (?, ?, ?, current_timestamp, datetime(current_timestamp, ?))`

	id := uuid.New()
	expiration := fmt.Sprintf("+%d days", expires)
	_, err := m.DB.Exec(stmt, id, title, content, expiration)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id uuid.UUID) (Snippet, error) {

	stmt := `select "id", "title", "content", "create_time", "expire_time" from "snippets"
	where expire_time > current_timestamp and id = ?`

	var s Snippet

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.CreateTime, &s.ExpireTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `select "id", "title", "content", "create_time", "expire_time" from "snippets"
	where expire_time > current_timestamp order by create_time desc limit 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreateTime, &s.ExpireTime)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
