package models

import "github.com/aureleoules/epitar/db"

type Stats struct {
	TotalFiles int `json:"total_files"`
	TotalNews  int `json:"total_news"`
	TotalSize  int `json:"total_size"`
}

func GetStats() (Stats, error) {
	var fileStats Stats
	err := db.DB.QueryRow("SELECT COUNT(*) AS total_files, SUM(size) AS total_size FROM files").Scan(&fileStats.TotalFiles, &fileStats.TotalSize)
	var newsStats Stats
	err = db.DB.QueryRow("SELECT COUNT(*) AS total_news, SUM(size) AS total_size FROM news").Scan(&newsStats.TotalNews, &newsStats.TotalSize)

	return Stats{
		TotalFiles: fileStats.TotalFiles,
		TotalNews:  newsStats.TotalNews,
		TotalSize:  fileStats.TotalSize + newsStats.TotalSize,
	}, err
}
