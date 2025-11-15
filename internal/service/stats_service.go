package service

import (
	"context"

	"Avito/internal/repository"
	"Avito/internal/schema"
)

type StatsServiceImpl struct {
	statsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) *StatsServiceImpl {
	return &StatsServiceImpl{
		statsRepo: statsRepo,
	}
}

func (s *StatsServiceImpl) GetStatistics(ctx context.Context) (*schema.Statistics, error) {
	stats, err := s.statsRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *StatsServiceImpl) GetUserStats(ctx context.Context, userID string) (*schema.UserStats, error) {
	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
