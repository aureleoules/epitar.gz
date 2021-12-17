package models

import (
	"time"

	"github.com/aureleoules/epitar/db"
	"github.com/jmoiron/sqlx"
)

type FileMeta struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Summary   string    `db:"summary" json:"summary"`
	Size      int64     `db:"size" json:"size"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func GetFileMeta(id string) (*FileMeta, error) {
	var f FileMeta
	err := db.DB.Get(&f, "SELECT * FROM files WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func GetFilesMeta(ids []string) ([]FileMeta, error) {
	var files []FileMeta
	q, args, err := sqlx.In("SELECT * FROM files WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	err = db.DB.Select(&files, q, args...)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (f *FileMeta) Save() error {
	_, err := db.DB.Exec("INSERT INTO files (id, name, summary, size) VALUES (?, ?, ?, ?)", f.ID, f.Name, f.Summary, f.Size)
	return err
}
