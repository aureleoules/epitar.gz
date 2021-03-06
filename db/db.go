package db

import (
	"os"
	"path"

	"github.com/aureleoules/epitar/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS files (
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	size INTEGER NOT NULL,
	summary TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS file_origins (
	file_id VARCHAR(255) NOT NULL,
	module VARCHAR(255) NOT NULL,
	original_url VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS news (
	id VARCHAR(255) PRIMARY KEY,
	from_user VARCHAR(255) NOT NULL,
	newsgroups VARCHAR(255) NOT NULL,
	subject VARCHAR(255) NOT NULL,
	date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	message_id VARCHAR(255) NOT NULL,
	size INTEGER NOT NULL,
	summary TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

var DB *sqlx.DB

func Init() {
	os.MkdirAll(path.Join(config.Cfg.Index.Store, "files"), 0755)
	DB = sqlx.MustOpen("sqlite3", path.Join(config.Cfg.Index.Store, "store.db"))

	DB.MustExec(schema)
}

func Close() {
	DB.Close()
}
