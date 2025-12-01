package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/redis/go-redis/v9"
)

type PortfolioStore struct {
	rdb *redis.Client
}

func (ps *PortfolioStore) Get(ctx context.Context, portfolioID int64) (*store.Portfolio, error) {

	cacheKey := fmt.Sprintf("user-%v", portfolioID)

	data, err := ps.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var portfolio store.Portfolio
	if data != "" {
		err := json.Unmarshal([]byte(data), &portfolio)
		if err != nil {
			return nil, err
		}
	}
	return &portfolio, nil
}
func (ps *PortfolioStore) Set(ctx context.Context, portfolio *store.Portfolio) error {
	cacheKey := fmt.Sprintf("user-%v", portfolio.ID)

	data, err := json.Marshal(portfolio)
	if err != nil {
		return err
	}

	err = ps.rdb.SetEx(ctx, cacheKey, data, PortfolioExpTime).Err()
	if err != nil {
		return err
	}
	return nil
}
