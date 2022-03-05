package rate

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	minute := time.Now().Minute()
	key := "zA21X31:" + strconv.Itoa(minute)
	err := RateLimit(rdb, ctx, key, 3, PerMinute())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("handle next")
}
