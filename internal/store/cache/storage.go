package cache

import (
	"context"
	"time"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/redis/go-redis/v9"
)

var (
	UserExpTime      = time.Minute
	PortfolioExpTime = time.Minute
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
	Portfolio interface {
		Get(context.Context, int64) (*store.Portfolio, error)
		Set(context.Context, *store.Portfolio) error
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users:     &UserStore{rdb: rbd},
		Portfolio: &PortfolioStore{rdb: rbd},
	}
}
