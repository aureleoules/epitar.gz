package models

import (
	"time"

	"github.com/aureleoules/epitar/db"
	"github.com/jmoiron/sqlx"
)

type FileOrigin struct {
	FileID      string `db:"file_id" json:"file_id" `
	Module      string `db:"module" json:"module"`
	OriginalURL string `db:"original_url" json:"original_url"`
}

type FileMeta struct {
	ID        string       `db:"id" json:"id"`
	Name      string       `db:"name" json:"name"`
	Summary   string       `db:"summary" json:"summary"`
	Size      int64        `db:"size" json:"size"`
	Origins   []FileOrigin `db:"-" json:"origins"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
}

func GetFileOrigins(fileID string) ([]FileOrigin, error) {
	var origins []FileOrigin
	err := db.DB.Select(&origins, "SELECT * FROM file_origins WHERE file_id = ?", fileID)
	if err != nil {
		return nil, err
	}

	return origins, nil
}

func GetFileMeta(id string) (*FileMeta, error) {
	var f FileMeta
	err := db.DB.Get(&f, "SELECT * FROM files WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	f.Origins, err = GetFileOrigins(id)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func GetFilesMeta(ids []string) ([]FileMeta, error) {
	var files []FileMeta
	q, args, err := sqlx.In("SELECT f.* FROM files AS f WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	err = db.DB.Select(&files, q, args...)
	if err != nil {
		return nil, err
	}

	for i := range files {
		files[i].Origins, err = GetFileOrigins(files[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (f *FileMeta) Save() error {
	_, err := db.DB.Exec("INSERT INTO files (id, name, summary, size) VALUES (?, ?, ?, ?)", f.ID, f.Name, f.Summary, f.Size)
	return err
}

func (f *FileMeta) AddOrigin(module string, originalURL string) error {
	_, err := db.DB.Exec("INSERT INTO file_origins (file_id, module, original_url) VALUES (?, ?, ?)", f.ID, module, originalURL)
	return err
}
