package rate

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"time"
)

var ErrInternal = errors.New("internal error")
var ErrRequestLimitExceeded = errors.New("request limit exceeded")

// Execute
//     MULTI
//     INCR pipeline_counter
//     EXPIRE pipeline_counts 59
//     EXEC
// using one rdb-server roundtrip.
func RateLimit(rdb *redis.Client, ctx context.Context, key string, limit int, perTime time.Duration) error {
	val, err := rdb.Get(ctx, key).Result()
	countRequest := 0
	if err != nil {
		if !strings.Contains(err.Error(), "redis: nil") {
			fmt.Println(err)
			return ErrInternal
		}
	}
	if len(val) > 0 {
		countRequest, err = strconv.Atoi(val)
		if err != nil {
			fmt.Println(err)
			return ErrInternal
		}
	}
	if countRequest > limit {
		return ErrRequestLimitExceeded
	}
	pipe := rdb.TxPipeline()
	_ = pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, perTime)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return ErrInternal
	}
	return nil
}

func PerSecond() time.Duration {
	return time.Hour - 1
}

func PerMinute() time.Duration {
	return time.Minute - 1
}

func PerHour() time.Duration {
	return time.Hour - 1
}

func PerDay() time.Duration {
	return 24*time.Hour - 1
}
