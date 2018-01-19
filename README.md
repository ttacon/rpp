rpp
===

`rpp` is a helper library for creating a Redis Pool (that is prepared).


## Usage

```go
const (
    MAX_ACTIVE = 5
    MAX_IDLE   = 5
)

var (
    err error
    rpp *redis.Pool
)

if rpp, err = RPP("redis://x:foo@10.10.10.10:8443/5", MAX_ACTIVE, MAX_IDLE); err != nil {
    fmt.Println("There was an error, err: ", err)
    return
}

conn := rpp.Get()

_, _ = conn.Do("HGETALL", "data:foo")
```
