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
func RateLimitMiddleware(rdb *redis.Client, ctx context.Context, key string) error {
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
	fmt.Println("contRequest", countRequest)
	if countRequest > 3 {
		return ErrRequestLimitExceeded
	}
	pipe := rdb.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Minute-1)

	_, err = pipe.Exec(ctx)
	fmt.Println(incr.Val(), err)
	if err != nil {
		return ErrInternal
	}
	return nil
}
