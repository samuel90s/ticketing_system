package service

import (
	"ticketing-system/internal/config"
	"ticketing-system/internal/model"
)

type DashboardStats struct {
	Total      int64 `json:"total"`
	Open       int64 `json:"open"`
	InProgress int64 `json:"in_progress"`
	Closed     int64 `json:"closed"`
}

func GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats

	if err := config.DB.Model(&model.Ticket{}).Count(&stats.Total).Error; err != nil {
		return stats, err
	}

	if err := config.DB.Model(&model.Ticket{}).Where("status = ?", "open").Count(&stats.Open).Error; err != nil {
		return stats, err
	}

	if err := config.DB.Model(&model.Ticket{}).Where("status = ?", "in_progress").Count(&stats.InProgress).Error; err != nil {
		return stats, err
	}

	if err := config.DB.Model(&model.Ticket{}).Where("status = ?", "closed").Count(&stats.Closed).Error; err != nil {
		return stats, err
	}

	return stats, nil
}
