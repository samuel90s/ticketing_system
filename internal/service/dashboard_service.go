package service

import (
	"ticketing-system/internal/config"
	"ticketing-system/internal/model"
)

type DashboardStats struct {
	Total      int64 `json:"total"`
	Open       int64 `json:"open"`
	Replied    int64 `json:"replied"`
	OnProgress int64 `json:"on_progress"`
	Done       int64 `json:"done"`
}

func GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats

	if err := config.DB.Model(&model.Ticket{}).Count(&stats.Total).Error; err != nil {
		return stats, err
	}
	if err := config.DB.Model(&model.Ticket{}).Where("status = ?", "open").Count(&stats.Open).Error; err != nil {
		return stats, err
	}
	if err := config.DB.Model(&model.Ticket{}).Where("status = ?", "replied").Count(&stats.Replied).Error; err != nil {
		return stats, err
	}
	if err := config.DB.Model(&model.Ticket{}).Where("status = ?", "on_progress").Count(&stats.OnProgress).Error; err != nil {
		return stats, err
	}
	if err := config.DB.Model(&model.Ticket{}).Where("status = ?", "done").Count(&stats.Done).Error; err != nil {
		return stats, err
	}

	return stats, nil
}
