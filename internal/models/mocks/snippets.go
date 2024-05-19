package mocks

import (
	"github.com/google/uuid"
	"snippetbox.doichevkostia.dev/internal/models"
	"time"
)

var SnippetID = uuid.New()

var mockSnippet = models.Snippet{
	ID:         SnippetID,
	Title:      "An old silent pond",
	Content:    "An old silent pond...",
	CreateTime: time.Now(),
	ExpireTime: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *SnippetModel) Get(id uuid.UUID) (models.Snippet, error) {
	switch id {
	case SnippetID:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
