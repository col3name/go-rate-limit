# Rate limiting for go-redis

## Installation

redis_rate supports 2 last Go versions and requires a Go version with
[modules](https://github.com/golang/go/wiki/Modules) support. So make sure to initialize a Go module:

```shell
go mod init github.com/my/repo
```

And then install rate-limit:

```shell
go get github.com/go-redis/redis/v8
go get github.com/col3name/rate-limit
```

## Example Usage

```go
package main

import (
    "context"
    "fmt"
    rateLimit "github.com/col3name/rate-limit"
    "github.com/go-redis/redis/v8"
    "strconv"
    "time"
)

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    minute := time.Now().Minute()
    key := "zA21X31:" + strconv.Itoa(minute)
    err := rateLimit.RateLimit(rdb, ctx, key, 3, rateLimit.PerMinute())
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("handle next")
}
```

### Example Middleware

```go
package main

import (
    "context"
    rateLimit "github.com/col3name/rate-limit"
    "github.com/go-redis/redis/v8"
    "net/http"
    "strconv"
    "time"
)

func rateLimitMiddleware(h http.Handler, rdb *redis.Client) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, 4096)

        minute := time.Now().Minute()
        token := r.Header.Get("AccessToken")

        if len(token) == 0 {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        key := token + ":" + strconv.Itoa(minute)
        ctx := context.Background()
        err := rateLimit.RateLimit(rdb, ctx, key, 3, rateLimit.PerMinute())
        if err != nil {
            w.WriteHeader(http.StatusTooManyRequests)
            return
        }

        h.ServeHTTP(w, r)
    })
}

```