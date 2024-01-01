package entity

import (
	"math"

	"github.com/redis/go-redis/v9"
)

const (
	CachePageSize = 100
)

type RedisClient interface {
	redis.Cmdable
}

type Range struct {
	Start int64
	End   int64
}

func (r Range) Diff() int {
	return int(r.End - r.Start)
}

func NewRange(from, to uint64) Range {
	return Range{
		Start: offsetToRange(from),
		End:   offsetToRange(to),
	}
}

func offsetToRange(num uint64) int64 {
	return int64(math.Ceil(float64(num) / CachePageSize))
}
