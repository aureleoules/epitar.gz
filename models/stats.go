package models

import "github.com/aureleoules/epitar/db"

type Stats struct {
	TotalFiles int `json:"total_files"`
	TotalSize  int `json:"total_size"`
}

func GetStats() (Stats, error) {
	var stats Stats
	err := db.DB.QueryRow("SELECT COUNT(*) AS total_files, SUM(size) AS total_size FROM files").Scan(&stats.TotalFiles, &stats.TotalSize)
	return stats, err
}
