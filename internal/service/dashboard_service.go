package service

import "ticketing-system/internal/config"

type DashboardStats struct {
	Total  int64 `json:"total"`
	Open   int64 `json:"open"`
	Closed int64 `json:"closed"`
}

func GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats

	// total ticket
	config.DB.Model(&struct{}{}).Table("tickets").Count(&stats.Total)

	// open
	config.DB.Table("tickets").Where("status = ?", "open").Count(&stats.Open)

	// closed
	config.DB.Table("tickets").Where("status = ?", "closed").Count(&stats.Closed)

	return stats, nil
}
