package models

import (
	"strings"
	"time"

	"github.com/aureleoules/epitar/db"
	"github.com/jmoiron/sqlx"
)

type News struct {
	ID         string    `json:"id" db:"id"`
	From       string    `json:"from" db:"from_user"`
	Newsgroups string    `json:"newsgroups" db:"newsgroups"`
	Subject    string    `json:"subject" db:"subject"`
	Date       time.Time `json:"date" db:"date"`
	MessageID  string    `json:"message_id" db:"message_id"`
	Size       int       `json:"size" db:"size"`
	Summary    string    `json:"summary" db:"summary"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func GetNews(id string) (*News, error) {
	var news News
	err := db.DB.Get(&news, "SELECT * FROM news WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (f *News) Save() error {
	_, err := db.DB.Exec("INSERT INTO news (id, from_user, newsgroups, subject, date, message_id, size, summary) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", f.ID, f.From, f.Newsgroups, f.Subject, f.Date, f.MessageID, f.Size, f.Summary)
	return err
}

func GetNewsList(ids []string) ([]News, error) {
	var news []News
	q, args, err := sqlx.In("SELECT f.* FROM news AS f WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	err = db.DB.Select(&news, q, args...)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func GetUniqueNewsgroups() ([]string, error) {
	var newsgroups []string
	err := db.DB.Select(&newsgroups, "SELECT DISTINCT newsgroups FROM news")
	if err != nil {
		return nil, err
	}

	var uniqueNewsgroups []string
	for _, newsgroup := range newsgroups {
		if strings.Contains(newsgroup, ",") {
			continue
		}

		uniqueNewsgroups = append(uniqueNewsgroups, newsgroup)
	}

	return uniqueNewsgroups, nil
}
